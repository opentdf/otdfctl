package cmd

import (
	"fmt"

	"github.com/opentdf/tructl/internal/config"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

// configCmd is the command for managing configuration
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long: `
Manage configuration within 'tructl'.

Configuration is used to manage the configuration of the 'tructl' command line tool and updates the
config .yaml file in the root directory when changes have been made.
`,
}

var updateOutputFormatCmd = &cobra.Command{
	Use:   "output",
	Short: "Define the configured output format",
	Long: `
Define the configured output format for the 'tructl' command line tool. The only supported outputs at
this time are 'json' and styled CLI output, which is the default when unspecified.
`,
	Run: func(cmd *cobra.Command, args []string) {
		h := cli.NewHandler(cmd)
		defer h.Close()

		flagHelper := cli.NewFlagHelper(cmd)
		format := flagHelper.GetRequiredString("format")

		config.UpdateOutputFormat(format)
		fmt.Println(cli.SuccessMessage(fmt.Sprintf("Output format updated to %s", format)))
	},
}

func init() {
	updateOutputFormatCmd.Flags().String("format", "", "'json' or 'styled' as the configured output format")
	configCmd.AddCommand(updateOutputFormatCmd)
	rootCmd.AddCommand(configCmd)
}
