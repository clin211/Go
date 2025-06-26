package v1

import (
	"github.com/clin211/gin-learn/06-upload-file/api/v1/files"
	"github.com/clin211/gin-learn/06-upload-file/api/v1/images"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {

	f := &files.Files{}
	imageController := &images.ImageController{}

	api := router.Group("/api/v1")
	{
		// 文件上传
		uploadRouter := api.Group("/upload")
		{
			uploadRouter.POST("/file", f.File)               // 上传单文件
			uploadRouter.POST("/multiple", f.Multiple)       // 上传多个文件
			uploadRouter.POST("/folder", f.Folder)           // 上传文件夹
			uploadRouter.POST("/chunk/init", f.ChunkInit)    // 初始化分片上传，返回 uploadID
			uploadRouter.POST("/chunk", f.Chunk)             // 上传分片（服务端自动触发合并）
			uploadRouter.GET("/chunk/status", f.ChunkStatus) // 查询分片上传状态（可选，用于断点续传）
			uploadRouter.GET("/url", f.GetFileURL)           // 获取文件访问地址（支持 MinIO 对象存储）
			uploadRouter.GET("/meta/:file_name", f.Meta)     // 获取文件元信息（如大小、类型等）
			uploadRouter.DELETE("/:file_name", f.Delete)     // 删除文件（通过 file_name 或 object_key）
		}

		// 文件下载
		filesRouter := api.Group("/files")
		{
			filesRouter.GET("/download", f.Download) // 下载文件（支持自定义文件名）
		}

		// 图片处理
		imageRouter := api.Group("/images")
		{
			imageRouter.POST("/compress", imageController.Compress)                  // 压缩单张图片
			imageRouter.POST("/compress/multiple", imageController.CompressMultiple) // 批量压缩图片
		}
	}
}
