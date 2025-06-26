package http

import (
	"fmt"
	"mime/multipart"
	"strconv"

	"06-upload-file/internal/shared/logger"
	"06-upload-file/internal/shared/response"
	"06-upload-file/internal/upload/application"

	"github.com/gin-gonic/gin"
)

// UploadHandler 处理文件上传的HTTP请求
type UploadHandler struct {
	appService *application.UploadAppService
}

// NewUploadHandler 创建一个新的上传处理器
func NewUploadHandler(appService *application.UploadAppService) *UploadHandler {
	return &UploadHandler{
		appService: appService,
	}
}

// UploadSingle 处理单文件上传
// @Summary 单文件上传
// @Description 上传单个文件
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "要上传的文件"
// @Success 200 {object} response.Response{data=SingleUploadResponse} "上传成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /upload/single [post]
func (h *UploadHandler) UploadSingle(c *gin.Context) {
	// 从请求中获取文件
	file, err := c.FormFile("file")
	fmt.Println("file", file)
	if err != nil {
		logger.Error("获取文件失败: %v", err)
		response.BadRequest(c, "未上传任何文件")
		return
	}

	// 处理文件
	uploadedFile, err := h.appService.UploadFile(file)
	if err != nil {
		logger.Error("上传文件失败: %v", err)
		response.InternalError(c, "上传文件失败")
		return
	}

	// 返回结果
	response.Success(c, gin.H{
		"file": uploadedFile,
	})
}

// UploadMultiple 处理多文件上传
// @Summary 多文件上传
// @Description 一次上传多个文件
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "要上传的多个文件"
// @Success 200 {object} response.Response{data=MultipleUploadResponse} "上传成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /upload/multiple [post]
func (h *UploadHandler) UploadMultiple(c *gin.Context) {
	// 解析多部分表单
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		logger.Error("解析多部分表单失败: %v", err)
		response.BadRequest(c, "无效请求")
		return
	}

	// 获取文件
	form, err := c.MultipartForm()
	if err != nil {
		logger.Error("获取多部分表单失败: %v", err)
		response.BadRequest(c, "无效请求")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		response.BadRequest(c, "未上传任何文件")
		return
	}

	// 处理文件
	uploadedFiles, err := h.appService.UploadMultipleFiles(files)
	if err != nil {
		logger.Error("上传文件失败: %v", err)
		response.InternalError(c, "上传文件失败")
		return
	}

	// 返回结果
	response.Success(c, gin.H{
		"files": uploadedFiles,
	})
}

// UploadFolder 处理文件夹上传
// @Summary 文件夹上传
// @Description 上传完整文件夹，保留目录结构
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "要上传的文件夹中的文件"
// @Param folderName formData string true "文件夹名称"
// @Success 200 {object} response.Response{data=FolderUploadResponse} "上传成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /upload/folder [post]
func (h *UploadHandler) UploadFolder(c *gin.Context) {
	// 解析多部分表单
	if err := c.Request.ParseMultipartForm(32 << 20); err != nil {
		logger.Error("解析多部分表单失败: %v", err)
		response.BadRequest(c, "无效请求")
		return
	}

	// 获取文件
	form, err := c.MultipartForm()
	if err != nil {
		logger.Error("获取多部分表单失败: %v", err)
		response.BadRequest(c, "无效请求")
		return
	}

	// 将表单文件转换为带路径的映射
	fileMap := make(map[string]*multipart.FileHeader)
	for path, fileHeaders := range form.File {
		if len(fileHeaders) > 0 {
			fileMap[path] = fileHeaders[0]
		}
	}

	if len(fileMap) == 0 {
		response.BadRequest(c, "未上传任何文件")
		return
	}

	// 处理文件夹
	uploadedFolder, err := h.appService.UploadFolder(fileMap)
	if err != nil {
		logger.Error("上传文件夹失败: %v", err)
		response.InternalError(c, "上传文件夹失败")
		return
	}

	// 返回结果
	response.Success(c, gin.H{
		"folder": uploadedFolder,
	})
}

