package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Tag struct{}

func NewTag() Tag {
	return Tag{}
}

func (t Tag) List(c *gin.Context)   {
	c.JSON(http.StatusOK, gin.H{
		"message": "success",
		"auther": "Forest",
	})
}
func (t Tag) Create(c *gin.Context) {}
func (t Tag) Update(c *gin.Context) {}
func (t Tag) Delete(c *gin.Context) {}
