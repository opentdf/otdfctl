package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// const sanitizationPrompt = "<<SYS>> Alongside the user's prompt, you may also be provided snippets of our documentation to help guide your response. The included documentation is not exhaustive, but will be helpful in contextualizing the needs of the user to the nuances of the codebase and platform.  You are a helpful, respectful and honest assistant for a data security company called Virtru. Your goal is to help users of all kinds use our products and understand how to get the most out of our products via explaining concepts and troubleshooting potential solutions. Always answer as helpfully as possible, while being safe. Your answers should not include any harmful, unethical, racist, sexist, toxic, dangerous, or illegal content. Please ensure that your responses are socially unbiased and positive in nature. If a question does not make any sense, or is not factually coherent, explain why instead of answering something not correct. If you don't know the answer to a question, please don't share false information. The User input is as follows: <</SYS>>"
const sanitizationPrompt = "Be as brief as possible with your response to the following query from the user: "

// SanitizeInput appends the sanitization prompt to the user's input.
func SanitizeInput(input string) string {
	return sanitizationPrompt + input
}

// extractKeywordsFromLLM makes the API call to the LLM and parses the response.
func extractKeywordsFromLLM(input string) ([]string, error) {
	requestBody, err := json.Marshal(map[string]interface{}{
		"model":  chatConfig.Chat.Model,
		"prompt": fmt.Sprintf("From the following prompt, extract between 1 and 8 keywords/acronyms that are most relevant to the list of keywords provided: %v. Prompt: %s", keywords, input),
	})
	if err != nil {
		return nil, err
	}
	resp, err := http.Post(chatConfig.Chat.ApiURL, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	// Read the streaming response
	var completeResponse string
	decoder := json.NewDecoder(resp.Body)
	for decoder.More() {
		var chunk map[string]interface{}
		if err := decoder.Decode(&chunk); err != nil {
			return nil, err
		}
		if response, ok := chunk["response"].(string); ok {
			completeResponse += response
		}
	}
	// Extract keywords from the complete response
	// Assuming the response format is something like: "Here are 6 most relevant keywords extracted from the prompt:\n\n1. keyword1\n2. keyword2\n..."
	lines := strings.Split(completeResponse, "\n")
	var keywords []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "1.") || strings.HasPrefix(line, "2.") || strings.HasPrefix(line, "3.") || strings.HasPrefix(line, "4.") || strings.HasPrefix(line, "5.") || strings.HasPrefix(line, "6.") || strings.HasPrefix(line, "7.") || strings.HasPrefix(line, "8.") {
			keyword := strings.TrimPrefix(line, strings.Split(line, " ")[0])
			keywords = append(keywords, strings.TrimSpace(keyword))
		}
	}
	if len(keywords) == 0 {
		return nil, fmt.Errorf("no keywords found in response")
	}
	//format keywords and return
	cleanedKeywords := formatKeywords(keywords)
	return strings.Split(cleanedKeywords, ", "), nil
}

// formatKeywords formats the keywords as specified
func formatKeywords(keywords []string) string {
	var formatted []string
	for _, keyword := range keywords {
		keyword = strings.ReplaceAll(keyword, "**", "")
		keyword = strings.TrimSpace(keyword)
		formatted = append(formatted, keyword)
	}
	return strings.Join(formatted, ", ")
}

// extractKeywords is a wrapper around extractKeywordsFromLLM to handle errors gracefully.
func extractKeywords(input string) string {
	keywords, err := extractKeywordsFromLLM(input)
	if err != nil {
		// Fallback to dummy implementation in case of error
		return "OpenTDF, otdfctl, troubleshooting"
	}
	return formatKeywords(keywords)
}

// Keywords to pull from
var keywords = []string{
	"OpenTDF", "TDF", "otdfctl", "Virtru", "data security", "encryption", "decryption",
	"attribute-based access control", "ABAC Attribute-Based Access Control", "ZTDF (Zero Trust Data Format)", "policy",
	"authentication", "configuration", "key management", "compliance", "regulations",
	"help", "troubleshooting", "platform", "documentation",
	"user", "products", "access", "control", "attributes", "resources",
	"interactive", "commands", "keys", "protection",
	"verification", "enforcement",
}
