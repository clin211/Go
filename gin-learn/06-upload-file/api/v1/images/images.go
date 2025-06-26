package images

import (
	"github.com/gin-gonic/gin"
)

type ImageController struct{}

type ImageAction interface {
	Compress(c *gin.Context)         // 压缩单张图片
	CompressMultiple(c *gin.Context) // 批量压缩图片
}

// 检查是否实现了 ImageAction 接口
var _ ImageAction = &ImageController{}
