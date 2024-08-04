package chat

import (
	"fmt"
	"strings"
	"time"
)

var responseTime time.Duration
var totalTokens int
var timeBeforeFirstToken time.Duration

// Tracks total number of tokens in the response
func trackStats(response []byte) {
	responseTokens := len(strings.Fields(string(response)))
	totalTokens += responseTokens
}

// Prints the response time, time before the first token, and total tokens if verbose mode is enabled. Then resets the stats once complete.
func printAndResetStats(startTime time.Time) {
	responseTime = time.Since(startTime)
	if chatConfig.Chat.Verbose {
		fmt.Printf("\nTotal Response Time: %v\n", responseTime)
		fmt.Printf("Time Before First Token: %v\n", timeBeforeFirstToken)
		fmt.Printf("Total Tokens: %d\n", totalTokens)
	}

	// Reset stats
	responseTime = 0
	totalTokens = 0
	timeBeforeFirstToken = 0
}
