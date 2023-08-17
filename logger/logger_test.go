package logger

import (
	"bytes"
	"os"
	"testing"
)

func TestNewLogger(t *testing.T) {
	tempLogFile := "temp_log.txt"
	logger, err := NewLogger(tempLogFile)
	defer func() {
		logger.Close()
		os.Remove(tempLogFile)
	}()

	if err != nil {
		t.Errorf("Expected no error, but got %v", err)
	}

	message := "Test log message\n"
	logger.WriteLog(message)

	fileContents, err := os.ReadFile(tempLogFile)
	if err != nil {
		t.Errorf("Error reading log file: %v", err)
	}
	if !bytes.Contains(fileContents, []byte(message)) {
		t.Errorf("Log message not found in file")
	}
}
