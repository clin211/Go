package domain

// FileID 表示文件的唯一标识符
type FileID string

// FilePath 表示存储系统中的有效路径
type FilePath string

// UploadID 表示分片上传会话的唯一标识符
type UploadID string

// FileMetadata 包含关于文件的附加信息
type FileMetadata struct {
	ContentType string
	Extension   string
	IsPublic    bool
}
