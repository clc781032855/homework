package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"os"
)

type Student struct {
	ID    int
	Name  string
	Age   int
	Grade string
}

// 创建 students 表
func createTable(db *sql.DB) error {
	createSQL := `
	CREATE TABLE IF NOT EXISTS students (
    id INT AUTO_INCREMENT PRIMARY KEY,  
    name VARCHAR(50) NOT NULL,
    age INT NOT NULL,
    grade VARCHAR(50) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4; `

	_, err := db.Exec(createSQL)
	return err
}

// 插入新学生记录
func insertStudent(db *sql.DB) {
	var name, grade string
	var age int

	fmt.Print("请输入学生姓名: ")
	fmt.Scanln(&name)

	fmt.Print("请输入学生年龄: ")
	fmt.Scanln(&age)

	fmt.Print("请输入学生年级: ")
	fmt.Scanln(&grade)

	// 题目要求的插入操作
	insertSQL := "INSERT INTO students (name, age, grade) VALUES (?, ?, ?)"
	result, err := db.Exec(insertSQL, name, age, grade)
	if err != nil {
		log.Println("插入失败:", err)
		return
	}

	id, _ := result.LastInsertId()
	fmt.Printf("插入成功! 学生ID: %d\n", id)
}

// 查询学生信息
func queryStudents(db *sql.DB, condition string) {
	querySQL := "SELECT id, name, age, grade FROM students WHERE " + condition

	rows, err := db.Query(querySQL)
	if err != nil {
		log.Println("查询失败:", err)
		return
	}
	defer rows.Close()

	fmt.Println("\n学生列表:")
	fmt.Println("ID\t姓名\t年龄\t年级")
	fmt.Println("----------------------------")

	count := 0
	for rows.Next() {
		var s Student
		if err := rows.Scan(&s.ID, &s.Name, &s.Age, &s.Grade); err != nil {
			log.Println("读取数据失败:", err)
			continue
		}
		fmt.Printf("%d\t%s\t%d\t%s\n", s.ID, s.Name, s.Age, s.Grade)
		count++
	}

	fmt.Printf("共找到 %d 名学生\n", count)
}

// 更新学生年级
func updateStudentGrade(db *sql.DB) {
	var name, newGrade string

	fmt.Print("请输入要更新的学生姓名: ")
	fmt.Scanln(&name)

	fmt.Print("请输入新的年级: ")
	fmt.Scanln(&newGrade)

	updateSQL := "UPDATE students SET grade = ? WHERE name = ?"
	result, err := db.Exec(updateSQL, newGrade, name)
	if err != nil {
		log.Println("更新失败:", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		fmt.Printf("成功更新 %d 名学生的年级\n", rowsAffected)
	} else {
		fmt.Println("没有找到匹配的学生")
	}
}

// 删除年龄小于15岁的学生
func deleteYoungStudents(db *sql.DB) {
	deleteSQL := "DELETE FROM students WHERE age < 15"
	result, err := db.Exec(deleteSQL)
	if err != nil {
		log.Println("删除失败:", err)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	fmt.Printf("成功删除 %d 名年龄小于15岁的学生\n", rowsAffected)
}

func main() {
	//连接到数据库
	db, err := sql.Open("mysql", "root:123456@tcp(127.0.0.1:3306)/gorm")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	// 创建 students 表
	if err := createTable(db); err != nil {
		log.Fatal(err)
	}
	// 命令行交互
	for {
		fmt.Println("\n学生信息管理系统")
		fmt.Println("1. 插入新学生")
		fmt.Println("2. 查询年龄大于18岁的学生")
		fmt.Println("3. 更新学生年级")
		fmt.Println("4. 删除年龄小于15岁的学生")
		fmt.Println("5. 显示所有学生")
		fmt.Println("6. 退出")
		fmt.Print("请选择操作: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			insertStudent(db)
		case 2:
			queryStudents(db, "age > 18")
		case 3:
			updateStudentGrade(db)
		case 4:
			deleteYoungStudents(db)
		case 5:
			queryStudents(db, "1=1") // 查询所有学生
		case 6:
			fmt.Println("程序退出")
			os.Exit(0)
		default:
			fmt.Println("无效选择，请重新输入")
		}
	}
}
