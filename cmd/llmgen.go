package cmd

import (
	"github.com/spf13/cobra"
)

var llmgenCmd = &cobra.Command{
	Use:   "llmgen",
	Short: "Just a simple test mode for our llmgen integration",
	Run: func(cmd *cobra.Command, args []string) {
		println("llmgen test - PASS")
	},
}

func init() {
	rootCmd.AddCommand(llmgenCmd)
}
