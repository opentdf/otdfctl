package cmd

import (
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/tui"
	"github.com/spf13/cobra"
)

func init() {
	cmd := man.Docs.GetCommand("interactive",
		man.WithRun(func(cmd *cobra.Command, args []string) {
			c := cli.New(cmd, args)
			h := NewHandler(c)
			//nolint:errcheck // error does not need to be checked
			tui.StartTea(h)
		}),
	)
	RootCmd.AddCommand(&cmd.Command)
}
