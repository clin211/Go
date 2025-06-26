package domain

import (
	"io"
)

// FileRepository 定义文件存储所需的操作接口
type FileRepository interface {
	SaveFile(name string, content io.Reader, size int64) (*File, error)
	SaveChunk(uploadID string, index int, content io.Reader, size int64) error
	MergeChunks(uploadID string, totalChunks int, fileName string) (*File, error)
	GetFile(id string) (*File, error)
	DeleteFile(id string) error
}

// UploadDomainService 处理文件上传的领域逻辑
type UploadDomainService interface {
	ValidateChunk(chunk *Chunk) error
	MergeChunks(uploadID string, total int, fileName string) (*File, error)
	ValidateFile(name string, size int64, mimeType string) error
}

// UploadDomainServiceImpl 实现UploadDomainService接口
type UploadDomainServiceImpl struct {
	repo FileRepository
}

// NewUploadDomainService 创建UploadDomainService的新实例
func NewUploadDomainService(repo FileRepository) UploadDomainService {
	return &UploadDomainServiceImpl{
		repo: repo,
	}
}

// ValidateChunk 验证分片是否有效可处理
func (s *UploadDomainServiceImpl) ValidateChunk(chunk *Chunk) error {
	// 在此添加分片验证逻辑
	// 例如，检查索引是否有效，总数是否一致等
	return nil
}

// MergeChunks 处理合并文件分片的逻辑
func (s *UploadDomainServiceImpl) MergeChunks(uploadID string, total int, fileName string) (*File, error) {
	return s.repo.MergeChunks(uploadID, total, fileName)
}

// ValidateFile 上传前验证文件
func (s *UploadDomainServiceImpl) ValidateFile(name string, size int64, mimeType string) error {
	// 在此添加文件验证逻辑
	// 例如，检查文件大小限制，允许的扩展名等
	return nil
}
