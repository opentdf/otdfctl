package cmd

import (
	"github.com/opentdf/otdfctl/tui"
	"github.com/spf13/cobra"
)

var interactiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Interactive mode",
	Run: func(cmd *cobra.Command, args []string) {
		tui.StartTea()
	},
}

func init() {
	rootCmd.AddCommand(interactiveCmd)
}
