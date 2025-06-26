package files

import "github.com/gin-gonic/gin"

func (u *Files) Chunk(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "chunk ok",
	})
}
