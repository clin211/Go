package storage

import (
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// LocalStorage 实现了Storage接口，提供本地文件系统存储功能
// 它将上传的文件保存到指定的本地目录中
type LocalStorage struct {
	BasePath string // 文件存储的根目录路径
}

// NewLocalStorage 创建一个新的LocalStorage实例
// 参数:
//   - basePath: 文件存储的根目录路径
//
// 返回值:
//   - 初始化好的LocalStorage指针
func NewLocalStorage(basePath string) *LocalStorage {
	return &LocalStorage{BasePath: basePath}
}

// Save 将上传的文件保存到本地文件系统
// 参数:
//   - fileHeader: 包含上传文件信息和数据的multipart.FileHeader
//   - dstPath: 目标存储路径（相对于BasePath的路径）
//
// 返回值:
//   - error: 如果保存过程中发生错误，返回相应的错误信息；否则返回nil
func (s *LocalStorage) Save(fileHeader *multipart.FileHeader, dstPath string) error {
	// 打开上传的文件
	src, err := fileHeader.Open()
	if err != nil {
		return err
	}
	defer src.Close() // 确保文件句柄被关闭，防止资源泄漏

	// 构建完整的文件存储路径
	fullPath := filepath.Join(s.BasePath, dstPath)

	// 创建必要的目录结构
	// 使用MkdirAll可以创建多级目录，确保存储路径存在
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	// 创建目标文件
	dst, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer dst.Close() // 确保文件句柄被关闭

	// 将上传的文件内容复制到目标文件
	// io.Copy实现了高效的数据传输，适用于大文件
	_, err = io.Copy(dst, src)
	return err
}

// SaveFromBytes 将字节数据保存为文件
// 参数:
//   - data: 要保存的字节数据
//   - dstPath: 目标存储路径（相对于BasePath的路径）
//
// 返回值:
//   - error: 如果保存过程中发生错误，返回相应的错误信息；否则返回nil
func (s *LocalStorage) SaveFromBytes(data []byte, dstPath string) error {
	// 构建完整的文件存储路径
	fullPath := filepath.Join(s.BasePath, dstPath)

	// 创建必要的目录结构
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return err
	}

	// 创建并写入文件
	return os.WriteFile(fullPath, data, 0644)
}

// GetURL 获取已存储文件的访问URL
// 参数:
//   - objectKey: 文件的唯一标识符/路径
//
// 返回值:
//   - string: 文件的访问URL（本地存储通常返回相对路径）
//   - error: 如果生成URL过程中发生错误，返回相应的错误信息；否则返回nil
func (s *LocalStorage) GetURL(objectKey string) (string, error) {
	// 本地存储无法直接生成公网访问URL
	// 返回一个相对路径，需要配合Web服务器（如Nginx）提供静态文件服务
	return "/static/" + objectKey, nil
}

// GetURLWithFilename 获取已存储文件的访问URL，并指定下载时的文件名
// 参数:
//   - objectKey: 文件的唯一标识符/路径
//   - filename: 下载时显示的文件名
//
// 返回值:
//   - string: 文件的访问URL（本地存储通常返回相对路径）
//   - error: 如果生成URL过程中发生错误，返回相应的错误信息；否则返回nil
func (s *LocalStorage) GetURLWithFilename(objectKey string, filename string) (string, error) {
	// 生成访问URL（在前端下载时将使用指定的文件名）
	// 这里我们添加一个查询参数，这个参数会在处理下载请求时被使用
	return "/static/" + objectKey + "?filename=" + filename, nil
}

// Delete 从本地文件系统中删除指定文件
// 参数:
//   - objectKey: 要删除的文件的唯一标识符/路径
//
// 返回值:
//   - error: 如果删除过程中发生错误，返回相应的错误信息；否则返回nil
func (s *LocalStorage) Delete(objectKey string) error {
	// 构建完整的文件路径并执行删除操作
	return os.Remove(filepath.Join(s.BasePath, objectKey))
}
