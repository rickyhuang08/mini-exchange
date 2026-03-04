package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/rickyhuang08/mini-exchange.git/config"
)

const (
	AccessLog = "yb-access.log"
	ErrorLog  = "yb-error.log"

	LogLevelInfo  = "INFO"
	LogLevelError = "ERROR"
)

type Logger struct {
	CfgLogger  *config.Logger
	AccessFile *os.File
	ErrorFile  *os.File
	AccessLog  *log.Logger
	ErrorLog   *log.Logger
}

func NewLogger(cfgLogger config.Logger) (*Logger, error) {
	// Build absolute paths for log files
	accessDir := cfgLogger.AccessPath
	errorDir := cfgLogger.ErrorPath

	accessFilePath := filepath.Join(accessDir, AccessLog)
	errorFilePath := filepath.Join(errorDir, ErrorLog)

	// Ensure log directories exist
	if err := os.MkdirAll(accessDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create access log dir: %w", err)
	}
	if err := os.MkdirAll(errorDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("create error log dir: %w", err)
	}

	// Rename existing log files if they exist
	timestamp := time.Now().Format("20060102_150405")
	if _, err := os.Stat(accessFilePath); err == nil {
		backupPath := filepath.Join(accessDir, fmt.Sprintf("yb-access-%s.log", timestamp))
		if err := os.Rename(accessFilePath, backupPath); err != nil {
			return nil, fmt.Errorf("rotate access log: %w", err)
		}
	}
	if _, err := os.Stat(errorFilePath); err == nil {
		backupPath := filepath.Join(errorDir, fmt.Sprintf("yb-error-%s.log", timestamp))
		if err := os.Rename(errorFilePath, backupPath); err != nil {
			return nil, fmt.Errorf("rotate error log: %w", err)
		}
	}

	// Create new log files
	accessFile, err := os.OpenFile(accessFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open access log: %w", err)
	}

	errorFile, err := os.OpenFile(errorFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return nil, fmt.Errorf("open error log: %w", err)
	}
	
	return &Logger{
		CfgLogger:  &cfgLogger,
		AccessFile: accessFile,
		ErrorFile:  errorFile,
		AccessLog:  log.New(accessFile, "ACCESS: ", log.Ldate|log.Ltime|log.LUTC),
		ErrorLog:   log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.LUTC|log.Lshortfile),
	}, nil
}

func (l *Logger) Close() {
	if l.AccessFile != nil {
		l.AccessFile.Close()
	}
	if l.ErrorFile != nil {
		l.ErrorFile.Close()
	}
}

func (l *Logger) LogLevel(level, message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf("[%s] [%s] %s", timestamp, level, message)

	// Log all levels (INFO, DEBUG, ERROR, etc.) to error log file
	l.ErrorLog.Println(formatted)
}

func (l *Logger) LogAccess(message string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	formatted := fmt.Sprintf("[%s] [ACCESS] %s", timestamp, message)

	// Only logs HTTP or middleware access-related messages
	l.AccessLog.Println(formatted)
}