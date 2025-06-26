package middleware

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorResponse 定义错误响应的结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

// Recovery 返回一个恢复中间件，用于处理可能的panic
// @Summary 恢复中间件
// @Description 捕获并处理API请求过程中的panic，防止服务器崩溃
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 检查是否已写入响应
				if c.Writer.Status() != http.StatusOK {
					c.Next()
					return
				}

				// 记录时间
				timestamp := time.Now().Format("2006/01/02 - 15:04:05")

				// 获取请求信息
				requestInfo := fmt.Sprintf("[Recovery] %s panic recovered:\n%s\n%s\n%s\n",
					timestamp, strings.Repeat("-", 100), fmt.Sprintf("%v", err), strings.Repeat("-", 100))

				// 获取堆栈信息
				stack := stack(3)

				// 输出到错误日志
				log.Printf("%s\n%s", requestInfo, stack)

				// 返回统一的错误响应
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Code:    http.StatusInternalServerError,
					Message: "服务器内部错误",
				})

				// 中止后续中间件的执行
				c.Abort()
			}
		}()

		// 继续处理请求
		c.Next()
	}
}

// stack 返回格式化的堆栈跟踪
func stack(skip int) []byte {
	buf := new(bytes.Buffer)

	// 跳过指定数量的调用栈
	var lines [][]byte
	var lastFile string

	for i := skip; ; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}

		// 打印文件名:行号
		fmt.Fprintf(buf, "%s:%d (0x%x)\n", file, line, pc)

		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}

			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}

		// 打印函数名
		fmt.Fprintf(buf, "\t%s: %s\n", function(pc), source(lines, line))
	}

	return buf.Bytes()
}

// function 返回函数名
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}

	name := []byte(fn.Name())

	// 解析函数名
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}

	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}

	name = bytes.Replace(name, centerDot, dot, -1)

	return name
}

// source 返回对应行的源代码
func source(lines [][]byte, n int) []byte {
	n-- // 调整为0索引

	if n < 0 || n >= len(lines) {
		return dunno
	}

	return bytes.TrimSpace(lines[n])
}

// LoggerConfig 定义日志中间件的配置
type LoggerConfig struct {
	// 输出写入器
	Output io.Writer
	// 是否禁用控制台颜色
	DisableColor bool
	// 是否禁用请求路径记录
	DisablePath bool
	// 自定义日志格式
	CustomFormat string
}

// DefaultLoggerConfig 返回默认的日志配置
func DefaultLoggerConfig() LoggerConfig {
	return LoggerConfig{
		Output:       os.Stdout,
		DisableColor: false,
		DisablePath:  false,
		CustomFormat: "",
	}
}

// Logger 返回一个日志中间件，使用默认配置
func Logger() gin.HandlerFunc {
	return LoggerWithConfig(DefaultLoggerConfig())
}

// LoggerWithConfig 返回一个基于给定配置的日志中间件
func LoggerWithConfig(config LoggerConfig) gin.HandlerFunc {
	// 设置默认输出
	if config.Output == nil {
		config.Output = os.Stdout
	}

	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		latency := end.Sub(start)

		// 获取状态码和客户端IP
		statusCode := c.Writer.Status()
		clientIP := c.ClientIP()
		method := c.Request.Method

		// 处理查询参数
		if raw != "" {
			path = path + "?" + raw
		}

		// 根据状态码选择颜色
		var statusColor, methodColor, resetColor string
		if !config.DisableColor {
			statusColor = colorForStatus(statusCode)
			methodColor = colorForMethod(method)
			resetColor = "\033[0m"
		}

		// 使用自定义格式或默认格式
		if config.CustomFormat != "" {
			// 替换自定义格式中的占位符
			logLine := strings.ReplaceAll(config.CustomFormat, "${status}", fmt.Sprintf("%d", statusCode))
			logLine = strings.ReplaceAll(logLine, "${method}", method)
			logLine = strings.ReplaceAll(logLine, "${path}", path)
			logLine = strings.ReplaceAll(logLine, "${latency}", latency.String())
			logLine = strings.ReplaceAll(logLine, "${ip}", clientIP)
			logLine = strings.ReplaceAll(logLine, "${time}", end.Format("2006/01/02 - 15:04:05"))

			fmt.Fprintln(config.Output, logLine)
		} else {
			// 使用默认格式
			fmt.Fprintf(config.Output, "%v |%s %3d %s| %13v | %15s |%s %-7s %s %#v\n",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, resetColor,
				latency,
				clientIP,
				methodColor, method, resetColor,
				path,
			)
		}
	}
}

// 辅助函数，根据状态码返回对应的控制台颜色
func colorForStatus(code int) string {
	switch {
	case code >= 200 && code < 300:
		return "\033[97;42m" // 绿色
	case code >= 300 && code < 400:
		return "\033[90;47m" // 灰色
	case code >= 400 && code < 500:
		return "\033[97;43m" // 黄色
	default:
		return "\033[97;41m" // 红色
	}
}

// 辅助函数，根据HTTP方法返回对应的控制台颜色
func colorForMethod(method string) string {
	switch method {
	case "GET":
		return "\033[97;44m" // 蓝色
	case "POST":
		return "\033[97;42m" // 绿色
	case "PUT":
		return "\033[97;43m" // 黄色
	case "DELETE":
		return "\033[97;41m" // 红色
	default:
		return "\033[97;45m" // 紫色
	}
}
