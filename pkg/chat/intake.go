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

func UserInputLoop(logger *Logger) {
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

			HandleUserInput(line, logger)
		}
	}
}

func HandleUserInput(input string, logger *Logger) {
	sanitizedInput := SanitizeInput(input)

	// if verbosity is enabled, print the sanitized input, otherwise just log it
	if chatConfig.Chat.Verbose {
		fmt.Printf("\n%s\n\n", sanitizedInput)
	}
	logger.Log(fmt.Sprintf("User: %s", input))
	logger.Log(fmt.Sprintf("Sanitized: %s", sanitizedInput))

	// Channels to receive results
	keywordChan := make(chan []string)
	apiResponseChan := make(chan *http.Response)
	errorChan := make(chan error)
	clearChan := make(chan bool)
	// Start keyword extraction in a goroutine
	go func() {
		keywords, err := ExtractKeywordsFromLLM(sanitizedInput)
		if err != nil {
			errorChan <- err
			return
		}
		keywordChan <- keywords
	}()

	// Start main API call in a goroutine
	go func() {
		requestBody, err := CreateRequestBody(sanitizedInput)
		if err != nil {
			errorChan <- err
			return
		}

		resp, err := SendRequest(requestBody)
		if err != nil {
			errorChan <- err
			return
		}
		apiResponseChan <- resp
	}()

	done := make(chan bool)
	go LoadingAnimation(done)

	startTime := time.Now() // Start timing before sending the request

	// Wait for both results
	var keywords []string
	var resp *http.Response
	for i := 0; i < 2; i++ {
		select {
		case kw := <-keywordChan:
			keywords = kw
			done <- true // Stop the loading animation
			<-clearChan  // Wait for the animation to clear
			fmt.Println()
			fmt.Printf("\rKeywords: [%s]\n", strings.Join(keywords, ", "))
			fmt.Println()
			logger.Log(fmt.Sprintf("Keywords: [%s]", strings.Join(keywords, ", ")))
		case r := <-apiResponseChan:
			resp = r
		case err := <-errorChan:
			ReportError("during chat", err)
			done <- true
			<-clearChan
			return
		}
	}

	defer resp.Body.Close()
	ProcessResponse(resp, logger, startTime)
}

func CreateRequestBody(userInput string) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"model":      chatConfig.Chat.Model,
		"prompt":     userInput,
		"stream":     true,
		"tokenLimit": chatConfig.Chat.TokenLimit,
		"options": map[string]interface{}{
			"useGpu": true,
		},
	})
}

func SendRequest(body []byte) (*http.Response, error) {
	return http.Post(chatConfig.Chat.ApiURL, "application/json", bytes.NewBuffer(body))
}

func ProcessResponse(resp *http.Response, logger *Logger, startTime time.Time) {
	responseScanner := bufio.NewScanner(resp.Body)
	var responseBuffer bytes.Buffer
	var tokenBuffer []string
	tokenCount := 0
	firstTokenReceived := false

	for responseScanner.Scan() {
		if !firstTokenReceived {
			timeBeforeFirstToken = time.Since(startTime)
			firstTokenReceived = true
		}

		result, err := DecodeResponse(responseScanner.Bytes())
		if err != nil {
			ReportError("decoding response", err)
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
		TrackStats(responseScanner.Bytes())

		// Log every logLength tokens
		if tokenCount >= chatConfig.Chat.LogLength {
			LogWithTimestamp(logger, strings.Join(tokenBuffer, ""))
			tokenBuffer = tokenBuffer[:0] // Reset the buffer
			tokenCount = 0
		}
	}

	// Log any remaining tokens
	if tokenCount > 0 {
		LogWithTimestamp(logger, strings.Join(tokenBuffer, ""))
	}

	PrintAndResetStats(startTime)
}

func LogWithTimestamp(logger *Logger, message string) {
	// Remove newline characters from the message
	cleanedMessage := strings.ReplaceAll(message, "\n", "")
	timestamp := time.Now().Format(time.RFC3339)
	logger.Log(fmt.Sprintf("%s: %s", timestamp, cleanedMessage))
}

func DecodeResponse(data []byte) (map[string]interface{}, error) {
	var result map[string]interface{}
	err := json.Unmarshal(data, &result)
	return result, err
}

func ReportError(action string, err error) {
	fmt.Fprintf(os.Stderr, "Error %s: %v\n", action, err)
}

func LoadingAnimation(done chan bool) {
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
