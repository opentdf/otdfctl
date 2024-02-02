package cmd

import (
	"fmt"

	"github.com/opentdf/tructl/pkg/llmgen"
	"github.com/spf13/cobra"
)

var llmgenCmd = &cobra.Command{
	Use:   "llmgen",
	Short: "Just a simple test mode for our llmgen integration",
	Run: func(cmd *cobra.Command, args []string) {
		var userQuery string
		// var route string
		// var function string
		if len(args) > 0 {
			userQuery = args[0]
		}

		// if len(args) > 1 {
		// 	route = args[1]
		// }

		if userQuery == "" {
			fmt.Println("Error: No query provided. Please provide a query.")
			return
		}
		print("thinking...\n\n")
		// llmgen.RawInference(function, userQuery)
		// llmgen.RawInference(userQuery)
		// if route == "cliAgent" {
		// 	llmgen.CLIAgent(userQuery)
		// } else {
		llmgen.KnowledgeBaseChat(userQuery)
		// }
	},
}

func init() {
	rootCmd.AddCommand(llmgenCmd)
}
