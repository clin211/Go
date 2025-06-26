package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// GormLogger 实现了gorm.Logger接口，使用zap作为底层日志
type GormLogger struct {
	ZapLogger *zap.Logger
	LogLevel  logger.LogLevel
	Config    logger.Config
}

// NewGormLogger 创建一个新的gorm日志适配器
func NewGormLogger(zapLogger *zap.Logger) *GormLogger {
	return &GormLogger{
		ZapLogger: zapLogger.With(zap.String("module", "gorm")),
		LogLevel:  logger.Info,
		Config: logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Info,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	newLogger.Config.LogLevel = level
	return &newLogger
}

// Info 打印信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.ZapLogger.Info(fmt.Sprintf(msg, args...))
	}
}

// Warn 打印警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.ZapLogger.Warn(fmt.Sprintf(msg, args...))
	}
}

// Error 打印错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, args ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.ZapLogger.Error(fmt.Sprintf(msg, args...))
	}
}

// Trace 打印SQL跟踪日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()
	fields := []zap.Field{
		zap.String("sql", sql),
		zap.Int64("rows", rows),
		zap.Duration("elapsed", elapsed),
	}

	switch {
	case err != nil && l.LogLevel >= logger.Error && !errors.Is(err, gorm.ErrRecordNotFound):
		l.ZapLogger.Error("SQL Error", append(fields, zap.Error(err))...)
	case elapsed > l.Config.SlowThreshold && l.LogLevel >= logger.Warn:
		l.ZapLogger.Warn("SQL Slow", fields...)
	case l.LogLevel >= logger.Info:
		l.ZapLogger.Debug("SQL", fields...)
	}
}
