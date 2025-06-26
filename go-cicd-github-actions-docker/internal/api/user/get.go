package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUser 处理获取单个用户
// @Summary 获取单个用户
// @Description 根据用户ID获取用户详细信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} model.UserResponse "成功返回用户信息"
// @Failure 400 {object} map[string]string "无效的用户ID"
// @Failure 404 {object} map[string]string "用户不存在"
// @Router /users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	user, err := h.userService.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, user)
}
