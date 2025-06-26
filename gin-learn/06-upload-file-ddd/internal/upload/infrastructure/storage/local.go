package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"06-upload-file/internal/upload/domain"

	"github.com/google/uuid"
)

// LocalStorageConfig 本地文件存储的配置
type LocalStorageConfig struct {
	UploadDir    string
	TempChunkDir string
	MaxFileSize  int64
	AllowedTypes []string
}

// LocalStorage 使用本地文件系统实现FileRepository接口
type LocalStorage struct {
	config LocalStorageConfig
}

// NewLocalStorage 创建LocalStorage的新实例
func NewLocalStorage(config LocalStorageConfig) (*LocalStorage, error) {
	// 如果目录不存在则创建
	err := os.MkdirAll(config.UploadDir, 0755)
	if err != nil {
		return nil, err
	}

	err = os.MkdirAll(config.TempChunkDir, 0755)
	if err != nil {
		return nil, err
	}

	return &LocalStorage{
		config: config,
	}, nil
}

// SaveFile 将文件存储到本地文件系统
func (s *LocalStorage) SaveFile(name string, content io.Reader, size int64) (*domain.File, error) {
	// 生成唯一的文件ID
	fileID := uuid.New().String()

	// 创建文件路径
	ext := filepath.Ext(name)
	fileName := fileID + ext
	filePath := filepath.Join(s.config.UploadDir, fileName)

	// 创建文件
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// 复制内容
	_, err = io.Copy(dst, content)
	if err != nil {
		return nil, err
	}

	// 创建并返回文件实体
	file := &domain.File{
		ID:       fileID,
		Name:     name,
		Size:     size,
		Path:     filePath,
		UploadAt: time.Now(),
	}

	return file, nil
}

// SaveChunk 存储分片上传中的一个文件分片
func (s *LocalStorage) SaveChunk(uploadID string, index int, content io.Reader, size int64) error {
	// 如果分片目录不存在则创建
	chunkDir := filepath.Join(s.config.TempChunkDir, uploadID)
	err := os.MkdirAll(chunkDir, 0755)
	if err != nil {
		return err
	}

	// 创建分片文件
	chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d", index))
	dst, err := os.Create(chunkPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	// 复制内容
	_, err = io.Copy(dst, content)
	if err != nil {
		return err
	}

	return nil
}

// MergeChunks 将分片合并成完整的文件
func (s *LocalStorage) MergeChunks(uploadID string, totalChunks int, fileName string) (*domain.File, error) {
	// 获取分片目录
	chunkDir := filepath.Join(s.config.TempChunkDir, uploadID)

	// 生成唯一的文件ID
	fileID := uuid.New().String()

	// 创建文件路径
	ext := filepath.Ext(fileName)
	finalFileName := fileID + ext
	finalPath := filepath.Join(s.config.UploadDir, finalFileName)

	// 创建最终文件
	dst, err := os.Create(finalPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	// 按正确顺序合并所有分片
	var totalSize int64 = 0
	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d", i))

		// 打开分片文件
		chunk, err := os.Open(chunkPath)
		if err != nil {
			return nil, err
		}

		// 获取分片大小
		info, err := chunk.Stat()
		if err != nil {
			chunk.Close()
			return nil, err
		}
		totalSize += info.Size()

		// 将分片复制到最终文件
		_, err = io.Copy(dst, chunk)
		chunk.Close()
		if err != nil {
			return nil, err
		}
	}

	// 清理分片
	os.RemoveAll(chunkDir)

	// 创建并返回文件实体
	file := &domain.File{
		ID:       fileID,
		Name:     fileName,
		Size:     totalSize,
		Path:     finalPath,
		UploadAt: time.Now(),
	}

	return file, nil
}

// GetFile 通过ID获取文件
func (s *LocalStorage) GetFile(id string) (*domain.File, error) {
	// 通常这里会涉及数据库查询
	// 为了演示目的简化实现

	// 在上传目录中搜索文件
	matches, err := filepath.Glob(filepath.Join(s.config.UploadDir, id+"*"))
	if err != nil {
		return nil, err
	}

	if len(matches) == 0 {
		return nil, fmt.Errorf("未找到文件")
	}

	filePath := matches[0]
	info, err := os.Stat(filePath)
	if err != nil {
		return nil, err
	}

	file := &domain.File{
		ID:       id,
		Name:     filepath.Base(filePath),
		Size:     info.Size(),
		Path:     filePath,
		UploadAt: info.ModTime(),
	}

	return file, nil
}

// DeleteFile 通过ID删除文件
func (s *LocalStorage) DeleteFile(id string) error {
	// 查找文件
	file, err := s.GetFile(id)
	if err != nil {
		return err
	}

	// 删除文件
	return os.Remove(file.Path)
}
