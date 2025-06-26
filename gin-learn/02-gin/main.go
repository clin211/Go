package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Article struct {
	title       string
	description string
	content     string
}

type User struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Age   int    `json:"age"`
}

type UserXML struct {
	Name  string `xml:"name"`
	Email string `xml:"email"`
	Age   int    `xml:"age"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("pages/*")

	r.GET("/user/:id", func(c *gin.Context) {
		id := c.Param("id") // 获取路径参数
		c.JSON(http.StatusOK, gin.H{
			"user_id": id,
		})
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "%v page", "home")
	})

	r.GET("/json1", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]interface{}{
			"status":  http.StatusOK,
			"message": "ok",
			"data":    "json1",
		})
	})

	r.GET("/json2", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "ok",
			"data":    "json2",
		})
	})

	r.GET("/json3", func(c *gin.Context) {
		article := &Article{
			title:       "gin",
			content:     "content",
			description: "description",
		}
		fmt.Println(article)
		c.JSON(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"message": "ok",
			"data":    article,
		})
	})

	//返回xml数据
	r.GET("/xml", func(c *gin.Context) {
		c.XML(http.StatusOK, gin.H{
			"status":  http.StatusOK,
			"success": true,
			"message": "返回xml",
		})
	})

	r.GET("/news", func(c *gin.Context) {
		c.HTML(http.StatusOK, "news.html", gin.H{
			"title":   "this is news",
			"content": "gin html 模板渲染",
		})
	})
	r.GET("/goods", func(c *gin.Context) {
		c.HTML(http.StatusOK, "goods.html", gin.H{
			"title":   "this is news",
			"content": "iPad mini Mega power. Mini sized.",
		})
	})

	r.GET("/users", func(c *gin.Context) {
		users := c.QueryMap("user") // 解析 user[name]=xxx&user[age]=xxx
		c.JSON(http.StatusOK, gin.H{
			"users": users,
		})
	})

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"password": password,
		})
	})
	r.POST("/register", func(c *gin.Context) {
		username := c.DefaultPostForm("username", "guest")
		role := c.DefaultPostForm("role", "user")

		c.JSON(http.StatusOK, gin.H{
			"username": username,
			"role":     role,
		})
	})

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upload file"})
			return
		}

		// 保存文件
		dst := "./uploads/" + file.Filename
		c.SaveUploadedFile(file, dst)

		c.JSON(http.StatusOK, gin.H{
			"filename": file.Filename,
			"size":     file.Size,
		})
	})

	r.POST("/multi-upload", func(c *gin.Context) {
		form, err := c.MultipartForm()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
			return
		}

		files := form.File["files"]
		var filenames []string

		for _, file := range files {
			dst := "./uploads/" + file.Filename
			c.SaveUploadedFile(file, dst)
			filenames = append(filenames, file.Filename)
		}

		c.JSON(http.StatusOK, gin.H{"uploaded_files": filenames})
	})

	r.POST("/json", func(c *gin.Context) {
		var user User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User received", "user": user})
	})

	r.POST("/xml", func(c *gin.Context) {
		var user UserXML
		if err := c.ShouldBindXML(&user); err != nil {
			c.XML(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.XML(http.StatusOK, gin.H{"message": "User received", "user": user})
	})

	r.POST("/bind", func(c *gin.Context) {
		var user User
		// 根据 Content-Type 类型自动绑定
		if err := c.ShouldBind(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "User received", "user": user})
	})

	r.Run(":8080")
}
