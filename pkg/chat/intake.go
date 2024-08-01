package chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var responseTime time.Duration
var totalTokens int

func init() {
	err := LoadConfig("chat_config.json")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading chat_config: %v\n", err)
		os.Exit(1)
	}
	configureChatCommand()
}

func configureChatCommand() {
	// TODO: Make more configurable without losing dynamic selection, keeping it accessible via command line flag.
	chatCmd.PersistentFlags().StringVar(&chat_config.Model, "model", chat_config.Model, "Model name for Ollama")
	RootCmd.AddCommand(chatCmd)
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start a chat session with a LLM helper aid",
	Long:  `This command starts an interactive chat session with a local LLM to help with setup, debugging, or generic troubleshooting`,
	Run:   runChatSession,
}

func runChatSession(cmd *cobra.Command, args []string) {
	fmt.Println("Starting chat session. Type 'exit' to end.")
	userInputLoop()
}

// TODO: Make additional 'exit criteria' for the chat session, CTRL+C, etc.
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

// Wraps the user's input and displaying the model's response
func handleUserInput(input string) {
	sanitizedInput := SanitizeInput(input)
	fmt.Printf("\n%s\n\n", sanitizedInput)
	requestBody, err := createRequestBody(sanitizedInput)
	if err != nil {
		reportError("creating request", err)
		return
	}

	done := make(chan bool)
	go loadingAnimation(done)

	resp, err := sendRequest(requestBody)
	if err != nil {
		reportError("during chat", err)
		done <- true
		return
	}
	defer resp.Body.Close()

	done <- true

	processResponse(resp)
}

// Constructs JSON payload for the model's API
func createRequestBody(userInput string) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"model":  chat_config.Model,
		"prompt": userInput,
		"stream": true,
	})
}

func sendRequest(body []byte) (*http.Response, error) {
	return http.Post(chat_config.ApiURL, "application/json", bytes.NewBuffer(body))
}

func trackStats(response []byte) {
	responseTokens := len(strings.Fields(string(response)))
	totalTokens += responseTokens
}

func printAndResetStats(startTime time.Time) {
	responseTime = time.Since(startTime)
	fmt.Printf("\nTotal Response Time: %v\n", responseTime)
	fmt.Printf("Total Tokens: %d\n", totalTokens)

	// Reset stats
	responseTime = 0
	totalTokens = 0
}

func processResponse(resp *http.Response) {
	responseScanner := bufio.NewScanner(resp.Body)
	startTime := time.Now()
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
		trackStats(responseScanner.Bytes())
	}
	printAndResetStats(startTime)
}

// Decodes a single JSON response from the model's API,
func decodeResponse(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	return result, err
}

func reportError(action string, err error) {
	fmt.Fprintf(os.Stderr, "Error %s: %v\n", action, err)
}

func loadingAnimation(done chan bool) {
	chars := []rune{'|', '/', '-', '\\'}
	for {
		select {
		case <-done:
			fmt.Print("\r") // Clear the loading animation
			return
		default:
			for _, char := range chars {
				fmt.Printf("\r%c", char)
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}
