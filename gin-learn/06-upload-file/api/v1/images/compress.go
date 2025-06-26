package images

import (
	"net/http"

	"github.com/clin211/gin-learn/06-upload-file/internal/config"
	"github.com/clin211/gin-learn/06-upload-file/internal/logic/images"
	"github.com/clin211/gin-learn/06-upload-file/internal/pkg/storage"
	"github.com/gin-gonic/gin"
)

// Compress 处理单张图片压缩
func (i *ImageController) Compress(c *gin.Context) {
	// 获取上传的文件
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "未上传文件或文件无效",
			"error":   err.Error(),
		})
		return
	}

	// 解析请求参数
	var req images.CompressRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 设置默认值
	if req.Quality <= 0 || req.Quality > 100 {
		req.Quality = 85 // 默认压缩质量85%
	}

	// 创建存储和业务逻辑处理器
	dir := config.Local.UploadDir
	storage := storage.NewLocalStorage(dir)
	logic := images.NewCompressImageLogic(storage)

	// 执行压缩
	result, err := logic.CompressImage(file, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "压缩失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "压缩成功",
		"data":    result,
	})
}

// CompressMultiple 处理多张图片批量压缩
func (i *ImageController) CompressMultiple(c *gin.Context) {
	// 获取上传的文件
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "获取上传文件失败",
			"error":   err.Error(),
		})
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "未上传任何文件",
		})
		return
	}

	// 解析请求参数
	var req images.CompressRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "请求参数无效",
			"error":   err.Error(),
		})
		return
	}

	// 设置默认值
	if req.Quality <= 0 || req.Quality > 100 {
		req.Quality = 85 // 默认压缩质量85%
	}

	// 创建存储和业务逻辑处理器
	dir := config.Local.UploadDir
	storage := storage.NewLocalStorage(dir)
	logic := images.NewCompressImageLogic(storage)

	// 执行批量压缩
	results, err := logic.BatchCompressImages(files, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "批量压缩失败",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "批量压缩成功",
		"data":    results,
	})
}
