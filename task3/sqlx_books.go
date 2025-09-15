package main

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Book 结构体映射 books 表字段
type Book struct {
	ID     int     `db:"id"`     // 主键ID
	Title  string  `db:"title"`  // 书名
	Author string  `db:"author"` // 作者
	Price  float64 `db:"price"`  // 价格
}

// 查询价格大于指定值的书籍
func getExpensiveBooks(db *sqlx.DB, minPrice float64) ([]Book, error) {
	var books []Book

	query := "SELECT id, title, author, price FROM books WHERE price > ?"

	err := db.Select(&books, query, minPrice)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func main() {

	db, err := sqlx.Connect("mysql", "root:123456@tcp(127.0.0.1:3306)/library?parseTime=true")
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	// 查询价格大于50元的书籍
	books, err := getExpensiveBooks(db, 50.0)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	fmt.Printf("\n价格大于 ¥50.00 的书籍:\n")
	fmt.Printf("%-5s %-40s %-20s %s\n", "ID", "书名", "作者", "价格")
	fmt.Println("--------------------------------------------------------------")

	for _, book := range books {
		fmt.Printf("%-5d %-40s %-20s ¥%6.2f\n",
			book.ID, book.Title, book.Author, book.Price)
	}
}
