package user

import (
	"fmt"
	"log"
	"net/http"

	"github.com/clin211/go-cicd-github-actions-docker/internal/model"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SignInRequest 表示登录请求
type SignInRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// SignInResponse 表示登录响应
type SignInResponse struct {
	User  model.UserResponse `json:"user"`
	Token string             `json:"token"`
}

// SignIn 处理用户登录
// @Summary 用户登录
// @Description 验证用户凭据并返回令牌
// @Tags 认证
// @Accept json
// @Produce json
// @Param credentials body SignInRequest true "用户凭据"
// @Success 200 {object} SignInResponse "成功登录，返回用户信息和令牌"
// @Failure 400 {object} map[string]string "请求参数错误"
// @Failure 401 {object} map[string]string "认证失败"
// @Router /auth/signin [post]
func (h *Handler) SignIn(c *gin.Context) {
	var req SignInRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Printf("登录请求参数绑定失败: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("尝试登录: 用户名=%s, 密码长度=%d", req.Username, len(req.Password))

	// 根据用户名获取用户
	userResp, err := h.userService.GetByUsername(req.Username)
	if err != nil {
		log.Printf("用户 %s 不存在或获取失败: %v", req.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户名或密码"})
		return
	}
	log.Printf("成功获取用户 %s 的基本信息", req.Username)

	// 获取完整用户信息以验证密码
	user, err := h.userService.GetFullUserByUsername(req.Username)
	if err != nil {
		log.Printf("获取用户 %s 的完整信息失败: %v", req.Username, err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无效的用户名或密码"})
		return
	}
	log.Printf("成功获取用户 %s 的完整信息，存储密码长度=%d", req.Username, len(user.Password))

	// 验证密码
	log.Printf("尝试验证密码: 密码=%s, 存储哈希=%s", req.Password, user.Password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		log.Printf("密码验证失败: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("无效的用户名或密码: %v", err)})
		return
	}
	log.Printf("密码验证成功")

	// 这里应该生成JWT token，简化版本直接返回用户信息
	c.JSON(http.StatusOK, SignInResponse{
		User:  *userResp,
		Token: "sample-token-" + req.Username, // 在实际应用中应使用JWT
	})
}