// InitChunkedUpload 初始化分片上传会话
// @Summary 初始化分片上传
// @Description 开始大文件分片上传前的初始化操作
// @Tags 上传
// @Accept json
// @Produce json
// @Param request body ChunkInitRequest true "初始化参数"
// @Success 200 {object} response.Response{data=ChunkInitResponse} "初始化成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /upload/chunk/init [post]
func (h *UploadHandler) InitChunkedUpload(c *gin.Context) {
	// 生成新的上传ID
	uploadID, err := h.appService.InitChunkUpload()
	if err != nil {
		logger.Error("初始化分片上传失败: %v", err)
		response.InternalError(c, "初始化上传失败")
		return
	}

	// 返回上传ID
	response.Success(c, gin.H{
		"upload_id": uploadID,
	})
}

// UploadChunk 处理单个文件分片的上传
// @Summary 上传文件分片
// @Description 上传单个文件分片
// @Tags 上传
// @Accept multipart/form-data
// @Produce json
// @Param upload_id formData string true "上传ID"
// @Param index formData integer true "分片索引，从0开始"
// @Param total formData integer true "总分片数"
// @Param chunk formData file true "分片数据"
// @Success 200 {object} response.Response{data=ChunkUploadResponse} "分片上传成功"
// @Failure 400 {object} response.Response "请求参数错误"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /upload/chunk [post]
func (h *UploadHandler) UploadChunk(c *gin.Context) {
	// 获取参数
	uploadID := c.PostForm("upload_id")
	if uploadID == "" {
		response.BadRequest(c, "缺少上传ID")
		return
	}

	indexStr := c.PostForm("index")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		response.BadRequest(c, "无效的索引")
		return
	}

	totalStr := c.PostForm("total")
	total, err := strconv.Atoi(totalStr)
	if err != nil {
		response.BadRequest(c, "无效的总数")
		return
	}

	// 获取分片文件
	chunk, err := c.FormFile("chunk")
	if err != nil {
		logger.Error("获取分片失败: %v", err)
		response.BadRequest(c, "未上传分片")
		return
	}

	// 打开分片
	chunkFile, err := chunk.Open()
	if err != nil {
		logger.Error("打开分片失败: %v", err)
		response.InternalError(c, "处理分片失败")
		return
	}
	defer chunkFile.Close()

	// 上传分片
	err = h.appService.UploadChunk(uploadID, index, total, chunkFile, chunk.Size)
	if err != nil {
		logger.Error("上传分片失败: %v", err)
		response.InternalError(c, "上传分片失败")
		return
	}

	// 返回成功
	response.Success(c, gin.H{
		"upload_id": uploadID,
		"index":     index,
		"total":     total,
	})
}

// MergeChunks 将所有已上传的分片合并为最终文件
// @Summary 合并文件分片
// @Description 所有分片上传完成后，合并为完整文件
// @Tags 上传
// @Accept json
// @Produce json
// @Param request body ChunkMergeRequest true "合并参数"
// @Success 200 {object} response.Response{data=ChunkMergeResponse} "合并成功"
// @Failure 400 {object} response.Response "请求参数错误或分片不完整"
// @Failure 500 {object} response.Response "服务器内部错误"
// @Router /upload/merge [post]
func (h *UploadHandler) MergeChunks(c *gin.Context) {
	// 解析请求体
	var req struct {
		UploadID    string `json:"upload_id"`
		FileName    string `json:"file_name"`
		TotalChunks int    `json:"total_chunks"`
		ContentType string `json:"content_type"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("解析合并请求失败: %v", err)
		response.BadRequest(c, "无效请求")
		return
	}

	if req.UploadID == "" {
		response.BadRequest(c, "缺少上传ID")
		return
	}

	// 合并分片
	fileInfo, err := h.appService.MergeChunks(req.UploadID, req.TotalChunks, req.FileName)
	if err != nil {
		logger.Error("合并分片失败: %v", err)
		response.InternalError(c, "合并分片失败")
		return
	}

	// 返回文件信息
	response.Success(c, gin.H{
		"file": fileInfo,
	})
}
