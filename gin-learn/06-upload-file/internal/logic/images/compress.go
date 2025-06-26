package images

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/clin211/gin-learn/06-upload-file/internal/config"
	"github.com/clin211/gin-learn/06-upload-file/internal/pkg/imageutil"
	"github.com/clin211/gin-learn/06-upload-file/internal/pkg/pathutil"
	"github.com/clin211/gin-learn/06-upload-file/internal/pkg/storage"
)

// CompressImageLogic 图片压缩业务逻辑结构体
type CompressImageLogic struct {
	Storage storage.Storage // 存储接口
}

// NewCompressImageLogic 创建图片压缩逻辑处理器实例
func NewCompressImageLogic(store storage.Storage) *CompressImageLogic {
	return &CompressImageLogic{Storage: store}
}

// CompressRequest 压缩请求参数
type CompressRequest struct {
	Quality      int    `json:"quality" form:"quality"`             // 压缩质量 (1-100)
	MaxWidth     uint   `json:"max_width" form:"max_width"`         // 最大宽度
	MaxHeight    uint   `json:"max_height" form:"max_height"`       // 最大高度
	Format       string `json:"format" form:"format"`               // 输出格式
	PreserveName bool   `json:"preserve_name" form:"preserve_name"` // 是否保留原始文件名
}

// CompressResult 压缩结果
type CompressResult struct {
	Filename       string  `json:"filename"`        // 文件名
	ObjectKey      string  `json:"object_key"`      // 存储的对象键
	OriginalSize   int64   `json:"original_size"`   // 原始大小(字节)
	CompressedSize int64   `json:"compressed_size"` // 压缩后大小(字节)
	Ratio          float64 `json:"ratio"`           // 压缩比率
	Width          int     `json:"width"`           // 宽度
	Height         int     `json:"height"`          // 高度
	Format         string  `json:"format"`          // 格式
	QualityScore   int     `json:"quality_score"`   // 质量评分
	PreviewURL     string  `json:"preview_url"`     // 预览URL
	Error          string  `json:"error,omitempty"` // 错误信息(如果有)
}

// CompressImage 压缩单张图片
func (c *CompressImageLogic) CompressImage(file *multipart.FileHeader, req CompressRequest) (*CompressResult, error) {
	// 校验是否为图片
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !isImageExtension(ext) {
		return nil, fmt.Errorf("不支持的图片格式: %s", ext)
	}

	// 设置压缩选项
	options := imageutil.CompressOptions{
		Quality:      req.Quality,
		MaxWidth:     req.MaxWidth,
		MaxHeight:    req.MaxHeight,
		Format:       req.Format,
		PreserveName: req.PreserveName,
	}

	// 执行压缩
	compressedData, imgInfo, err := imageutil.CompressImage(file, options)
	if err != nil {
		return nil, fmt.Errorf("压缩失败: %w", err)
	}

	// 生成唯一文件名
	filename := file.Filename
	downloadFilename := filename
	if !req.PreserveName {
		ext := "." + imgInfo.Format
		baseName := strings.TrimSuffix(filepath.Base(file.Filename), filepath.Ext(file.Filename))
		downloadFilename = baseName + "_compressed" + ext
	}

	// 生成存储路径
	objectKey := pathutil.GenerateFilePath(filename)

	// 保存压缩后的图片
	if err := c.Storage.SaveFromBytes(compressedData, objectKey); err != nil {
		return nil, fmt.Errorf("保存压缩图片失败: %w", err)
	}

	// 获取图片访问URL，指定下载时使用的文件名
	previewURL, err := c.Storage.GetURLWithFilename(objectKey, downloadFilename)
	if err != nil {
		previewURL = "/static/" + objectKey + "?filename=" + downloadFilename // 兜底操作
	}

	// 构建返回结果
	result := &CompressResult{
		Filename:       downloadFilename,
		ObjectKey:      objectKey,
		OriginalSize:   imgInfo.OriginalSize,
		CompressedSize: imgInfo.CompressedSize,
		Ratio:          imgInfo.Ratio,
		Width:          imgInfo.Width,
		Height:         imgInfo.Height,
		Format:         imgInfo.Format,
		QualityScore:   imgInfo.QualityScore,
		PreviewURL:     previewURL,
	}

	return result, nil
}

// BatchCompressImages 批量压缩图片
func (c *CompressImageLogic) BatchCompressImages(files []*multipart.FileHeader, req CompressRequest) ([]*CompressResult, error) {
	// 检查文件数量限制
	if len(files) > 100 {
		return nil, fmt.Errorf("一次最多处理100张图片，当前提交: %d", len(files))
	}

	results := make([]*CompressResult, 0, len(files))

	for _, file := range files {
		result, err := c.CompressImage(file, req)
		if err != nil {
			// 对于批量处理，单个失败不中断整体流程，而是记录错误
			results = append(results, &CompressResult{
				Filename: file.Filename,
				Error:    err.Error(),
			})
			continue
		}
		results = append(results, result)
	}

	return results, nil
}

// 判断文件扩展名是否是支持的图片格式
func isImageExtension(ext string) bool {
	// 从配置中获取允许的图片扩展名列表
	allowedExts := config.Local.AllowedExtensions

	// 确保扩展名格式一致（带点号）
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	// 转为小写进行比较
	ext = strings.ToLower(ext)

	for _, allowed := range allowedExts {
		if strings.ToLower(allowed) == ext {
			return true
		}
	}

	return false
}
