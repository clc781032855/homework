package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"os"
	"strings"
	"time"
)

// User 用户
type User struct {
	gorm.Model
	Username string    `gorm:"type:text;not null"`
	Password string    `gorm:"not null"`
	Email    string    `gorm:"not null"`
	Posts    []Post    // 一对多：一个用户有多篇文章
	Comments []Comment // 一对多：一个用户有多条评论
}

// Post 文章
type Post struct {
	gorm.Model
	Title    string    `gorm:"not null"`
	Content  string    `gorm:"type:text;not null"`
	UserID   uint      `gorm:"not null"` // 外键：关联User.ID
	User     User      // 属于：一篇文章属于一个用户
	Comments []Comment // 一对多：一篇文章有多条评论
}

// Comment 评论
type Comment struct {
	gorm.Model
	Content string `gorm:"type:text;not null"`
	UserID  uint   `gorm:"not null"` // 外键：关联User.ID
	PostID  uint   `gorm:"not null"` // 外键：关联Post.ID
	User    User   // 属于：一条评论属于一个用户
	Post    Post   // 属于：一条评论属于一篇文章
}

// 用户注册
func registerHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}

		hashpassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(500, gin.H{"error": "密码加密失败"})
			return
		}

		user := User{
			Username: input.Username,
			Password: string(hashpassword),
		}
		if err := db.Create(&user).Error; err != nil {
			c.JSON(500, gin.H{"error": "用户名已存在"})
			return
		}
		c.JSON(200, gin.H{"message": "用户注册成功"})
	}
}

// 用户登陆（生成JWT）
func loginHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var input struct {
			Username string `json:"username" binding:"required"` // 增加必填验证
			Password string `json:"password" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 查询用户
		var user User
		if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}

		// 验证密码
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
			return
		}

		// 生成JWT（关键步骤）
		claims := jwt.MapClaims{
			"sub": user.ID,                               // 核心：用户唯一标识（User.ID）
			"exp": time.Now().Add(24 * time.Hour).Unix(), // 过期时间：24小时（Unix时间戳）
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
			return
		}

		// 返回token（客户端需保存，后续请求携带）
		c.JSON(http.StatusOK, gin.H{
			"message": "登录成功",
			"token":   tokenString, // 客户端需存储此token（如localStorage）
		})
	}
}

// 定义路由
// 定义路由（包含公开路由和保护路由）
func setupAuthRoutes(r *gin.Engine, db *gorm.DB) {
	// 1. 公开路由（无需认证）：注册、登录、获取文章列表等
	publicGroup := r.Group("/api/public")
	{
		publicGroup.POST("/register", registerHandler(db)) // 注册
		publicGroup.POST("/login", loginHandler(db))       // 登录
	}

	// 2. 保护路由（需要认证）：创建文章、获取用户信息等
	protectedGroup := r.Group("/api/protected")
	protectedGroup.Use(authMiddleware()) // 应用认证中间件（所有该组下的路由都需验证JWT）
	{
		protectedGroup.GET("/user", getUserInfoHandler(db))  // 获取当前用户信息
		protectedGroup.POST("/posts", createPostHandler(db)) // 创建文章
		//protectedGroup.PUT("/posts/:id", updatePostHandler(db)) // 更新文章（可选）
	}
}

// authMiddleware JWT认证中间件（复用性高）
func authMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取Authorization头（格式：Bearer <token>）
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "缺少Authorization头"})
			c.Abort() // 终止请求，不再执行后续逻辑
			return
		}

		// 2. 解析Authorization格式
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的Authorization格式（需为Bearer <token>）"})
			c.Abort()
			return
		}

		// 3. 提取token字符串
		tokenString := parts[1]

		// 4. 验证token有效性（签名+过期时间）
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// 验证签名算法（必须为HS256）
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("无效的签名算法")
			}
			return jwtSecret, nil // 使用密钥验证签名
		})

		// 5. 处理验证结果
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token（已过期或被伪造）"})
			c.Abort()
			return
		}

		// 6. 提取用户ID（存入gin.Context，供后续Handler使用）
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的token载荷"})
			c.Abort()
			return
		}
		userId := uint(claims["sub"].(float64)) // 将float64转为uint（User.ID类型）
		c.Set("userId", userId)                 // 存入上下文，后续Handler可通过c.Get("userId")获取

		// 7. 继续执行后续逻辑（如创建文章、获取用户信息）
		c.Next()
	}
}

// 获取当前用户信息（需要认证）
func getUserInfoHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从上下文获取用户ID（由authMiddleware存入）
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}

		// 查询当前用户信息（避免返回密码哈希）
		var user User
		if err := db.First(&user, userId).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
			return
		}
		user.Password = "" // 隐藏敏感信息（密码哈希）

		// 返回用户信息
		c.JSON(http.StatusOK, gin.H{"data": user})
	}
}

// 创建文章（需要认证，示例）
func createPostHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userId, exists := c.Get("userId")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
			return
		}

		// 解析文章参数
		var input struct {
			Title   string `json:"title" binding:"required"`
			Content string `json:"content" binding:"required"`
		}
		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 创建文章（关联当前用户）
		post := Post{
			Title:   input.Title,
			Content: input.Content,
			UserID:  userId.(uint), // 关联当前用户ID
		}
		if err := db.Create(&post).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
			return
		}

		// 返回创建成功的文章
		c.JSON(http.StatusCreated, gin.H{"data": post})
	}
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET"))

// 主函数
// 主函数（启动服务器）
func main() {
	// 1. 连接数据库
	dsn := "root:123456@tcp(127.0.0.1:3306)/blog?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("数据库连接失败:" + err.Error())
	}

	// 2. 自动迁移表结构（创建/更新users、posts、comments表）
	err = db.AutoMigrate(&User{}, &Post{}, &Comment{})
	if err != nil {
		panic("自动迁移失败:" + err.Error())
	}
	fmt.Println("表迁移成功")

	// 3. 初始化Gin路由器并设置路由
	r := gin.Default()     // 初始化Gin引擎
	setupAuthRoutes(r, db) // 加载路由配置

	// 4. 启动服务器（监听8080端口）
	fmt.Println("服务器启动成功，监听端口: 8080")
	if err := r.Run(":8080"); err != nil {
		panic("服务器启动失败:" + err.Error())
	}
}
