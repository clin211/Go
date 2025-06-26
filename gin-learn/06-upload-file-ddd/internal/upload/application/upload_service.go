package application

import (
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"06-upload-file/internal/upload/domain"

	"github.com/google/uuid"
)

// UploadAppService 协调文件上传用例
type UploadAppService struct {
	domainService domain.UploadDomainService
	fileRepo      domain.FileRepository
}

// NewUploadAppService 创建一个新的上传应用服务
func NewUploadAppService(domainService domain.UploadDomainService, fileRepo domain.FileRepository) *UploadAppService {
	return &UploadAppService{
		domainService: domainService,
		fileRepo:      fileRepo,
	}
}

// UploadFile 处理单文件上传
func (s *UploadAppService) UploadFile(file *multipart.FileHeader) (*domain.File, error) {
	// 打开文件
	f, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// 验证文件
	err = s.domainService.ValidateFile(file.Filename, file.Size, file.Header.Get("Content-Type"))
	if err != nil {
		return nil, err
	}

	// 保存文件
	savedFile, err := s.fileRepo.SaveFile(file.Filename, f, file.Size)
	if err != nil {
		return nil, err
	}

	return savedFile, nil
}

// UploadMultipleFiles 处理批量文件上传
func (s *UploadAppService) UploadMultipleFiles(files []*multipart.FileHeader) ([]*domain.File, error) {
	result := make([]*domain.File, 0, len(files))

	for _, file := range files {
		uploadedFile, err := s.UploadFile(file)
		if err != nil {
			return nil, err
		}
		result = append(result, uploadedFile)
	}

	return result, nil
}

// UploadFolder 处理带有嵌套结构的文件夹上传
func (s *UploadAppService) UploadFolder(files map[string]*multipart.FileHeader) (*domain.Folder, error) {
	// 创建根文件夹
	rootFolder := &domain.Folder{
		Name:       "upload_" + time.Now().Format("20060102150405"),
		Path:       "",
		Children:   []*domain.File{},
		SubFolders: []*domain.Folder{},
	}

	// 处理每个文件
	for path, file := range files {
		// 打开文件
		f, err := file.Open()
		if err != nil {
			continue
		}

		// 提取目录结构
		dir, fileName := filepath.Split(path)

		// 保存文件
		savedFile, err := s.fileRepo.SaveFile(fileName, f, file.Size)
		if err != nil {
			f.Close()
			continue
		}
		f.Close()

		// 添加到文件夹结构
		if dir == "" {
			rootFolder.Children = append(rootFolder.Children, savedFile)
		} else {
			// 处理嵌套目录结构 - 目前简化实现
			rootFolder.Children = append(rootFolder.Children, savedFile)
		}
	}

	return rootFolder, nil
}

// UploadChunk 处理分片上传中的单个分片
func (s *UploadAppService) UploadChunk(uploadID string, index, total int, content io.Reader, size int64) error {
	chunk := &domain.Chunk{
		UploadID: uploadID,
		Index:    index,
		Total:    total,
		Size:     size,
	}

	// 验证分片
	err := s.domainService.ValidateChunk(chunk)
	if err != nil {
		return err
	}

	// 保存分片
	err = s.fileRepo.SaveChunk(uploadID, index, content, size)
	if err != nil {
		return err
	}

	return nil
}

// InitChunkUpload 初始化新的分片上传会话
func (s *UploadAppService) InitChunkUpload() (string, error) {
	// 生成唯一的上传ID
	uploadID := uuid.New().String()
	return uploadID, nil
}

// MergeChunks 将所有分片合并为最终文件
func (s *UploadAppService) MergeChunks(uploadID string, totalChunks int, fileName string) (*domain.File, error) {
	// 调用领域服务合并分片
	return s.domainService.MergeChunks(uploadID, totalChunks, fileName)
}
