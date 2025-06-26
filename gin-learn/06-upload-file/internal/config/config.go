package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config 配置结构体
type Config struct {
	Server ServerConfig `mapstructure:"server"`
	Logger LoggerConfig `mapstructure:"logger"`
	Local  LocalConfig  `mapstructure:"local"`
	AliOSS AliOSSConfig `mapstructure:"ali_oss"`
	MinIO  MinIOConfig  `mapstructure:"minio"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	Mode         string `mapstructure:"mode"`
	ReadTimeout  string `mapstructure:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

// LocalConfig 本地存储配置
type LocalConfig struct {
	UploadDir         string   `mapstructure:"upload_dir"`
	AllowedExtensions []string `mapstructure:"allowed_extensions"`
	MaxFileSize       int64    `mapstructure:"max_file_size"`
}

// AliOSSConfig 阿里云OSS配置
type AliOSSConfig struct {
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	BucketName      string `mapstructure:"bucket_name"`
	Endpoint        string `mapstructure:"endpoint"`
}

// MinIOConfig MinIO配置
type MinIOConfig struct {
	AccessKey  string `mapstructure:"access_key"`
	SecretKey  string `mapstructure:"secret_key"`
	BucketName string `mapstructure:"bucket_name"`
	Endpoint   string `mapstructure:"endpoint"`
}

// 包级别变量
var (
	Server ServerConfig
	Logger LoggerConfig
	Local  LocalConfig
	AliOSS AliOSSConfig
	MinIO  MinIOConfig
)

// Init 初始化配置
func Init(configPath string) error {
	viper.SetConfigFile(configPath)
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %w", err)
	}

	// 解析到包级别变量
	if err := viper.UnmarshalKey("server", &Server); err != nil {
		return fmt.Errorf("解析server配置失败: %w", err)
	}
	if err := viper.UnmarshalKey("logger", &Logger); err != nil {
		return fmt.Errorf("解析logger配置失败: %w", err)
	}
	if err := viper.UnmarshalKey("local", &Local); err != nil {
		return fmt.Errorf("解析local配置失败: %w", err)
	}
	if err := viper.UnmarshalKey("ali_oss", &AliOSS); err != nil {
		return fmt.Errorf("解析ali_oss配置失败: %w", err)
	}
	if err := viper.UnmarshalKey("minio", &MinIO); err != nil {
		return fmt.Errorf("解析minio配置失败: %w", err)
	}

	return nil
}
