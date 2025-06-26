package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"

	v1 "github.com/clin211/gin-learn/06-upload-file/api/v1"
	"github.com/clin211/gin-learn/06-upload-file/internal/config"
	"github.com/clin211/gin-learn/06-upload-file/internal/middleware"
)

func main() {
	// 初始化配置
	configPath := filepath.Join("configs", "config.yaml")
	if err := config.Init(configPath); err != nil {
		log.Fatalf("初始化配置失败: %v", err)
		return
	}

	// 现在可以通过 config.xx 来访问配置
	fmt.Println("config.Local.UploadDir: ", config.Local.UploadDir)

	// Gin 初始化
	gin.SetMode(config.Server.Mode)
	router := gin.Default()

	// 添加支持自定义下载文件名的静态文件处理
	router.GET("/static/*filepath", middleware.StaticFileWithFilename(config.Local.UploadDir))

	// 其他静态文件服务（不带下载功能）
	router.Static("/assets", "./api/v1/static")
	router.StaticFile("/", "./api/v1/static/index.html")

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "ok",
		})
	})

	v1.RegisterRoutes(router)

	// 支持从环境变量覆盖配置文件中的端口设置
	port := config.Server.Port
	if envPort := os.Getenv("PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	portStr := fmt.Sprintf(":%d", port)
	fmt.Println("server start at ", portStr)
	router.Run(portStr)
}
