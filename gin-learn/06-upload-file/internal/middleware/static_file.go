package middleware

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

// StaticFileWithFilename 返回一个处理静态文件请求的中间件，支持自定义下载文件名
// 参数:
//   - rootDir: 静态文件的根目录路径
//
// 返回值:
//   - gin.HandlerFunc: 处理静态文件请求的Gin中间件函数
func StaticFileWithFilename(rootDir string) gin.HandlerFunc {
	// 使用http.FileSystem创建文件服务
	fileServer := http.FileServer(gin.Dir(rootDir, false))

	// 返回中间件处理函数
	return func(c *gin.Context) {
		// 获取URL路径参数
		filePath := c.Param("filepath")
		// 检查请求URL中是否有filename参数
		filename := c.Query("filename")

		// 如果提供了filename参数，设置Content-Disposition头
		if filename != "" {
			// URL解码文件名，以支持中文等特殊字符
			decodedFilename, err := url.QueryUnescape(filename)
			if err == nil {
				filename = decodedFilename
			}

			// 根据文件类型确定是展示还是下载
			// 图片、视频等通常是展示，其他文件通常是下载
			contentDisposition := "attachment"
			ext := strings.ToLower(filepath.Ext(filePath))
			if isPreviewableType(ext) {
				contentDisposition = "inline"
			}

			// 设置Content-Disposition头，指定文件名
			c.Header("Content-Disposition", contentDisposition+"; filename=\""+filename+"\"")
		}

		// 修改请求路径，去掉前缀，以便http.FileServer正确处理
		c.Request.URL.Path = filePath

		// 使用http.FileServer处理请求
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

// isPreviewableType 判断文件扩展名是否是可预览的类型
func isPreviewableType(ext string) bool {
	previewableExts := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
		".bmp":  true,
		".webp": true,
		".svg":  true,
		".mp4":  true,
		".webm": true,
		".ogg":  true,
		".mp3":  true,
		".wav":  true,
		".pdf":  true,
	}

	return previewableExts[ext]
}
