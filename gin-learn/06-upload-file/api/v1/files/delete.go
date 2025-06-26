package files

import (
	"os"

	"github.com/gin-gonic/gin"
)

// Delete 通过文件名删除文件（先不考虑同名文件）
func (u *Files) Delete(c *gin.Context) {
	// 获取文件名
	fileName := c.Query("filename")

	// 删除文件
	err := os.Remove(fileName)
	if err != nil {
		c.JSON(500, gin.H{
			"message": "删除失败",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "删除成功",
	})
}
