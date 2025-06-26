package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/clin211/go-cicd-github-actions-docker/internal/api"
	"github.com/clin211/go-cicd-github-actions-docker/internal/config"
	"github.com/clin211/go-cicd-github-actions-docker/internal/middleware"
	"github.com/clin211/go-cicd-github-actions-docker/internal/pkg/db"
	"github.com/clin211/go-cicd-github-actions-docker/internal/pkg/log"
)

// @title 用于演示 CI/CD 的 API
// @version 1.0
// @description 使用 Github Actions 和 Docker 进行 CI/CD 的演示
// @termsOfService http://swagger.io/terms/

// @contact.name clin
// @contact.email 767425412lin@gmail.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:18080
// @BasePath /api/v1
func main() {
	// 定义命令行参数
	var configFile string
	flag.StringVar(&configFile, "f", "", "指定配置文件路径 (默认: configs/config.yaml)")
	flag.Parse()

	// 加载配置
	cfg, err := config.Init(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置文件失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log.Init(&log.Options{
		DisableCaller:     cfg.Log.DisableCaller,
		DisableStacktrace: cfg.Log.DisableStacktrace,
		Level:             cfg.Log.Level,
		Format:            cfg.Log.Format,
		OutputPaths:       cfg.Log.OutputPaths,
	})

	// 初始化数据库连接
	if err := db.InitAll(log.GetLogger()); err != nil {
		log.Fatalw("Failed to initialize database connections", "error", err)
	}
	defer func() {
		if err := db.CloseAll(); err != nil {
			log.Errorw("Failed to close database connections", "error", err)
		}
	}()

	// 设置Gin
	r := gin.New()

	// 中间件
	mws := []gin.HandlerFunc{
		middleware.NoCache,
		middleware.Cors,
		middleware.Secure,
		middleware.RequestID(),
		middleware.IpWho(),
		middleware.LogAccess(),
		middleware.Recovery(),
	}

	r.Use(mws...)

	// 注册API路由
	api.RegisterRoutes(r)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:    cfg.Server.Address,
		Handler: r,
	}

	// 在goroutine中启动服务器
	go func() {
		log.Infow("Starting server", "address", cfg.Server.Address)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalw("Failed to start server", "error", err)
		}
	}()

	// 创建一个通道接收信号
	quit := make(chan os.Signal, 1)
	// kill (无参数) 默认发送 syscall.SIGTERM
	// kill -2 发送 syscall.SIGINT
	// kill -9 发送 syscall.SIGKILL，但无法被捕获，所以不需要添加
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Infow("Shutting down server...")

	// 创建一个超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅地关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalw("Server forced to shutdown: %v", err)
	}

	log.Infow("Server exiting")
}
