package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ListUsers 处理获取所有用户
// @Summary 获取所有用户
// @Description 获取系统中的所有用户列表
// @Tags 用户管理
// @Accept json
// @Produce json
// @Success 200 {array} model.UserResponse "成功返回用户列表"
// @Failure 500 {object} map[string]string "服务器内部错误"
// @Router /users [get]
func (h *Handler) ListUsers(c *gin.Context) {
	users, err := h.userService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户失败"})
		return
	}

	c.JSON(http.StatusOK, users)
}
