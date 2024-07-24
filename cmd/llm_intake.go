package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Channel for LLM input
var modelFilePath string

func init() {
	// _ = man.Docs.GetDoc("llm-intake")
	chatCmd.PersistentFlags().StringVar(
		&modelFilePath,
		"model",
		"llama3",
		"Model name for Ollama",
	)
	RootCmd.AddCommand(chatCmd)
}

// // Function to handle LLM input
// func handleLLMInput() {
// 	for input := range llmInputChan {
// 		// Process the LLM input
// 		fmt.Println("Processing LLM input:", input)
// 		response, err := queryOllama(modelFilePath, input, true)
// 		if err != nil {
// 			log.Fatalf("Failed to query Ollama: %v", err)
// 		}
// 		fmt.Println("Ollama response:", response)
// 	}
// }

// // Non-blocking function to send LLM input
// func sendLLMInput(input string) {
// 	go func() {
// 		llmInputChan <- input
// 	}()
// }

func queryOllama(model, prompt string, stream bool) (string, error) {
	payload := map[string]interface{}{
		"model":  model,
		"prompt": prompt,
		"stream": stream,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result struct {
		Response string `json:"response"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	return result.Response, nil
}

// Define the chat command
var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat session with the model",
	Long:  `This command starts an interactive chat session with the model.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting chat session. Type 'exit' to end.")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("> ")
			scanned := scanner.Scan()
			if !scanned {
				return
			}

			line := scanner.Text()
			if strings.TrimSpace(line) == "exit" {
				fmt.Println("Ending chat session.")
				break
			}

			// Replace this with the actual function to send input to the model and receive a response
			response, err := queryOllama("modelFilePath", line, false) // Assuming queryOllama is adapted for chat
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error during chat: %v\n", err)
				continue
			}

			fmt.Println("Model:", response)
		}
	},
}
