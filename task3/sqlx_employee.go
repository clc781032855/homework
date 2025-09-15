package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"log"
)

// Employee 结构体映射 employees 表
type Employee struct {
	ID         int    `db:"id"`
	Name       string `db:"name"`
	Department string `db:"department"`
	Salary     int    `db:"salary"`
}

// 创建 employees 表
func createEmployeesTable(db *sqlx.DB) {
	schema := `
	CREATE TABLE IF NOT EXISTS employees (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(50) NOT NULL,
		department VARCHAR(50) NOT NULL,
		salary INT NOT NULL
	)`

	_, err := db.Exec(schema)
	if err != nil {
		log.Fatalf("创建表失败: %v", err)
	}
}

// 插入测试数据
func insertTestData(db *sqlx.DB) {
	// 清空表
	db.MustExec("DELETE FROM employees")

	// 插入测试数据
	employees := []Employee{
		{Name: "张三", Department: "技术部", Salary: 15000},
		{Name: "李四", Department: "技术部", Salary: 18000},
		{Name: "王五", Department: "市场部", Salary: 12000},
		{Name: "赵六", Department: "财务部", Salary: 22000},
		{Name: "钱七", Department: "技术部", Salary: 20000},
	}

	for _, emp := range employees {
		_, err := db.NamedExec(`
			INSERT INTO employees (name, department, salary)
			VALUES (:name, :department, :salary)
		`, emp)
		if err != nil {
			log.Fatalf("插入数据失败: %v", err)
		}
	}
}

// 查询特定部门的所有员工
func getEmployeesByDepartment(db *sqlx.DB, department string) ([]Employee, error) {
	var employees []Employee

	err := db.Select(&employees,
		"SELECT id, name, department, salary FROM employees WHERE department = ?",
		department)

	if err != nil {
		return nil, err
	}
	return employees, nil
}

// 查询工资最高的员工
func getHighestPaidEmployee(db *sqlx.DB) (Employee, error) {
	var employee Employee

	// 使用sqlx.Get获取单条记录
	err := db.Get(&employee,
		"SELECT id, name, department, salary FROM employees ORDER BY salary DESC LIMIT 1")

	if err != nil {
		return Employee{}, err
	}
	return employee, nil
}

func main() {
	// 连接到数据库
	db, err := sqlx.Connect("mysql", "root:123456@tcp(127.0.0.1:3306)/company?parseTime=true")
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	defer db.Close()

	// 创建 employees 表
	createEmployeesTable(db)

	// 插入测试数据
	insertTestData(db)

	// 1. 查询技术部所有员工
	fmt.Println("技术部员工列表:")
	techEmployees, err := getEmployeesByDepartment(db, "技术部")
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	for _, emp := range techEmployees {
		fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 薪水: %d\n",
			emp.ID, emp.Name, emp.Department, emp.Salary)
	}
	fmt.Println()

	// 2. 查询工资最高的员工
	fmt.Println("工资最高的员工:")
	highestPaidEmployee, err := getHighestPaidEmployee(db)
	if err != nil {
		log.Fatalf("查询失败: %v", err)
	}

	fmt.Printf("ID: %d, 姓名: %s, 部门: %s, 薪水: %d\n",
		highestPaidEmployee.ID, highestPaidEmployee.Name,
		highestPaidEmployee.Department, highestPaidEmployee.Salary)
}
