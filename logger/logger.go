package logger

import (
	"fmt"
	"os"
)

type logger struct {
	file *os.File
}

func NewLogger(fileName string) (*logger, error) {
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}
	return &logger{file: file}, nil
}

func (l *logger) WriteLog(message string) {
	fmt.Println(message)
	_, err := l.file.WriteString(message)
	if err != nil {
		fmt.Println("Błąd podczas zapisu do pliku:", err)
	}
}

func (l *logger) Close() {
	l.file.Close()
}
