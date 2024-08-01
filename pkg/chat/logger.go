package chat

import (
	"fmt"
	"os"
	"time"
)

type Logger struct {
	file *os.File
}

func NewLogger() (*Logger, error) {
	timestamp := time.Now().Format("20060102_150405")
	filename := fmt.Sprintf("pkg/chat/log/session_%s.txt", timestamp)
	file, err := os.Create(filename)
	if err != nil {
		return nil, fmt.Errorf("could not create log file: %v", err)
	}
	return &Logger{file: file}, nil
}

func (l *Logger) Log(message string) error {
	timestamp := time.Now().Format(time.RFC3339)
	_, err := l.file.WriteString(fmt.Sprintf("%s: %s\n", timestamp, message))
	return err
}

func (l *Logger) Close() {
	l.file.Close()
}
