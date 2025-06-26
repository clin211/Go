package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"
)

// 全局配置变量
var (
	once    sync.Once
	initErr error

	// 直接暴露配置项为包变量
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Log      LogConfig
)

// Config 保存应用程序的所有配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Log      LogConfig      `mapstructure:"log"`
}

// ServerConfig 保存所有服务器相关的配置
type ServerConfig struct {
	Mode    string `mapstructure:"mode"`
	Bind    string `mapstructure:"bind"`
	Port    int    `mapstructure:"port"`
	Address string
}

// DatabaseConfig 保存数据库相关配置
type DatabaseConfig struct {
	Master DBInfo `mapstructure:"master"`
	Slave  DBInfo `mapstructure:"slave"`
}

// DBInfo 保存单个数据库连接信息
type DBInfo struct {
	Type        string        `mapstructure:"type"`
	Host        string        `mapstructure:"host"`
	Port        int           `mapstructure:"port"`
	User        string        `mapstructure:"user"`
	Password    string        `mapstructure:"password"`
	DBName      string        `mapstructure:"db-name"`
	MaxOpen     int           `mapstructure:"max-open"`
	MaxIdle     int           `mapstructure:"max-idle"`
	MaxLifetime time.Duration `mapstructure:"max-lifetime"`
}

// RedisConfig 保存Redis相关配置
type RedisConfig struct {
	Addr     string `mapstructure:"addr"`
	Password string `mapstructure:"password"`
	PoolSize int    `mapstructure:"pool_size"`
	DB       int    `mapstructure:"db"`
}

// LogConfig 保存日志相关配置
type LogConfig struct {
	DisableCaller     bool     `mapstructure:"disable-caller"`
	DisableStacktrace bool     `mapstructure:"disable-stacktrace"`
	Level             string   `mapstructure:"level"`
	Format            string   `mapstructure:"format"`
	OutputPaths       []string `mapstructure:"output-paths"`
}

// Init 初始化全局配置，只会执行一次
// configFile: 如果不为空，则使用指定的配置文件，否则使用默认配置文件
func Init(configFile string) (*Config, error) {
	once.Do(func() {
		var cfg *Config
		cfg, initErr = loadConfig(configFile)
		if cfg != nil {
			// 将配置赋值给包级别变量
			Server = cfg.Server
			Database = cfg.Database
			Redis = cfg.Redis
			Log = cfg.Log
		}
	})

	if initErr != nil {
		return nil, initErr
	}

	return &Config{
		Server:   Server,
		Database: Database,
		Redis:    Redis,
		Log:      Log,
	}, nil
}

// 实际加载配置的内部函数
func loadConfig(configFile string) (*Config, error) {
	v := viper.New()

	// 设置默认名称和路径
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath("./configs")
	v.AddConfigPath(".")

	// 如果指定了配置文件，则使用指定的配置文件
	if configFile != "" {
		v.SetConfigFile(configFile)
		// 提取目录路径，如果不为空则添加到搜索路径中
		configDir := filepath.Dir(configFile)
		if configDir != "" && configDir != "." {
			v.AddConfigPath(configDir)
		}
	}

	// 设置键的映射函数，将连字符(-)替换为下划线(_)
	replacer := strings.NewReplacer("-", "_")
	v.SetEnvKeyReplacer(replacer)

	// 读取配置文件
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// 遍历所有键，创建替换了连字符的别名
	for _, key := range v.AllKeys() {
		if strings.Contains(key, "-") {
			val := v.Get(key)
			newKey := replacer.Replace(key)
			v.Set(newKey, val)
		}
	}

	// 从环境变量覆盖
	v.AutomaticEnv()

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to decode config into struct: %w", err)
	}

	// 设置服务器地址
	if cfg.Server.Bind != "" {
		cfg.Server.Address = fmt.Sprintf("%s:%d", cfg.Server.Bind, cfg.Server.Port)
	} else {
		cfg.Server.Address = fmt.Sprintf(":%d", cfg.Server.Port)
	}

	return &cfg, nil
}
