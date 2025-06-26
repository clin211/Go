package main

import (
	"fmt"
	"os"

	"06-upload-file/internal/shared/logger"
	"06-upload-file/internal/upload/application"
	"06-upload-file/internal/upload/domain"
	"06-upload-file/internal/upload/infrastructure/storage"
	uploadHttp "06-upload-file/internal/upload/interfaces/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "06-upload-file/docs" // 导入swagger生成的文档
)

// @title           文件上传系统 API
// @version         1.0
// @description     基于领域驱动设计的文件上传系统API，支持单文件上传、多文件上传、文件夹上传及大文件分片上传
// @termsOfService  http://swagger.io/terms/

// @contact.name   开发团队
// @contact.email  support@example.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// Config 应用程序配置结构
type Config struct {
	Server struct {
		Port    int    `mapstructure:"port"`
		Host    string `mapstructure:"host"`
		Timeout int    `mapstructure:"timeout"`
	} `mapstructure:"server"`

	Upload struct {
		Provider          string   `mapstructure:"provider"`
		MaxFileSize       int64    `mapstructure:"max_file_size"`
		MaxFiles          int      `mapstructure:"max_files"`
		AllowedExtensions []string `mapstructure:"allowed_extensions"`

		Local struct {
			UploadDir    string `mapstructure:"upload_dir"`
			TempChunkDir string `mapstructure:"temp_chunk_dir"`
		} `mapstructure:"local"`

		OSS struct {
			Bucket       string `mapstructure:"bucket"`
			Endpoint     string `mapstructure:"endpoint"`
			AccessKey    string `mapstructure:"access_key"`
			SecretKey    string `mapstructure:"secret_key"`
			URLPrefix    string `mapstructure:"url_prefix"`
			TempChunkDir string `mapstructure:"temp_chunk_dir"`
		} `mapstructure:"oss"`
	} `mapstructure:"upload"`

	Logging struct {
		Level string `mapstructure:"level"`
		File  string `mapstructure:"file"`
	} `mapstructure:"logging"`
}

func main() {
	// 加载配置
	config, err := loadConfig()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	initLogger(config)

	// 创建路由
	router := gin.Default()

	// 设置请求体最大大小
	router.MaxMultipartMemory = config.Upload.MaxFileSize

	// 根据配置的提供商初始化存储库
	var fileRepo domain.FileRepository
	if config.Upload.Provider == "oss" {
		// 阿里云OSS存储
		ossConfig := storage.OSSStorageConfig{
			BucketName:   config.Upload.OSS.Bucket,
			Endpoint:     config.Upload.OSS.Endpoint,
			AccessKey:    config.Upload.OSS.AccessKey,
			SecretKey:    config.Upload.OSS.SecretKey,
			URLPrefix:    config.Upload.OSS.URLPrefix,
			TempChunkDir: config.Upload.OSS.TempChunkDir,
		}
		ossStorage, err := storage.NewOSSStorage(ossConfig)
		if err != nil {
			logger.Fatal("初始化OSS存储失败: %v", err)
		}
		fileRepo = ossStorage
	} else {
		// 本地存储（默认）
		localConfig := storage.LocalStorageConfig{
			UploadDir:    config.Upload.Local.UploadDir,
			TempChunkDir: config.Upload.Local.TempChunkDir,
			MaxFileSize:  config.Upload.MaxFileSize,
			AllowedTypes: config.Upload.AllowedExtensions,
		}
		localStorage, err := storage.NewLocalStorage(localConfig)
		if err != nil {
			logger.Fatal("初始化本地存储失败: %v", err)
		}
		fileRepo = localStorage
	}

	// 创建服务层
	domainService := domain.NewUploadDomainService(fileRepo)
	appService := application.NewUploadAppService(domainService, fileRepo)

	// 创建HTTP处理器
	uploadHandler := uploadHttp.NewUploadHandler(appService)

	// 注册路由
	uploadHttp.RegisterRoutes(router, uploadHandler)

	// 设置Swagger
	router.GET("/api/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	logger.Info("服务器启动于 %s", addr)
	err = router.Run(addr)
	if err != nil {
		logger.Fatal("服务器启动失败: %v", err)
	}
}

// loadConfig 使用viper加载应用程序配置
func loadConfig() (*Config, error) {
	v := viper.New()

	// 设置配置文件名称、路径、类型
	v.SetConfigName("config")    // 配置文件名称（不带扩展名）
	v.SetConfigType("yaml")      // 配置文件类型
	v.AddConfigPath("./configs") // 配置文件路径

	// 设置默认值
	v.SetDefault("server.host", "localhost")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.timeout", 30)

	v.SetDefault("upload.provider", "local")
	v.SetDefault("upload.max_file_size", 10<<20) // 10MB
	v.SetDefault("upload.max_files", 10)
	v.SetDefault("upload.allowed_extensions", []string{".jpg", ".jpeg", ".png", ".gif", ".pdf", ".doc", ".docx"})

	v.SetDefault("upload.local.upload_dir", "./uploads")
	v.SetDefault("upload.local.temp_chunk_dir", "./temp_chunks")

	// 阿里云OSS配置默认值
	v.SetDefault("upload.oss.bucket", "")
	v.SetDefault("upload.oss.endpoint", "oss-cn-hangzhou.aliyuncs.com")
	v.SetDefault("upload.oss.access_key", "")
	v.SetDefault("upload.oss.secret_key", "")
	v.SetDefault("upload.oss.url_prefix", "")
	v.SetDefault("upload.oss.temp_chunk_dir", "./temp_chunks")

	v.SetDefault("logging.level", "info")

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// 配置文件未找到，使用默认值
			fmt.Println("配置文件未找到，使用默认值")
		} else {
			// 其他错误
			return nil, fmt.Errorf("读取配置文件失败: %w", err)
		}
	}

	// 可选：支持从环境变量读取配置
	v.SetEnvPrefix("UPLOAD") // 环境变量前缀
	v.AutomaticEnv()         // 自动绑定环境变量

	// 解析配置到结构体
	var config Config
	if err := v.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("解析配置失败: %w", err)
	}

	// 创建必要的目录
	ensureDirectories(&config)

	return &config, nil
}

// ensureDirectories 创建文件存储所需的目录
func ensureDirectories(config *Config) {
	// 创建上传目录
	os.MkdirAll(config.Upload.Local.UploadDir, 0755)

	// 创建临时分片目录
	os.MkdirAll(config.Upload.Local.TempChunkDir, 0755)
}

// initLogger 根据配置的级别初始化日志
func initLogger(config *Config) {
	var level logger.Level

	switch config.Logging.Level {
	case "debug":
		level = logger.DEBUG
	case "info":
		level = logger.INFO
	case "warn":
		level = logger.WARN
	case "error":
		level = logger.ERROR
	default:
		level = logger.INFO
	}

	logger.Init(level)
}
