package user

import (
	"github.com/clin211/go-cicd-github-actions-docker/internal/repository"
	"github.com/clin211/go-cicd-github-actions-docker/internal/service"
	"github.com/gin-gonic/gin"
)

// Handler 处理用户相关的所有HTTP请求
type Handler struct {
	userService *service.UserService
}

// NewHandler 创建一个新的用户处理器
func NewHandler(repo repository.UserRepository) *Handler {
	userService := service.NewUserService(repo)
	return &Handler{
		userService: userService,
	}
}

// RegisterRoutes 注册所有用户相关路由
func (h *Handler) RegisterRoutes(router *gin.RouterGroup) {
	users := router.Group("/users")
	{
		// 用户列表和创建
		users.GET("", h.ListUsers)
		users.POST("", h.SignUp)

		// 单个用户操作
		users.GET("/:id", h.GetUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}

	// 认证相关
	auth := router.Group("/auth")
	{
		auth.POST("/signin", h.SignIn)
		auth.POST("/signup", h.SignUp)
	}
}
