package logger

import (
	"fmt"
	"log"
	"os"
	"time"
)

// Level 表示日志级别
type Level int

const (
	// DEBUG 调试消息级别
	DEBUG Level = iota
	// INFO 信息消息级别
	INFO
	// WARN 警告消息级别
	WARN
	// ERROR 错误消息级别
	ERROR
	// FATAL 致命消息级别（记录后退出）
	FATAL
)

var (
	// 默认日志记录器
	defaultLogger *Logger
	// 级别名称
	levelNames = map[Level]string{
		DEBUG: "DEBUG",
		INFO:  "INFO",
		WARN:  "WARN",
		ERROR: "ERROR",
		FATAL: "FATAL",
	}
)

// Logger 是一个自定义日志记录器
type Logger struct {
	level  Level
	logger *log.Logger
}

// Init 初始化默认日志记录器
func Init(level Level) {
	defaultLogger = NewLogger(level)
}

// NewLogger 创建具有指定级别的新日志记录器
func NewLogger(level Level) *Logger {
	return &Logger{
		level:  level,
		logger: log.New(os.Stdout, "", 0),
	}
}

// log 在指定级别记录消息
func (l *Logger) log(level Level, format string, args ...interface{}) {
	if level < l.level {
		return
	}

	now := time.Now().Format("2006-01-02 15:04:05")
	msg := fmt.Sprintf(format, args...)
	l.logger.Printf("[%s] %s %s", levelNames[level], now, msg)

	if level == FATAL {
		os.Exit(1)
	}
}

// Debug 记录调试消息
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(DEBUG, format, args...)
}

// Info 记录信息消息
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(INFO, format, args...)
}

// Warn 记录警告消息
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(WARN, format, args...)
}

// Error 记录错误消息
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(ERROR, format, args...)
}

// Fatal 记录致命消息并退出
func (l *Logger) Fatal(format string, args ...interface{}) {
	l.log(FATAL, format, args...)
}

// Debug 使用默认日志记录器记录调试消息
func Debug(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, args...)
	}
}

// Info 使用默认日志记录器记录信息消息
func Info(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, args...)
	}
}

// Warn 使用默认日志记录器记录警告消息
func Warn(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(format, args...)
	}
}

// Error 使用默认日志记录器记录错误消息
func Error(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, args...)
	}
}

// Fatal 使用默认日志记录器记录致命消息并退出
func Fatal(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Fatal(format, args...)
	}
}
