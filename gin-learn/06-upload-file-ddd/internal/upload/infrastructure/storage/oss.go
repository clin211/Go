package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"06-upload-file/internal/shared/logger"
	"06-upload-file/internal/upload/domain"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/google/uuid"
)

// OSSStorageConfig 包含阿里云OSS对象存储的配置
type OSSStorageConfig struct {
	BucketName   string
	Endpoint     string
	AccessKey    string
	SecretKey    string
	URLPrefix    string // URL前缀，用于生成可访问的URL
	TempChunkDir string // 临时存储分片的目录
}

// OSSStorage 使用阿里云OSS实现FileRepository接口
type OSSStorage struct {
	config    OSSStorageConfig
	client    *oss.Client
	bucket    *oss.Bucket
	urlPrefix string
}

// NewOSSStorage 创建OSSStorage的新实例
func NewOSSStorage(config OSSStorageConfig) (*OSSStorage, error) {
	// 创建OSS客户端
	client, err := oss.New(config.Endpoint, config.AccessKey, config.SecretKey)
	if err != nil {
		return nil, fmt.Errorf("创建OSS客户端失败: %w", err)
	}

	// 获取存储空间
	bucket, err := client.Bucket(config.BucketName)
	if err != nil {
		return nil, fmt.Errorf("获取OSS存储空间失败: %w", err)
	}

	// 确保分片临时目录存在
	err = os.MkdirAll(config.TempChunkDir, 0755)
	if err != nil {
		return nil, fmt.Errorf("创建临时分片目录失败: %w", err)
	}

	// 处理URL前缀
	urlPrefix := config.URLPrefix
	if urlPrefix == "" {
		urlPrefix = fmt.Sprintf("https://%s.%s/", config.BucketName, config.Endpoint)
	}

	return &OSSStorage{
		config:    config,
		client:    client,
		bucket:    bucket,
		urlPrefix: urlPrefix,
	}, nil
}

// SaveFile 将文件上传到OSS
func (s *OSSStorage) SaveFile(name string, content io.Reader, size int64) (*domain.File, error) {
	// 生成唯一的文件ID和路径
	fileID := uuid.New().String()
	objectKey := fmt.Sprintf("%s/%s", time.Now().Format("2006/01/02"), fileID)

	// 获取文件扩展名
	ext := filepath.Ext(name)
	if ext != "" {
		objectKey = objectKey + ext
	}

	// 上传文件到OSS
	err := s.bucket.PutObject(objectKey, content)
	if err != nil {
		return nil, fmt.Errorf("上传文件到OSS失败: %w", err)
	}

	// 构建可访问的URL
	fileURL := s.urlPrefix + objectKey

	// 返回文件实体
	file := &domain.File{
		ID:       fileID,
		Name:     name,
		Size:     size,
		Path:     objectKey,
		MimeType: "", // OSS会自动检测或可以从请求中获取
		UploadAt: time.Now(),
	}

	logger.Info("文件已上传到OSS: %s, URL: %s", objectKey, fileURL)
	return file, nil
}

// SaveChunk 保存分片上传期间的一个文件分片
func (s *OSSStorage) SaveChunk(uploadID string, index int, content io.Reader, size int64) error {
	// 为分片创建一个临时文件
	chunkDir := filepath.Join(s.config.TempChunkDir, uploadID)
	err := os.MkdirAll(chunkDir, 0755)
	if err != nil {
		return fmt.Errorf("创建分片目录失败: %w", err)
	}

	chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d", index))
	chunkFile, err := os.Create(chunkPath)
	if err != nil {
		return fmt.Errorf("创建分片文件失败: %w", err)
	}
	defer chunkFile.Close()

	// 将内容复制到临时文件
	_, err = io.Copy(chunkFile, content)
	if err != nil {
		return fmt.Errorf("写入分片数据失败: %w", err)
	}

	logger.Info("上传ID %s 的分片 %d 已保存", uploadID, index)
	return nil
}

// MergeChunks 通过合并所有分片并上传为单个对象来完成分片上传
func (s *OSSStorage) MergeChunks(uploadID string, totalChunks int, fileName string) (*domain.File, error) {
	chunkDir := filepath.Join(s.config.TempChunkDir, uploadID)

	// 生成唯一的文件ID和路径
	fileID := uuid.New().String()
	objectKey := fmt.Sprintf("%s/%s", time.Now().Format("2006/01/02"), fileID)

	// 获取文件扩展名
	ext := filepath.Ext(fileName)
	if ext != "" {
		objectKey = objectKey + ext
	}

	// 准备分片上传
	imur, err := s.bucket.InitiateMultipartUpload(objectKey)
	if err != nil {
		return nil, fmt.Errorf("初始化分片上传失败: %w", err)
	}

	// 上传各个分片
	var parts []oss.UploadPart
	var totalSize int64

	for i := 0; i < totalChunks; i++ {
		chunkPath := filepath.Join(chunkDir, fmt.Sprintf("%d", i))

		// 检查分片是否存在
		if _, err := os.Stat(chunkPath); os.IsNotExist(err) {
			// 取消分片上传
			s.bucket.AbortMultipartUpload(imur)
			return nil, fmt.Errorf("上传ID %s 的分片 %d 未找到", uploadID, i)
		}

		// 获取分片文件大小
		chunkInfo, err := os.Stat(chunkPath)
		if err != nil {
			s.bucket.AbortMultipartUpload(imur)
			return nil, fmt.Errorf("获取分片信息失败: %w", err)
		}
		totalSize += chunkInfo.Size()

		// 打开分片文件
		chunk, err := os.Open(chunkPath)
		if err != nil {
			s.bucket.AbortMultipartUpload(imur)
			return nil, fmt.Errorf("打开分片文件失败: %w", err)
		}

		// 上传分片
		part, err := s.bucket.UploadPart(imur, chunk, chunkInfo.Size(), i+1)
		chunk.Close()

		if err != nil {
			s.bucket.AbortMultipartUpload(imur)
			return nil, fmt.Errorf("上传分片 %d 失败: %w", i, err)
		}

		parts = append(parts, part)
	}

	// 完成分片上传
	_, err = s.bucket.CompleteMultipartUpload(imur, parts)
	if err != nil {
		return nil, fmt.Errorf("完成分片上传失败: %w", err)
	}

	// 清理临时分片文件
	os.RemoveAll(chunkDir)

	// 构建可访问的URL
	fileURL := s.urlPrefix + objectKey

	// 返回文件实体
	file := &domain.File{
		ID:       fileID,
		Name:     fileName,
		Size:     totalSize,
		Path:     objectKey,
		MimeType: "", // 获取MIME类型，通常根据扩展名或从请求中获取
		UploadAt: time.Now(),
	}

	logger.Info("分片已合并并上传到OSS: %s, URL: %s", objectKey, fileURL)
	return file, nil
}

// GetFile 获取文件元数据
func (s *OSSStorage) GetFile(id string) (*domain.File, error) {
	// 在实际应用中，你可能需要从数据库获取文件元数据
	// 这里简化实现，仅用于演示架构
	logger.Info("获取文件 %s 的元数据", id)
	return &domain.File{
		ID:   id,
		Name: id,
		Path: id,
	}, nil
}

// DeleteFile 从OSS中删除文件
func (s *OSSStorage) DeleteFile(id string) error {
	// 在实际应用中，你需要根据ID查找文件路径
	// 这里简化实现，假设ID就是对象键
	err := s.bucket.DeleteObject(id)
	if err != nil {
		return fmt.Errorf("从OSS删除文件失败: %w", err)
	}

	logger.Info("文件已从OSS删除: %s", id)
	return nil
}
