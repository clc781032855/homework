package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type User struct {
	gorm.Model
	Name       string `gorm:"size:100;not null"`
	Email      string `gorm:"size:255;not null;unique"`
	Password   string `gorm:"size:255;not null"`
	PostsCount int    `gorm:"default:0"` // 新增：文章数量统计字段
	Posts      []Post `gorm:"foreignKey:UserID"`
}

type Post struct {
	gorm.Model
	Title         string    `gorm:"size:200;not null"`
	Content       string    `gorm:"type:text;not null"`
	UserID        uint      `gorm:"index;not null"`
	Comments      []Comment `gorm:"foreignKey:PostID"`
	CommentStatus string    `gorm:"size:20;default:'no_comments'"` // 新增：评论状态字段
}

type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	PostID  uint   `gorm:"index;not null"`
	UserID  uint   `gorm:"index;not null"`
}

// 创建测试数据
func createTestData(db *gorm.DB) {
	// 清空现有数据
	db.Exec("DELETE FROM comments")
	db.Exec("DELETE FROM posts")
	db.Exec("DELETE FROM users")

	// 创建用户
	user := User{
		Name:     "张三",
		Email:    "zhangsan@example.com",
		Password: "pass123",
	}
	db.Create(&user)

	fmt.Printf("用户初始文章数量: %d\n", user.PostsCount)
}

// 演示钩子函数效果
func demonstrateHooks(db *gorm.DB) {

	var user User
	db.First(&user)

	// 创建文章
	fmt.Println("\n创建第一篇文章...")
	post1 := Post{
		Title:   "GORM钩子函数详解",
		Content: "本文讲解GORM钩子函数的使用方法...",
		UserID:  user.ID,
	}
	db.Create(&post1)

	// 查询用户更新后的文章数量
	db.First(&user)
	fmt.Printf("用户文章数量更新为: %d\n", user.PostsCount)

	// 添加评论
	fmt.Println("\n添加评论...")
	comment1 := Comment{
		Content: "很有帮助的文章！",
		PostID:  post1.ID,
		UserID:  user.ID,
	}
	db.Create(&comment1)

	// 检查文章评论状态
	var updatedPost Post
	db.First(&updatedPost, post1.ID)
	fmt.Printf("文章评论状态: %s\n", updatedPost.CommentStatus)

	// 删除评论 - 会触发AfterDelete钩子
	fmt.Println("\n删除评论...")
	db.Delete(&comment1)

	// 检查文章评论状态更新
	db.First(&updatedPost, post1.ID)
	fmt.Printf("删除评论后文章状态: %s\n", updatedPost.CommentStatus)

	// 创建第二篇文章（不带评论）
	fmt.Println("\n创建第二篇文章（不带评论）...")
	post2 := Post{
		Title:   "钩子函数最佳实践",
		Content: "分享钩子函数的使用技巧...",
		UserID:  user.ID,
	}
	db.Create(&post2)

	// 检查第二篇文章初始状态
	var newPost Post
	db.First(&newPost, post2.ID)
	fmt.Printf("新文章初始评论状态: %s\n", newPost.CommentStatus)
}

// Post的BeforeCreate钩子：更新用户的文章数量统计
func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Printf("执行Post.BeforeCreate钩子: 更新用户 %d 的文章数量\n", p.UserID)

	// 原子操作：用户的文章计数+1
	result := tx.Model(&User{}).
		Where("id = ?", p.UserID).
		Update("posts_count", gorm.Expr("posts_count + ?", 1))

	if result.Error != nil {
		return result.Error
	}

	// 设置初始评论状态
	p.CommentStatus = "no_comments"

	return nil
}

// Comment的AfterDelete钩子：检查文章的评论数量
func (c *Comment) AfterDelete(tx *gorm.DB) (err error) {
	fmt.Printf("执行Comment.AfterDelete钩子: 检查文章 %d 的评论数量\n", c.PostID)

	// 查询当前文章的评论数量
	var count int64
	if err := tx.Model(&Comment{}).
		Where("post_id = ?", c.PostID).
		Count(&count).Error; err != nil {
		return err
	}

	// 如果评论数量为0，更新文章评论状态
	if count == 0 {
		fmt.Printf("文章 %d 已无评论, 更新状态\n", c.PostID)
		return tx.Model(&Post{}).
			Where("id = ?", c.PostID).
			Update("comment_status", "no_comments").Error
	}

	return nil
}

func main() {

	dsn := "root:123456@tcp(127.0.0.1:3306)/gorm?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}

	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		log.Fatalf("自动迁移失败: %v", err)
	}

	fmt.Println("数据库表创建成功!")

	createTestData(db)

	demonstrateHooks(db)
}
