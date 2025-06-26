package api

import (
	"github.com/clin211/go-cicd-github-actions-docker/internal/api/user"
	"github.com/clin211/go-cicd-github-actions-docker/internal/repository"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 设置所有API路由
func RegisterRoutes(r *gin.Engine) {
	// 创建依赖
	userRepo := repository.NewMySQLUserRepository()
	userHandler := user.NewHandler(userRepo)

	// API版本组
	v1 := r.Group("/api/v1")

	// 注册用户模块路由
	userHandler.RegisterRoutes(v1)

	// 健康检查
	v1.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}
