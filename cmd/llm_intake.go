package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

// Channel for LLM input
var llmInputChan chan string

// LLMIntakeCmd is the command for handling LLM input
var llmIntakeCmd = &cobra.Command{
	Use:   man.Docs.GetDoc("llm-intake").Use,
	Short: man.Docs.GetDoc("llm-intake").Short,
	Long:  man.Docs.GetDoc("llm-intake").Long,
	Run: func(cmd *cobra.Command, args []string) {
		// Start the LLM input handler goroutine
		go handleLLMInput()

		// Example usage of sending LLM input
		for _, input := range args {
			sendLLMInput(input)
		}
	},
}

func init() {
	doc := man.Docs.GetDoc("llm-intake")
	llmIntakeCmd.PersistentFlags().StringVar(
		&llmInput,
		doc.GetDocFlag("input").Name,
		doc.GetDocFlag("input").Default,
		doc.GetDocFlag("input").Description,
	)
	RootCmd.AddCommand(llmIntakeCmd)
}

// Function to handle LLM input
func handleLLMInput() {
	for input := range llmInputChan {
		// Process the LLM input
		fmt.Println("Processing LLM input:", input)
		// TODO: Add actual LLM processing logic here
	}
}

// Non-blocking function to send LLM input
func sendLLMInput(input string) {
	go func() {
		llmInputChan <- input
	}()
}

var llmInput string
