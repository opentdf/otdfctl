package cmd

import (
	"fmt"

	configgenerated "github.com/opentdf/otdfctl/cmd/generated/config"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/spf13/cobra"
)

// handleConfigOutput implements the business logic for the config output command
func handleConfigOutput(cmd *cobra.Command, req *configgenerated.OutputRequest) error {
	c := cli.New(cmd, []string{})
	h := NewHandler(c)
	defer h.Close()

	format := req.Flags.Format

	err := config.UpdateOutputFormat(cfgKey, format)
	if err != nil {
		c.ExitWithError("Failed to update output format", err)
	}

	c.Println(cli.SuccessMessage(fmt.Sprintf("Output format updated to %s", format)))
	return nil
}

// handleConfig implements the parent config command (shows help if called without subcommands)
func handleConfig(cmd *cobra.Command, req *configgenerated.ConfigRequest) error {
	return cmd.Help()
}

func init() {
	// Create commands using generated constructors with handler functions
	configCmd := configgenerated.NewConfigCommand(handleConfig)
	outputCmd := configgenerated.NewOutputCommand(handleConfigOutput)

	// Add subcommand to parent
	configCmd.AddCommand(outputCmd)

	// Add to root command
	RootCmd.AddCommand(configCmd)
}
