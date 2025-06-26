package files

import (
	"net/http"

	"github.com/clin211/gin-learn/06-upload-file/internal/config"
	"github.com/clin211/gin-learn/06-upload-file/internal/logic/upload"
	"github.com/clin211/gin-learn/06-upload-file/internal/pkg/storage"
	"github.com/gin-gonic/gin"
)

func (u *Files) File(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文件未上传"})
		return
	}

	dir := config.Local.UploadDir
	storage := storage.NewLocalStorage(dir)
	logic := upload.NewFileUploadLogic(storage)
	objectKey, err := logic.Upload(file)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "上传失败: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   "上传成功",
		"objectKey": objectKey,
	})
}
