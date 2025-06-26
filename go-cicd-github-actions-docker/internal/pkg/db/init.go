package db

import (
	"go.uber.org/zap"
)

// InitAll 初始化所有数据库连接
func InitAll(logger *zap.Logger) error {
	// 初始化MySQL
	if err := Initialize(logger); err != nil {
		return err
	}

	// 初始化Redis
	if err := InitializeRedis(logger); err != nil {
		return err
	}

	return nil
}

// CloseAll 关闭所有数据库连接
func CloseAll() error {
	// 关闭MySQL
	if err := Close(); err != nil {
		return err
	}

	// 关闭Redis
	if err := CloseRedis(); err != nil {
		return err
	}

	return nil
}
