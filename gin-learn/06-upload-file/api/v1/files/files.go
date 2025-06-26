package files

import (
	"github.com/gin-gonic/gin"
)

type Files struct{}

type FileAction interface {
	File(c *gin.Context)        // 上传单文件
	Multiple(c *gin.Context)    // 上传多个文件
	Folder(c *gin.Context)      // 上传文件夹
	Chunk(c *gin.Context)       // 上传分片
	ChunkInit(c *gin.Context)   // 初始化分片上传，返回 uploadID
	ChunkStatus(c *gin.Context) // 查询分片上传状态（可选，用于断点续传）
	Meta(c *gin.Context)        // 获取文件元信息（如大小、类型等）
	GetFileURL(c *gin.Context)  // 获取文件访问地址（支持 MinIO 对象存储）
	Delete(c *gin.Context)      // 删除文件（通过 file_id 或 object_key）
	Download(c *gin.Context)    // 下载文件（支持自定义文件名）
}

// 检查是否实现了 Uploader 接口
var _ FileAction = &Files{}
