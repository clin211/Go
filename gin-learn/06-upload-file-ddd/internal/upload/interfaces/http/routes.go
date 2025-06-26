package http

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册所有上传路由
func RegisterRoutes(router *gin.Engine, handler *UploadHandler) {
	// 提供静态HTML页面
	router.Static("/static", "./web")
	router.StaticFile("/", "./web/index.html")

	// 创建上传组
	uploadGroup := router.Group("/upload")
	{
		// 单文件上传
		uploadGroup.POST("/single", handler.UploadSingle)

		// 多文件上传
		uploadGroup.POST("/multiple", handler.UploadMultiple)

		// 文件夹上传
		uploadGroup.POST("/folder", handler.UploadFolder)

		// 分片上传
		chunkedGroup := uploadGroup.Group("/chunk")
		{
			// 初始化分片上传
			chunkedGroup.POST("/init", handler.InitChunkedUpload)

			// 上传单个分片
			chunkedGroup.POST("", handler.UploadChunk)

			// 合并分片
			chunkedGroup.POST("/merge", handler.MergeChunks)
		}
	}
}
