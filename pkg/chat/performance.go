package chat

import (
	"fmt"
	"strings"
	"time"
)

var responseTime time.Duration
var totalTokens int
var timeBeforeFirstToken time.Duration

func trackStats(response []byte) {
	responseTokens := len(strings.Fields(string(response)))
	totalTokens += responseTokens
}

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
