package user

import (
	"net/http"
	"strconv"

	"github.com/clin211/go-cicd-github-actions-docker/internal/model"
	"github.com/gin-gonic/gin"
)

// UpdateUser 处理更新用户
// @Summary 更新用户信息
// @Description 根据ID更新用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param user body model.UserUpdate true "用户更新信息"
// @Success 200 {object} model.UserResponse "成功返回更新后的用户信息"
// @Failure 400 {object} map[string]string "请求参数错误或业务错误"
// @Failure 404 {object} map[string]string "用户不存在"
// @Router /users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var userData model.UserUpdate
	if err := c.ShouldBindJSON(&userData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.Update(uint(id), userData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
