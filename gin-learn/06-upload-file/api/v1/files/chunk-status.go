package files

import "github.com/gin-gonic/gin"

// 查询分片上传状态（可选，用于断点续传）
func (u *Files) ChunkStatus(c *gin.Context) {
	// 获取文件名
	fileName := c.Query("filename")

	// 获取分片上传状态
	chunkStatus := c.Query("chunkStatus")

	c.JSON(200, gin.H{
		"message":     "chunk status ok",
		"fileName":    fileName,
		"chunkStatus": chunkStatus,
	})
}
