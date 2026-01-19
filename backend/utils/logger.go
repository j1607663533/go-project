package utils

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
	DebugLogger *log.Logger
)

// InitLogger 初始化日志
func InitLogger() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	DebugLogger = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// LogInfo 记录信息日志
func LogInfo(format string, v ...interface{}) {
	if InfoLogger == nil {
		InitLogger()
	}
	InfoLogger.Printf(format, v...)
}

// LogError 记录错误日志
func LogError(format string, v ...interface{}) {
	if ErrorLogger == nil {
		InitLogger()
	}
	ErrorLogger.Printf(format, v...)
}

// LogDebug 记录调试日志
func LogDebug(format string, v ...interface{}) {
	if DebugLogger == nil {
		InitLogger()
	}
	DebugLogger.Printf(format, v...)
}
