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

var modelName string
var apiURL = "http://localhost:11434/api/generate"

func init() {
	configureChatCommand()
}

func configureChatCommand() {
	chatCmd.PersistentFlags().StringVar(&modelName, "model", "llama3", "Model name for Ollama")
	RootCmd.AddCommand(chatCmd)
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat session with the model",
	Long:  `This command starts an interactive chat session with the model.`,
	Run:   runChatSession,
}

func runChatSession(cmd *cobra.Command, args []string) {
	fmt.Println("Starting chat session. Type 'exit' to end.")
	userInputLoop()
}

func userInputLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			return
		}

		line := scanner.Text()
		if strings.TrimSpace(line) == "exit" {
			fmt.Println("Ending chat session.")
			break
		}

		handleUserInput(line)
	}
}

func handleUserInput(input string) {
	requestBody, err := createRequestBody(input)
	if err != nil {
		reportError("creating request", err)
		return
	}

	resp, err := sendRequest(requestBody)
	if err != nil {
		reportError("during chat", err)
		return
	}
	defer resp.Body.Close()

	processResponse(resp)
}

func createRequestBody(input string) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"model":  modelName,
		"prompt": input,
		"stream": true,
	})
}

func sendRequest(body []byte) (*http.Response, error) {
	return http.Post(apiURL, "application/json", bytes.NewBuffer(body))
}

func processResponse(resp *http.Response) {
	responseScanner := bufio.NewScanner(resp.Body)
	for responseScanner.Scan() {
		result, err := decodeResponse(responseScanner.Bytes())
		if err != nil {
			reportError("decoding response", err)
			continue
		}

		if response, ok := result["response"]; ok {
			fmt.Print(response)
		}
		if done, ok := result["done"].(bool); ok && done {
			fmt.Println()
			break
		}
	}
}

func decodeResponse(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	return result, err
}

func reportError(action string, err error) {
	fmt.Fprintf(os.Stderr, "Error %s: %v\n", action, err)
}
