package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/cmd/common"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var cfgKey string

func updateOutput(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := common.NewHandler(c)
	defer h.Close()

	format := c.Flags.GetRequiredString("format")

	err := config.UpdateOutputFormat(cfgKey, format)
	if err != nil {
		c.ExitWithError("Failed to update output format", err)
	}

	c.Println(cli.SuccessMessage(fmt.Sprintf("Output format updated to %s", format)))
}

var (
	outputDoc = man.Docs.GetCommand("config/output", man.WithRun(updateOutput))
	configDoc = man.Docs.GetCommand("config", man.WithSubcommands(outputDoc))
	Cmd       = &configDoc.Command
)

func init() {
	outputDoc.Flags().String(
		outputDoc.GetDocFlag("format").Name,
		outputDoc.GetDocFlag("format").Default,
		outputDoc.GetDocFlag("format").Description,
	)
}
