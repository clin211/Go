package main

import (
	"encoding/xml"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Userinfo struct {
	Username string `json:"username" form:"username" binding:"required"`
	Email    string `json:"email" form:"email" binding:"required"`
}

type Article struct {
	Title   string `json:"title" xml:"title" binding:"required"`
	Content string `json:"content" xml:"content" binding:"required"`
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("pages/*")

	// GET 请求传值
	r.GET("/", func(c *gin.Context) {
		// Query() 获取请求时携带的参数 DefaultQuery()如果没有获取到携带的参数则第二个参数为默认值
		username := c.Query("username")
		email := c.Query("email")
		page := c.DefaultQuery("page", "1")

		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "ok",
			"data": map[string]interface{}{
				"username": username,
				"email":    email,
				"page":     page,
			},
		})
	})

	r.GET("/user", func(c *gin.Context) {
		c.HTML(http.StatusOK, "user.html", gin.H{})
	})

	r.POST("/add-user", func(c *gin.Context) {
		// 获取表单传递的数据
		username := c.PostForm("username")
		email := c.PostForm("email")
		hobby := c.DefaultPostForm("hobby", "react")
		fmt.Println(username, email)
		c.JSON(http.StatusOK, gin.H{
			"code":    http.StatusOK,
			"message": "ok",
			"data": map[string]interface{}{
				"username": username,
				"email":    email,
				"hobby":    hobby,
			},
		})
	})

	// 获取GET POST传递的数据绑定到结构体上
	r.GET("/get-user", func(c *gin.Context) {
		user := &Userinfo{}
		err := c.ShouldBind(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, user)
		}
	})

	r.POST("/add", func(c *gin.Context) {
		user := &Userinfo{}

		if err := c.ShouldBind(&user); err == nil {
			c.JSON(http.StatusOK, user)
		} else {
			c.JSON(http.StatusOK, gin.H{
				"err": err.Error(),
			})
		}
	})

	// 获取post xml数据
	r.POST("/xml", func(c *gin.Context) {
		xmlSliceData, _ := c.GetRawData() // 获取c.Request.Body

		article := &Article{}

		if err := xml.Unmarshal(xmlSliceData, &article); err == nil {
			c.JSON(http.StatusOK, article)
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
		}
	})

	// 动态路由传值
	r.GET("/detail/:id", func(c *gin.Context) {
		id := c.Param("id")

		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	})
	r.Run(":9090")
}
