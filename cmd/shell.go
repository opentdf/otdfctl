package cmd

import (
	"github.com/opentdf/otdfctl/cmd/common"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/tui"
	"github.com/spf13/cobra"
)

func init() {
	cmd := man.Docs.GetCommand("shell",
		man.WithRun(func(cmd *cobra.Command, args []string) {
			c := cli.New(cmd, args)
			profile := common.InitProfile(c)
			h := common.NewHandler(c)
			tui.StartShell(h, profile.Name())
		}),
	)
	RootCmd.AddCommand(&cmd.Command)
}
