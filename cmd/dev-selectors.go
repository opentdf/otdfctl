package cmd

import (
	"fmt"
	"strings"

	selectorsgenerated "github.com/opentdf/otdfctl/cmd/generated/dev/selectors"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/spf13/cobra"
)

var selectors []string

// handleDevSelectorsGenerate implements the business logic for the generate command
func handleDevSelectorsGenerate(cmd *cobra.Command, req *selectorsgenerated.GenerateRequest) error {
	c := cli.New(cmd, []string{})
	handler := NewHandler(c)
	defer handler.Close()

	subject := c.Flags.GetRequiredString("subject")

	flattened, err := handlers.FlattenSubjectContext(subject)
	if err != nil {
		cli.ExitWithError("Failed to parse subject context keys and values", err)
	}

	rows := [][]string{}
	for _, item := range flattened {
		rows = append(rows, []string{item.Key, fmt.Sprintf("%v", item.Value)})
	}

	t := cli.NewTabular(rows...)
	cli.PrintSuccessTable(cmd, "", t)
	return nil
}

// handleDevSelectorsTest implements the business logic for the test command
func handleDevSelectorsTest(cmd *cobra.Command, req *selectorsgenerated.TestRequest) error {
	c := cli.New(cmd, []string{})
	handler := NewHandler(c)
	defer handler.Close()

	subject := c.Flags.GetRequiredString("subject")
	
	// Convert single selector string to slice for compatibility with existing logic
	var selectorsList []string
	if req.Flags.Selector != "" {
		selectorsList = strings.Split(req.Flags.Selector, ",")
	}

	if len(selectorsList) == 0 {
		cli.ExitWithError("Must provide at least one selector", nil)
	}

	flattened, err := handlers.FlattenSubjectContext(subject)
	if err != nil {
		cli.ExitWithError("Failed to process subject context keys and values", err)
	}

	rows := [][]string{}
	for _, item := range flattened {
		for _, selector := range selectorsList {
			selector = strings.TrimSpace(selector)
			if selector == item.Key {
				rows = append(rows, []string{item.Key, fmt.Sprintf("%v", item.Value)})
			}
		}
	}

	t := cli.NewTabular(rows...)
	cli.PrintSuccessTable(cmd, "", t)
	return nil
}

func init() {
	// Create commands using generated code with handler functions
	genCmd := selectorsgenerated.NewGenerateCommand(handleDevSelectorsGenerate)
	testCmd := selectorsgenerated.NewTestCommand(handleDevSelectorsTest)

	// Create selectors parent command
	selectorsCmd := &cobra.Command{
		Use:     "selectors",
		Aliases: []string{"sel"},
		Short:   "Selectors",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	
	// Add subcommands
	selectorsCmd.AddCommand(genCmd)
	selectorsCmd.AddCommand(testCmd)

	// Export for use by dev.go
	DevSelectorsCmd = selectorsCmd
}

// DevSelectorsCmd exports the selectors command for use by dev.go
var DevSelectorsCmd *cobra.Command
