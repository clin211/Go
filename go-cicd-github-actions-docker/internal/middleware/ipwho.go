package middleware

import (
	"fmt"

	"github.com/clin211/go-cicd-github-actions-docker/internal/pkg/known"
	"github.com/gin-gonic/gin"
)

// RequestID 是一个 Gin 中间件，用来在每一个 HTTP 请求的 context, response 中注入 `X-Request-ID` 键值对.
func IpWho() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ips := ipwhos.NewWhoisLib(c)
		// reply, err := ips.GetHostIpDetail()
		// if err != nil {
		// 	fmt.Println("error: ", err)
		// }
		// ip := reply.IP
		ip := c.ClientIP()
		fmt.Println("middleware ip: ", ip)
		// 将 RequestID 保存在 gin.Context 中，方便后边程序使用
		c.Set(known.XIpWhoKey, ip)

		c.Next()
	}
}
