package files

import "github.com/gin-gonic/gin"

func (u *Files) Meta(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "meta ok",
	})
}
