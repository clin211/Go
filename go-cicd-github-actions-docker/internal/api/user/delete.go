package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeleteUser 处理删除用户
// @Summary 删除用户
// @Description 根据ID删除用户
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Success 204 "删除成功，无内容返回"
// @Failure 400 {object} map[string]string "无效的用户ID"
// @Failure 404 {object} map[string]string "用户不存在"
// @Router /users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	err = h.userService.Delete(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.Status(http.StatusNoContent)
}
