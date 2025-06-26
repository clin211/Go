package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/clin211/go-cicd-github-actions-docker/internal/pkg/log"
	"github.com/gin-gonic/gin"
)

// bodyLogWriter 用于捕获响应体
type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 实现ResponseWriter接口，同时捕获响应数据
func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LogAccess 返回访问日志中间件
// 记录请求的开始和结束，包含请求体和响应体内容
func LogAccess() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 保存请求体
		reqBody, _ := io.ReadAll(c.Request.Body)
		c.Request.Body = io.NopCloser(bytes.NewReader(reqBody))

		// 记录开始时间
		start := time.Now()

		// 包装响应写入器以捕获响应体
		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}
		c.Writer = blw

		// 记录请求开始
		accessLog(c, "access_start", time.Since(start), reqBody, nil)

		// 请求结束时记录完整信息
		defer func() {
			accessLog(c, "access_end", time.Since(start), reqBody, blw.body.String())
		}()

		// 继续处理请求
		c.Next()
	}
}

// accessLog 记录访问日志的详细信息
func accessLog(c *gin.Context, accessType string, dur time.Duration, body []byte, dataOut interface{}) {
	req := c.Request

	// 获取请求相关信息
	bodyStr := string(body)
	query := req.URL.RawQuery
	path := req.URL.Path
	token := c.Request.Header.Get("Authorization")

	// 记录结构化日志
	log.C(c).Infow("AccessLog",
		"type", accessType,
		"ip", c.ClientIP(),
		"token", token,
		"method", req.Method,
		"path", path,
		"query", query,
		"body", bodyStr,
		"output", dataOut,
		"time(ms)", int64(dur/time.Millisecond))
}
