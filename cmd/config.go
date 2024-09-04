package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

func config_updateOutput(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	format := c.Flags.GetRequiredString("format")

	err := config.UpdateOutputFormat(cfgKey, format)
	if err != nil {
		c.ExitWithError("Failed to update output format", err)
	}
	c.Println(cli.SuccessMessage(fmt.Sprintf("Output format updated to %s", format)))
}

func init() {
	outputCmd := man.Docs.GetCommand("config/output",
		man.WithRun(config_updateOutput),
	)
	outputCmd.Flags().String(
		outputCmd.GetDocFlag("format").Name,
		outputCmd.GetDocFlag("format").Default,
		outputCmd.GetDocFlag("format").Description,
	)

	cmd := man.Docs.GetCommand("config",
		man.WithSubcommands(outputCmd),
	)
	RootCmd.AddCommand(&cmd.Command)
}
