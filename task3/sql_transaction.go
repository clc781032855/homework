package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

type Account struct {
	ID      int
	Name    string
	Balance int
}

// 创建 students 表
func createTable_accounts(db *sql.DB) error {
	createSQL := `
	CREATE TABLE IF NOT EXISTS accounts (
    id INT AUTO_INCREMENT PRIMARY KEY,  
    balance INT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; `

	_, err := db.Exec(createSQL)
	return err
}

// 创建 transactions 表
func createTable_transactions(db *sql.DB) error {
	createSQL := `
	CREATE TABLE IF NOT EXISTS transactions (
    id INT AUTO_INCREMENT PRIMARY KEY,  
    from_account_id INT NOT NULL,
    to_account_id INT NOT NULL,
    amount INT NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; `

	_, err := db.Exec(createSQL)
	return err
}

// 检查账户余额
func checkBalance(db *sql.DB, accountID int) (int, error) {
	var balance int
	query := "SELECT balance FROM accounts WHERE id = ?"
	err := db.QueryRow(query, accountID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}

func updateBalance(db *sql.DB, accountID int, amount int) error {
	query := "UPDATE accounts SET balance = balance + ? WHERE id = ?"
	_, err := db.Exec(query, amount, accountID)
	return err
}

// 转账事务
func transfer(db *sql.DB, fromAccountID int, toAccountID int, amount int) error {
	// 检查转账金额是否大于0
	if amount <= 0 {
		return errors.New("转账金额必须大于0")
	}

	// 开始事务
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("开始事务失败: %v", err)
	}

	// 1. 检查转出账户余额 (使用FOR UPDATE锁定行)
	var fromBalance int
	err = tx.QueryRow("SELECT balance FROM accounts WHERE id = ? FOR UPDATE", fromAccountID).Scan(&fromBalance)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("查询转出账户余额失败: %v", err)
	}

	// 2. 检查余额是否足够
	if fromBalance < amount {
		tx.Rollback()
		return fmt.Errorf("账户 %d 余额不足 (当前余额: %d, 需要: %d)", fromAccountID, fromBalance, amount)
	}

	// 3. 更新转出账户余额
	_, err = tx.Exec("UPDATE accounts SET balance = balance - ? WHERE id = ?", amount, fromAccountID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新转出账户失败: %v", err)
	}

	// 4. 更新转入账户余额
	_, err = tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, toAccountID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("更新转入账户失败: %v", err)
	}

	// 5. 记录交易信息
	_, err = tx.Exec(`
		INSERT INTO transactions (from_account_id, to_account_id, amount)
		VALUES (?, ?, ?)
	`, fromAccountID, toAccountID, amount)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("记录交易失败: %v", err)
	}

	// 提交事务
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %v", err)
	}

	return nil
}

// 创建账户
func createAccount(db *sql.DB, balance int) error {
	_, err := db.Exec("INSERT INTO accounts (balance) VALUES (?)", balance)
	return err
}

// 初始化账户
func initAccounts(db *sql.DB) {
	// 清空表
	db.Exec("DELETE FROM transactions")
	db.Exec("DELETE FROM accounts")

	// 创建初始账户
	balances := []int{500, 300}

	for _, bal := range balances {
		if err := createAccount(db, bal); err != nil {
			log.Printf("初始化账户失败: %v\n", err)
		}
	}
}

func main() {
	// 连接到数据库
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8&parseTime=True")
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}
	defer db.Close()

	// 创建 accounts 表
	if err := createTable_accounts(db); err != nil {
		log.Fatal("创建accounts表失败:", err)
	}

	// 创建 transactions 表
	if err := createTable_transactions(db); err != nil {
		log.Fatal("创建transactions表失败:", err)
	}

	// 初始化测试账户
	initAccounts(db)
	fmt.Println("初始化账户完成：账户1余额500，账户2余额300")

	// 执行转账操作：从账户1向账户2转100元
	fmt.Println("\n开始执行转账：账户1 → 账户2 (100元)")

	err = transfer(db, 1, 2, 100)
	if err != nil {
		log.Fatal("转账失败:", err)
	}

	fmt.Println("转账成功！")

	// 查询转账后余额
	bal1, _ := checkBalance(db, 1)
	bal2, _ := checkBalance(db, 2)

	fmt.Println("\n转账后余额：")
	fmt.Printf("账户1: %d\n", bal1)
	fmt.Printf("账户2: %d\n", bal2)

	// 尝试转账失败的情况（余额不足）
	fmt.Println("\n尝试转账失败案例：账户1 → 账户2 (500元)")

	err = transfer(db, 1, 2, 500)
	if err != nil {
		fmt.Println("预期中的转账失败:", err)
	} else {
		fmt.Println("转账成功（这不应该发生）")
	}

	// 查询转账后余额（应保持不变）
	bal1, _ = checkBalance(db, 1)
	bal2, _ = checkBalance(db, 2)

	fmt.Println("\n转账尝试后余额（应保持不变）：")
	fmt.Printf("账户1: %d\n", bal1)
	fmt.Printf("账户2: %d\n", bal2)
}
