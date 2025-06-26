package db

import (
	"fmt"
	"sync"
	"time"

	"github.com/clin211/go-cicd-github-actions-docker/internal/config"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	masterDB *gorm.DB
	slaveDB  *gorm.DB
	once     sync.Once
)

// GetMasterDB 获取主数据库连接
func GetMasterDB() *gorm.DB {
	return masterDB
}

// GetSlaveDB 获取从数据库连接
func GetSlaveDB() *gorm.DB {
	return slaveDB
}

// Initialize 初始化数据库连接
func Initialize(logger *zap.Logger) error {
	var err error

	once.Do(func() {
		// 初始化主库连接
		masterDB, err = connect(config.Database.Master, logger)
		if err != nil {
			logger.Error("Failed to connect to master database", zap.Error(err))
			return
		}
		logger.Info("Connected to master database")

		// 如果从库配置与主库不同，则初始化从库连接
		if config.Database.Slave.Host != config.Database.Master.Host ||
			config.Database.Slave.Port != config.Database.Master.Port ||
			config.Database.Slave.DBName != config.Database.Master.DBName {
			slaveDB, err = connect(config.Database.Slave, logger)
			if err != nil {
				logger.Error("Failed to connect to slave database", zap.Error(err))
				return
			}
			logger.Info("Connected to slave database")
		} else {
			// 主从配置相同，使用同一个连接
			slaveDB = masterDB
			logger.Info("Using master connection for slave operations")
		}
	})

	return err
}

// Close 关闭数据库连接
func Close() error {
	if masterDB != nil {
		db, err := masterDB.DB()
		if err != nil {
			return err
		}
		if err := db.Close(); err != nil {
			return err
		}
	}

	// 如果从库与主库不是同一个连接，则单独关闭
	if slaveDB != nil && slaveDB != masterDB {
		db, err := slaveDB.DB()
		if err != nil {
			return err
		}
		if err := db.Close(); err != nil {
			return err
		}
	}

	return nil
}

// 连接到数据库
func connect(dbConfig config.DBInfo, logger *zap.Logger) (*gorm.DB, error) {
	if dbConfig.Type != "mysql" {
		return nil, fmt.Errorf("unsupported database type: %s", dbConfig.Type)
	}

	// 构建DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dbConfig.User,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.DBName,
	)

	// 创建gorm日志适配器
	gormLogger := NewGormLogger(logger)

	// GORM配置
	gormConfig := &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true, // 使用单数表名
		},
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// 设置最大连接数
	sqlDB.SetMaxOpenConns(dbConfig.MaxOpen)
	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdle)
	// 设置连接最大生命周期
	sqlDB.SetConnMaxLifetime(dbConfig.MaxLifetime * time.Second)

	return db, nil
}
