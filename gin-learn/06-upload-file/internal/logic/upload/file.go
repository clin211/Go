package upload

import (
	"mime/multipart"

	"github.com/clin211/gin-learn/06-upload-file/internal/pkg/pathutil"
	"github.com/clin211/gin-learn/06-upload-file/internal/pkg/storage"
	"github.com/clin211/gin-learn/06-upload-file/internal/pkg/validator"
)

// FileUploadLogic 文件上传业务逻辑结构体
// 通过组合Storage接口实现对不同存储方式的支持
type FileUploadLogic struct {
	Storage storage.Storage // 存储接口，支持本地存储、对象存储等多种方式
}

// NewFileUploadLogic 创建文件上传逻辑处理器实例
// 参数:
//   - store: 实现了Storage接口的存储实例
//
// 返回值:
//   - *FileUploadLogic: 初始化后的上传逻辑处理器
func NewFileUploadLogic(store storage.Storage) *FileUploadLogic {
	return &FileUploadLogic{Storage: store}
}

// Upload 处理文件上传的核心方法
// 完整的处理流程包括：验证文件 -> 生成存储路径 -> 保存文件 -> 返回文件标识
// 参数:
//   - file: 用户上传的文件信息
//
// 返回值:
//   - string: 文件的唯一标识符/存储路径
//   - error: 处理过程中可能发生的错误
func (l *FileUploadLogic) Upload(file *multipart.FileHeader) (string, error) {
	// 校验文件合法性
	// 验证文件类型和大小是否符合系统要求
	if err := validator.ValidateFile(file); err != nil {
		return "", err
	}

	// 生成存储路径
	// 基于日期和UUID创建唯一的文件路径，避免文件名冲突
	objectKey := pathutil.GenerateFilePath(file.Filename)

	// 存储文件
	// 将文件保存到配置的存储系统中（可能是本地文件系统、云存储等）
	if err := l.Storage.Save(file, objectKey); err != nil {
		return "", err
	}

	// 返回文件的唯一标识符
	// 该标识符可用于后续获取文件URL或执行其他操作
	return objectKey, nil
}
