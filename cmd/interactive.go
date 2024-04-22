package cmd

import (
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/tui"
	"github.com/spf13/cobra"
)

func init() {
	cmd := man.Docs.GetCommand("interactive",
		man.WithRun(func(cmd *cobra.Command, args []string) {
			tui.StartTea()
		}),
	)
	RootCmd.AddCommand(&cmd.Command)
}
