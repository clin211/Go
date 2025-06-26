package main

import (
	"embed"
	"html/template"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

var static embed.FS

func main() {
	router := gin.Default()

	// 注册函数
	router.SetFuncMap(template.FuncMap{
		"formatDate": func(t time.Time) string {
			return t.Format("2006-01-02 15:04:05")
		},
		"add":      lo.Sum[int],         // 求和
		"max":      lo.Max[int],         // 最大值
		"min":      lo.Min[int],         // 最小值
		"contains": lo.Contains[string], // 判断是否包含
	})

	// router.LoadHTMLGlob("templates/*")
	router.LoadHTMLGlob("templates/**/*") // 加载多级目录模板
	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		// 文件路径是基于加载模板文件下的路径
		c.HTML(200, "index.html", gin.H{
			"title":     "Gin Template Demo",
			"username":  "长林啊",
			"isAdmin":   false,
			"isLoginIn": true,
			"skills":    []string{"Go", "Gin", "MySQL", "Redis"},
			"date":      time.Now(),
			"content":   "<script>alert('XSS')</script>",
			"numbers":   []int{3, 5, 7, 2},
			"word":      "hello",
			"words":     []string{"hello", "world", "gin"},
		})
	})

	// 路由
	router.GET("/detail", func(c *gin.Context) {
		c.HTML(http.StatusOK, "detail.html", gin.H{
			"Title": "Gin 模板继承",
		})
	})

	router.GET("/products", func(c *gin.Context) {
		c.HTML(http.StatusOK, "products.html", nil)
	})

	router.Run(":8080")
}
