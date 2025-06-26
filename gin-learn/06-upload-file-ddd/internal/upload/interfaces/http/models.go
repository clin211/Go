package http

import "06-upload-file/internal/upload/domain"

// FileResponse 表示Swagger文档的文件响应
type FileResponse struct {
	Name     string `json:"name" example:"example.jpg"`
	Path     string `json:"path" example:"/uploads/2023/04/example.jpg"`
	Size     int64  `json:"size" example:"1024000"`
	URL      string `json:"url" example:"http://localhost:8080/uploads/2023/04/example.jpg"`
	MimeType string `json:"mimeType" example:"image/jpeg"`
}

// SingleUploadResponse 表示单文件上传响应
type SingleUploadResponse struct {
	File FileResponse `json:"file"`
}

// MultipleUploadResponse 表示多文件上传响应
type MultipleUploadResponse struct {
	Files []FileResponse `json:"files"`
	Total int            `json:"total" example:"3"`
}

// FolderUploadResponse 表示文件夹上传响应
type FolderUploadResponse struct {
	FolderName string         `json:"folderName" example:"my-project"`
	FolderPath string         `json:"folderPath" example:"/uploads/folders/my-project"`
	Files      []FileResponse `json:"files"`
	Total      int            `json:"total" example:"3"`
}

// ChunkInitResponse 表示分片上传初始化响应
type ChunkInitResponse struct {
	UploadID string `json:"upload_id" example:"c0e0d7b1-9eb1-4c70-a9c5-69f7ca214285"`
}

// ChunkUploadResponse 表示分片上传响应
type ChunkUploadResponse struct {
	UploadID string `json:"upload_id" example:"c0e0d7b1-9eb1-4c70-a9c5-69f7ca214285"`
	Index    int    `json:"index" example:"3"`
	Total    int    `json:"total" example:"20"`
}

// ChunkMergeResponse 表示分片合并响应
type ChunkMergeResponse struct {
	File FileResponse `json:"file"`
}

// ChunkInitRequest 表示分片上传初始化请求
type ChunkInitRequest struct {
	FileName    string `json:"fileName" example:"large-video.mp4"`
	FileSize    int64  `json:"fileSize" example:"104857600"`
	TotalChunks int    `json:"totalChunks" example:"20"`
	MimeType    string `json:"mimeType" example:"video/mp4"`
}

// ChunkMergeRequest 表示分片合并请求
type ChunkMergeRequest struct {
	UploadID    string `json:"upload_id" example:"c0e0d7b1-9eb1-4c70-a9c5-69f7ca214285"`
	FileName    string `json:"file_name" example:"large-video.mp4"`
	TotalChunks int    `json:"total_chunks" example:"20"`
	ContentType string `json:"content_type" example:"video/mp4"`
}

// FromDomainFile 将domain.File转换为FileResponse
func FromDomainFile(file *domain.File) FileResponse {
	return FileResponse{
		Name:     file.Name,
		Path:     file.Path,
		Size:     file.Size,
		URL:      "/uploads/" + file.Path,
		MimeType: file.MimeType,
	}
}

// FromDomainFiles 将domain.File切片转换为FileResponse切片
func FromDomainFiles(files []*domain.File) []FileResponse {
	result := make([]FileResponse, len(files))
	for i, file := range files {
		result[i] = FromDomainFile(file)
	}
	return result
}
