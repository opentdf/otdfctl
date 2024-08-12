package cli

import (
	"fmt"
	"io"
	"os"
)

func ReadFromFileOrPipe(file string, pipe *os.File) []byte {
	if file != "" {
		return ReadFromFile(file)
	} else {
		if pipe == nil {
			pipe = os.Stdin
		}
		return ReadFromPipe(pipe)
	}
}

func ReadFromPipe(in *os.File) []byte {
	stat, err := in.Stat()
	if err != nil {
		ExitWithError("failed to read stat from stdin", err)
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		buf, err := io.ReadAll(in)
		if err != nil {
			ExitWithError("failed to scan bytes from stdin", err)
		}
		return buf
	}
	return nil
}

func ReadFromFile(filePath string) []byte {
	fileToEncrypt, err := os.Open(filePath)
	if err != nil {
		ExitWithError(fmt.Sprintf("Failed to open file at path: %s", filePath), err)
	}
	defer fileToEncrypt.Close()

	bytes, err := io.ReadAll(fileToEncrypt)
	if err != nil {
		ExitWithError(fmt.Sprintf("Failed to read bytes from file at path: %s", filePath), err)
	}
	return bytes
}
