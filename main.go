package main

import (
	"file-service/config"
	"file-service/controllers"
	"file-service/database"
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// 初始化数据库
	database.InitDB()
	defer database.CloseDB()

	// 创建Gin路由器
	router := gin.Default()

	// 创建文件控制器
	fileController := controllers.NewFileController()

	// 定义路由
	router.POST("/upload/", fileController.UploadFiles)
	router.GET("/download/:fileName", fileController.DownloadFile)
	router.GET("/files/", fileController.GetFileList)

	// 设置静态文件服务（使用更具体的路径前缀）
	router.Static("/static", "./static")

	// 为根路径添加重定向到静态文件首页
	router.GET("/", func(c *gin.Context) {
		c.Redirect(302, "/static/index.html")
	})

	// 启动服务器
	port := config.AppConfig.Server.Port
	log.Printf("Server starting on port %d", port)
	if err := router.Run(":" + strconv.Itoa(port)); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
