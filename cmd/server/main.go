package main

import (
	"log"
	"net/http"

	"gopass/internal/database"
	"gopass/internal/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	if err := database.InitDB("gopass.db"); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// 创建路由器
	r := gin.Default()

	// 添加中间件
	r.Use(handlers.CORSMiddleware())

	// 静态文件服务
	r.Static("/static", "./web/static")
	r.LoadHTMLGlob("web/templates/*")

	// 前端路由
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "GoPass - Password Manager",
		})
	})

	r.GET("/dashboard", func(c *gin.Context) {
		c.HTML(http.StatusOK, "dashboard.html", gin.H{
			"title": "Dashboard - GoPass",
		})
	})

	// API路由
	api := r.Group("/api")
	{
		// 公开路由
		api.POST("/register", handlers.Register)
		api.POST("/login", handlers.Login)

		// 需要认证的路由
		auth := api.Group("/")
		auth.Use(handlers.AuthMiddleware())
		{
			// 密码管理
			auth.POST("/passwords", handlers.CreatePassword)
			auth.GET("/passwords", handlers.GetPasswords)
			auth.GET("/passwords/:id", handlers.GetPassword)
			auth.PUT("/passwords/:id", handlers.UpdatePassword)
			auth.DELETE("/passwords/:id", handlers.DeletePassword)

			// 密码生成
			auth.POST("/generate-password", handlers.GeneratePassword)

			// 分类管理
			auth.GET("/categories", handlers.GetCategories)
			auth.POST("/categories", handlers.CreateCategory)

			// 数据导入导出
			auth.GET("/export", handlers.ExportData)
			auth.POST("/import", handlers.ImportData)
		}
	}

	log.Println("Server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
