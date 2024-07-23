package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
	"gorgonia.org/gorgonia"
)

// Channel for LLM input
var llmInputChan = make(chan string)
var modelFilePath string

// LLMIntakeCmd is the command for handling LLM input
var llmIntakeCmd = &cobra.Command{
	Use:   man.Docs.GetDoc("llm-intake").Use,
	Short: man.Docs.GetDoc("llm-intake").Short,
	Long:  man.Docs.GetDoc("llm-intake").Long,
	Run: func(cmd *cobra.Command, args []string) {
		// Load the model
		model, err := loadModel(modelFilePath)
		if err != nil {
			log.Fatalf("Failed to load model: %v", err)
		}

		// Start the LLM input handler goroutine
		go handleLLMInput(model)

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
	llmIntakeCmd.PersistentFlags().StringVar(
		&modelFilePath,
		"model",
		"default_model.onnx",
		"Path to the model file",
	)
	RootCmd.AddCommand(llmIntakeCmd)
}

// Function to load the model
func loadModel(filePath string) (*gorgonia.ExprGraph, error) {
	model := gorgonia.NewGraph()
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// TODO: Load the model into the graph (this is a placeholder)
	// You need to implement the actual model loading logic here

	return model, nil
}

// Function to handle LLM input
func handleLLMInput(model *gorgonia.ExprGraph) {
	for input := range llmInputChan {
		// Process the LLM input
		fmt.Println("Processing LLM input:", input)

		// Convert input to tensor
		// inputTensor := tensor.New(tensor.Of(tensor.Float32), tensor.WithShape(1, len(input)), tensor.Of(tensor.Float32))
		// TODO: Fill the tensor with input data

		// Create a new VM to run the graph
		machine := gorgonia.NewTapeMachine(model)
		if err := machine.RunAll(); err != nil {
			log.Fatalf("Failed to run the model: %v", err)
		}

		// TODO: Extract and print the output from the model
	}
}

// Non-blocking function to send LLM input
func sendLLMInput(input string) {
	go func() {
		llmInputChan <- input
	}()
}

var llmInput string
