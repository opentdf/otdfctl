package utils

import (
	"fmt"
	"io"
	"os"
)

func ReadBytesFromFile(filePath string) ([]byte, error) {
	fileToEncrypt, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file at path %s: %w", filePath, err)
	}
	defer fileToEncrypt.Close()

	bytes, err := io.ReadAll(fileToEncrypt)
	if err != nil {
		return nil, fmt.Errorf("failed to read bytes from file at path %s: %w", filePath, err)
	}
	return bytes, nil
}
