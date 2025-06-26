package domain

import (
	"time"
)

// File 表示一个已完成上传的文件
type File struct {
	ID       string
	Name     string
	Size     int64
	Path     string
	MimeType string
	UploadAt time.Time
}

// Chunk 表示分片上传中的一个文件片段
type Chunk struct {
	UploadID string
	Index    int
	Total    int
	TempPath string
	Size     int64
}

// Folder 表示包含文件的目录结构
type Folder struct {
	Name       string
	Path       string
	Children   []*File
	SubFolders []*Folder
}
