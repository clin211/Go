package files

import "github.com/gin-gonic/gin"

func (u *Files) Multiple(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "multiple ok",
	})
}
