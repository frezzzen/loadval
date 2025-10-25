package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	WARNING
	ERROR
)

const (
	maxLogSize  = 10 * 1024 * 1024 // 10MB
	maxLogFiles = 5
)

type Logger struct {
	mu          sync.Mutex
	logFile     *os.File
	logger      *log.Logger
	logLevel    LogLevel
	logFilePath string
}

var globalLogger *Logger

func (l LogLevel) String() string {
	switch l {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARNING"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func InitLogger() error {
	appData := os.Getenv("LOCALAPPDATA")
	if appData == "" {
		appData = os.Getenv("APPDATA")
	}
	if appData == "" {
		return fmt.Errorf("LOCALAPPDATA/APPDATA environment variable is not set")
	}

	logDir := filepath.Join(appData, "LOADVAL", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	logFilePath := filepath.Join(logDir, "loadval.log")

	if info, err := os.Stat(logFilePath); err == nil && info.Size() > maxLogSize {
		rotateLogs(logDir)
	}

	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)

	globalLogger = &Logger{
		logFile:     logFile,
		logger:      log.New(multiWriter, "", 0),
		logLevel:    INFO,
		logFilePath: logFilePath,
	}

	globalLogger.Info("Logger initialized successfully")
	globalLogger.Info(fmt.Sprintf("Log file location: %s", logFilePath))

	return nil
}

func rotateLogs(logDir string) {
	for i := maxLogFiles - 1; i >= 0; i-- {
		var oldPath, newPath string
		if i == 0 {
			oldPath = filepath.Join(logDir, "loadval.log")
		} else {
			oldPath = filepath.Join(logDir, fmt.Sprintf("loadval.log.%d", i))
		}
		newPath = filepath.Join(logDir, fmt.Sprintf("loadval.log.%d", i+1))

		if _, err := os.Stat(oldPath); err == nil {
			if i == maxLogFiles-1 {
				os.Remove(oldPath) // Remove the oldest log
			} else {
				os.Rename(oldPath, newPath)
			}
		}
	}
}

func (l *Logger) log(level LogLevel, message string) {
	if level < l.logLevel {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] [%s] %s", timestamp, level.String(), message)
	l.logger.Println(logMessage)

	if info, err := os.Stat(l.logFilePath); err == nil && info.Size() > maxLogSize {
		l.rotateLog()
	}
}

func (l *Logger) rotateLog() {
	l.logFile.Close()

	logDir := filepath.Dir(l.logFilePath)
	rotateLogs(logDir)

	logFile, err := os.OpenFile(l.logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to reopen log file after rotation: %v\n", err)
		return
	}

	multiWriter := io.MultiWriter(os.Stdout, logFile)
	l.logFile = logFile
	l.logger = log.New(multiWriter, "", 0)
}

func (l *Logger) Debug(message string) {
	l.log(DEBUG, message)
}

func (l *Logger) Info(message string) {
	l.log(INFO, message)
}

func (l *Logger) Warning(message string) {
	l.log(WARNING, message)
}

func (l *Logger) Error(message string) {
	l.log(ERROR, message)
}

func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.logLevel = level
}

func (l *Logger) Close() {
	if l.logFile != nil {
		l.logFile.Close()
	}
}

func LogDebug(message string) {
	if globalLogger != nil {
		globalLogger.Debug(message)
	}
}

func LogInfo(message string) {
	if globalLogger != nil {
		globalLogger.Info(message)
	}
}

func LogWarning(message string) {
	if globalLogger != nil {
		globalLogger.Warning(message)
	}
}

func LogError(message string) {
	if globalLogger != nil {
		globalLogger.Error(message)
	}
}

func GetLogFilePath() string {
	if globalLogger != nil {
		return globalLogger.logFilePath
	}
	return ""
}

func CloseLogger() {
	if globalLogger != nil {
		globalLogger.Info("Closing logger")
		globalLogger.Close()
	}
}
