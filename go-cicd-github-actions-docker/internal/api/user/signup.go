package user

import (
	"net/http"

	"github.com/clin211/go-cicd-github-actions-docker/internal/model"
	"github.com/gin-gonic/gin"
)

// SignUp 处理用户注册
// @Summary 用户注册
// @Description 创建新用户账号
// @Tags 用户管理,认证
// @Accept json
// @Produce json
// @Param user body model.UserCreate true "用户注册信息"
// @Success 201 {object} model.UserResponse "成功创建用户"
// @Failure 400 {object} map[string]string "请求参数错误或用户已存在"
// @Router /auth/signup [post]
func (h *Handler) SignUp(c *gin.Context) {
	var userData model.UserCreate
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Create(userData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, user)
}
