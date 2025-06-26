package files

import (
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/clin211/gin-learn/06-upload-file/internal/config"
	"github.com/gin-gonic/gin"
)

// Download 处理文件下载请求
// 支持通过objectKey下载文件并指定下载文件名
// 请求参数:
//   - objectKey: 文件存储的对象键 (必填)
//   - filename: 下载时的文件名 (可选，默认使用objectKey的文件名部分)
//
// 响应:
//   - 成功: 返回文件内容，设置适当的Content-Disposition头
//   - 失败: 返回错误信息
func (f *Files) Download(c *gin.Context) {
	// 获取查询参数
	objectKey := c.Query("objectKey")
	if objectKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "必须提供objectKey参数",
		})
		return
	}

	// 获取可选的filename参数，如果没有提供，使用objectKey的文件名部分
	filename := c.Query("filename")
	if filename == "" {
		filename = filepath.Base(objectKey)
	}

	// URL解码文件名，以支持中文等特殊字符
	decodedFilename, err := url.QueryUnescape(filename)
	if err == nil {
		filename = decodedFilename
	}

	// 根据文件类型确定是展示还是下载
	contentDisposition := "attachment"
	ext := strings.ToLower(filepath.Ext(filename))
	if isPreviewableType(ext) {
		contentDisposition = "inline"
	}

	// 设置Content-Disposition头，指定文件名
	c.Header("Content-Disposition", contentDisposition+"; filename=\""+filename+"\"")

	// 使用静态文件服务器提供文件
	// 直接重定向到静态文件服务器，让它处理文件内容的提供
	c.File(filepath.Join(config.Local.UploadDir, objectKey))
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
