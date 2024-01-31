package llmgen

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
)

// LLMGenOutputObject struct to hold the output of the command.
type LLMGenOutputObject struct {
	Output        string
	OriginalQuery string
	ExtraData     string
}

type LLMGenArgs struct {
	Query    string
	Function string
}

func NewLLMGenArgs(args ...string) LLMGenArgs {
	var query string
	var function string

	if len(args) > 0 {
		query = args[0]
	}

	if len(args) > 1 {
		function = args[1]
	}

	return LLMGenArgs{
		Query:    query,
		Function: function,
	}
}

func _run_binary(args LLMGenArgs) (*exec.Cmd, *bufio.Reader, error) {

	// Setup the command to run the external binary.
	cmd := exec.Command("./bin/LLMGen/LLMGen", args.Query, args.Function)

	// Create a pipe to the standard output of the command.
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Failed to create stdout pipe: %v", err)
	}

	// Start the command.
	if err := cmd.Start(); err != nil {
		log.Fatalf("Failed to start command: %v", err)
	}

	// Use a buffered reader to read the command's output in real-time.
	reader := bufio.NewReader(stdoutPipe)

	// Wait for the command to finish.
	// if err := cmd.Wait(); err != nil {
	// 	log.Fatalf("Command execution failed: %v", err)
	// }
	return cmd, reader, nil
}

func _collectAndPrintOutput(reader *bufio.Reader) (string, error) {
	var output []byte
	for {
		// Read each line of the output.
		line, err := reader.ReadBytes('\n')
		// line = line[:len(line)-1]
		output = append(output, line...)
		fmt.Print(string(line))

		// Break the loop if an error occurred.
		if err != nil {
			if err != io.EOF {
				log.Printf("Error reading stdout: %v", err)
			}
			break
		}
	}
	return string(output), nil
}

// RawInference executes a given command and streams the output.
// It returns LLMGenOutputObject containing the complete output after execution.
// func RawInference(function string, query string) (LLMGenOutputObject, error) {
func RawInference(query string) (LLMGenOutputObject, error) {
	_, reader, _ := _run_binary(NewLLMGenArgs(query, "/api/raw"))
	var output string
	output, _ = _collectAndPrintOutput(reader)

	// Return the captured output.
	return LLMGenOutputObject{
		Output:        output,
		OriginalQuery: query,
		ExtraData:     string(output),
	}, nil
}

func KnowledgeBaseChat(query string) (LLMGenOutputObject, error) {
	_, reader, _ := _run_binary(NewLLMGenArgs(query, "/api/knowledgebase/chat"))
	var output string
	output, _ = _collectAndPrintOutput(reader)

	// Return the captured output.
	return LLMGenOutputObject{
		Output:        output,
		OriginalQuery: query,
		ExtraData:     string(output),
	}, nil
}

// RawInference executes a given command and streams the output.
// It returns LLMGenOutputObject containing the complete output after execution.
// func RawInference(function string, query string) (LLMGenOutputObject, error) {
func Classify(query string) (LLMGenOutputObject, error) {
	_, reader, _ := _run_binary(NewLLMGenArgs(query, "/api/dlp/classify"))

	var output string
	output, _ = _collectAndPrintOutput(reader)

	// Return the captured output.
	return LLMGenOutputObject{
		Output:        output,
		OriginalQuery: query,
	}, nil
}
