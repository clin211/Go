package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// StartTrace 生成traceID和spanID
func StartTrace() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("traceid")
		parentSpanID := c.Request.Header.Get("spanid")
		spanID := generateSpanID(c.Request.RemoteAddr)

		if traceID == "" { // 如果traceID为空，证明是链路的发端，把它设置成此次的spanID
			traceID = spanID // trace标识整个请求的链路, span则标识链路中的不同服务
		}

		c.Set("traceid", traceID)
		c.Set("spanid", spanID)
		c.Set("pspanid", parentSpanID)

		// 添加到响应头，便于追踪
		c.Header("X-Trace-ID", traceID)
		c.Header("X-Span-ID", spanID)

		c.Next()
	}
}

// generateSpanID 生成span标识符
func generateSpanID(prefix string) string {
	// 生成8字节的随机数据
	b := make([]byte, 8)
	_, err := rand.Read(b)
	if err != nil {
		// 如果随机数生成失败，使用时间戳作为备选
		return fmt.Sprintf("%s-%x", prefix, time.Now().UnixNano())
	}

	return fmt.Sprintf("%s-%s", prefix, hex.EncodeToString(b))
}
