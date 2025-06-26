package files

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// 获取文件访问地址（支持 MinIO 对象存储）
func (u *Files) GetFileURL(c *gin.Context) {
	// 获取文件名
	fileName := c.Query("filename")

	// 获取文件访问地址
	fileURL := fmt.Sprintf("http://%s/%s", "127.0.0.1:8080", fileName)

	c.JSON(200, gin.H{
		"message": "ok",
		"fileURL": fileURL,
	})
}
