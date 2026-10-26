package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

var (
	logFile *os.File
	logger  *log.Logger
)

func InitLogger() error {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Create log file with timestamp
	timestamp := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logsDir, fmt.Sprintf("grafisearch-%s.log", timestamp))
	
	var err error
	logFile, err = os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return fmt.Errorf("failed to create log file: %v", err)
	}

	// Create multi-writer to write to both console and file
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	logger = log.New(multiWriter, "", log.LstdFlags|log.Lmicroseconds)

	// Replace the default logger
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	LogInfo("Logger initialized", map[string]interface{}{
		"log_file": logPath,
		"timestamp": time.Now(),
	})

	return nil
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}

func LogInfo(message string, data map[string]interface{}) {
	logWithData("INFO", message, data)
}

func LogError(message string, err error, data map[string]interface{}) {
	if data == nil {
		data = make(map[string]interface{})
	}
	data["error"] = err.Error()
	logWithData("ERROR", message, data)
}

func LogWarn(message string, data map[string]interface{}) {
	logWithData("WARN", message, data)
}

func logWithData(level, message string, data map[string]interface{}) {
	logMsg := fmt.Sprintf("[%s] %s", level, message)
	
	if len(data) > 0 {
		logMsg += " |"
		for key, value := range data {
			logMsg += fmt.Sprintf(" %s=%v", key, value)
		}
	}
	
	log.Println(logMsg)
}