package cmd

import (
	"runtime"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var authCmd = man.Docs.GetCommand("auth", man.WithHiddenFlags(
	"with-client-creds",
	"with-client-creds-file",
))

func init() {
	RootCmd.AddCommand(&authCmd.Command)

	authCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// not supported on linux
		if runtime.GOOS == "linux" {
			cli.ExitWithWarning(
				"Warning: Keyring storage is not available on Linux. Please use the `--with-client-creds` flag or the" +
					"`--with-client-creds-file` flag to provide client credentials securely.",
			)
		}
	}
}
