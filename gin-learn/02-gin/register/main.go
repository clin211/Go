package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 地址信息
type Address struct {
	City    string `json:"city" xml:"city" form:"city" binding:"required"`
	ZipCode string `json:"zipcode" xml:"zipcode" form:"zipcode" binding:"required,len=6"`
}

// 注册请求
type RegisterRequest struct {
	Username        string  `json:"username" xml:"username" form:"username" binding:"required,min=3,max=20"`
	Email           string  `json:"email" xml:"email" form:"email" binding:"required,email"`
	Phone           string  `json:"phone" xml:"phone" form:"phone" binding:"required,e164"`
	Password        string  `json:"password" xml:"password" form:"password" binding:"required,min=8,max=32"`
	ConfirmPassword string  `json:"confirm_password" xml:"confirm_password" form:"confirm_password" binding:"required,eqfield=Password"`
	Age             int     `json:"age" xml:"age" form:"age" binding:"required,gte=18,lte=65"`
	Address         Address `json:"address" xml:"address" binding:"required"`
}

// 响应结构体
type RegisterResponse struct {
	Message   string      `json:"message"`
	User      interface{} `json:"user"`
	RequestID string      `json:"request_id"`
}

// 处理注册请求
func RegisterHandler(c *gin.Context) {
	var req RegisterRequest

	// 解析 JSON、XML 或 Form
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 获取 Header
	clientVersion := c.GetHeader("X-Client-Version")
	if clientVersion == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "X-Client-Version header is required"})
		return
	}

	c.Header("X-Request-ID", "1234567890")
	// 返回 JSON 响应
	c.JSON(http.StatusOK, RegisterResponse{
		Message: "注册成功",
		User:    req,
	})
}

func main() {
	r := gin.Default()
	r.POST("/register", RegisterHandler)
	r.Run(":8080")
}
