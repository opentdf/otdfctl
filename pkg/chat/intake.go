package chat

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

var responseTime time.Duration
var totalTokens int

// What are other invocation methods for the chat model? --ask? --file? --batch? --interactive?

// TODO add timing/performance metrics for _before_ the model begins responding not just when the first response comes back

// TODO: add a 'one-off' --ask flag to allow for a single question to be asked and answered, DRYing existing chat code

func userInputLoop(logger *Logger) {
	scanner := bufio.NewScanner(os.Stdin)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	for {
		fmt.Print("> ")
		select {
		case <-sigChan:
			fmt.Println("\nReceived interrupt signal. Ending chat session.")
			return
		default:
			if !scanner.Scan() {
				return
			}

			line := scanner.Text()
			if strings.TrimSpace(line) == "exit" || strings.TrimSpace(line) == "quit" {
				fmt.Println("Ending chat session.")
				return
			}

			handleUserInput(line, logger)
		}
	}
}

func handleUserInput(input string, logger *Logger) {
	sanitizedInput := SanitizeInput(input)
	fmt.Printf("\n%s\n\n", sanitizedInput)
	logger.Log(fmt.Sprintf("User: %s", input))
	logger.Log(fmt.Sprintf("Sanitized: %s", sanitizedInput))

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

	processResponse(resp, logger)
}

func createRequestBody(userInput string) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"model":      chatConfig.Chat.Model,
		"prompt":     userInput,
		"stream":     true,
		"tokenLimit": chatConfig.Chat.TokenLimit,
	})
}

func sendRequest(body []byte) (*http.Response, error) {
	return http.Post(chatConfig.Chat.ApiURL, "application/json", bytes.NewBuffer(body))
}

func trackStats(response []byte) {
	responseTokens := len(strings.Fields(string(response)))
	totalTokens += responseTokens
}

func printAndResetStats(startTime time.Time) {
	responseTime = time.Since(startTime)
	if chatConfig.Chat.Verbose {
		fmt.Printf("\nTotal Response Time: %v\n", responseTime)
		fmt.Printf("Total Tokens: %d\n", totalTokens)
	}

	// Reset stats
	responseTime = 0
	totalTokens = 0
}

func processResponse(resp *http.Response, logger *Logger) {
	responseScanner := bufio.NewScanner(resp.Body)
	startTime := time.Now()
	var responseBuffer bytes.Buffer
	var tokenBuffer []string
	tokenCount := 0

	for responseScanner.Scan() {
		result, err := decodeResponse(responseScanner.Bytes())
		if err != nil {
			reportError("decoding response", err)
			continue
		}

		if response, ok := result["response"]; ok {
			fmt.Print(response)
			responseBuffer.WriteString(fmt.Sprintf("%s", response))
			tokenBuffer = append(tokenBuffer, fmt.Sprintf("%s", response))
			tokenCount++
		}
		if done, ok := result["done"].(bool); ok && done {
			fmt.Println()
			break
		}
		trackStats(responseScanner.Bytes())

		// Log every logLength tokens
		if tokenCount >= chatConfig.Chat.LogLength {
			logWithTimestamp(logger, strings.Join(tokenBuffer, ""))
			tokenBuffer = tokenBuffer[:0] // Reset the buffer
			tokenCount = 0
		}
	}

	// Log any remaining tokens
	if tokenCount > 0 {
		logWithTimestamp(logger, strings.Join(tokenBuffer, ""))
	}

	printAndResetStats(startTime)
}

func logWithTimestamp(logger *Logger, message string) {
	// Remove newline characters from the message
	cleanedMessage := strings.ReplaceAll(message, "\n", "")
	timestamp := time.Now().Format(time.RFC3339)
	logger.Log(fmt.Sprintf("%s: %s", timestamp, cleanedMessage))
}

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
