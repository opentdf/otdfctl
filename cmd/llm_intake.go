package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Channel for LLM input
var modelName string // Renamed from modelFilePath to modelName

// New variable for the API URL
var apiURL = "http://localhost:11434/api/generate"

func init() {
	chatCmd.PersistentFlags().StringVar(
		&modelName,
		"model",
		"llama3",
		"Model name for Ollama",
	)
	RootCmd.AddCommand(chatCmd)
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

			requestBody, err := json.Marshal(map[string]interface{}{
				"model":  modelName,
				"prompt": line,
				"stream": true,
			})
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error creating request: %v\n", err)
				continue
			}

			resp, err := http.Post(apiURL, "application/json", bytes.NewBuffer(requestBody))
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error during chat: %v\n", err)
				continue
			}
			defer resp.Body.Close()

			// New: Process each line of the response as it arrives
			responseScanner := bufio.NewScanner(resp.Body)
			for responseScanner.Scan() {
				var result map[string]interface{}
				if err := json.Unmarshal(responseScanner.Bytes(), &result); err != nil {
					fmt.Fprintf(os.Stderr, "Error decoding response: %v\n", err)
					continue
				}

				// Assuming the response is correctly formatted
				if response, ok := result["response"]; ok {
					fmt.Print(response) // Print without newline to continue the sentence
				}
				if done, ok := result["done"]; ok && done.(bool) {
					fmt.Println() // Print a newline when done
					break
				}
			}
		}
	},
}
