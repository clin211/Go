package validator

import (
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/clin211/gin-learn/06-upload-file/internal/config"
)

// ValidateFile 验证上传文件是否符合系统要求
// 主要检查两个方面：
// 1. 文件类型是否在允许列表中
// 2. 文件大小是否超过系统限制
//
// 参数:
//   - f: 上传的文件信息
//
// 返回值:
//   - error: 如果验证失败返回描述性错误；验证通过返回nil
func ValidateFile(f *multipart.FileHeader) error {
	// 从配置中获取允许的文件扩展名列表
	allowedExt := config.Local.AllowedExtensions

	// 提取上传文件的扩展名并转为小写（确保大小写不敏感的比较）
	ext := strings.ToLower(filepath.Ext(f.Filename))

	// 检查文件扩展名是否在允许列表中
	valid := false
	for _, a := range allowedExt {
		if a == ext {
			valid = true
			break
		}
	}

	// 如果文件类型不被允许，返回相应错误
	if !valid {
		return errors.New("文件类型不被允许")
	}

	// 检查文件大小是否超过系统配置的最大限制
	if f.Size > config.Local.MaxFileSize {
		return errors.New("文件大小超过限制")
	}

	// 验证通过，返回nil
	return nil
}
