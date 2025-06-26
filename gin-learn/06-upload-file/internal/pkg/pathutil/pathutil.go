package pathutil

import (
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateFilePath 生成文件存储路径
// 参数:
//   - filename: 原始文件名
//
// 返回值:
//   - string: 生成的存储路径
func GenerateFilePath(filename string) string {
	// 获取当前日期并格式化为 YYYY/MM/DD 结构
	now := time.Now()
	datePath := now.Format("2006/01/02")

	// 生成UUID作为文件名的一部分，避免文件名冲突
	uniqueID := uuid.New().String()

	// 获取文件扩展名
	ext := filepath.Ext(filename)
	ext = strings.ToLower(ext)

	// 构建安全的文件名：UUID + 原始扩展名
	safeFilename := uniqueID + ext

	// 返回完整的存储路径
	return filepath.Join(datePath, safeFilename)
}
