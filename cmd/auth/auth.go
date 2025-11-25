package auth

import (
	"runtime"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var (
	authCmd = man.Docs.GetCommand("auth", man.WithHiddenFlags(
		"with-client-creds",
		"with-client-creds-file",
	))

	Cmd = &authCmd.Command
)

// InitCommands sets up all auth subcommands and their flags.
// Call this explicitly from main before executing the root command.
func InitCommands() {
	// Set up platform-specific warning for Linux
	authCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		// not supported on linux
		if runtime.GOOS == "linux" {
			cli.ExitWithWarning(
				"Warning: Keyring storage is not available on Linux. Please use the `--with-client-creds` flag or the" +
					"`--with-client-creds-file` flag to provide client credentials securely.",
			)
		}
	}

	// Register all subcommands
	Cmd.AddCommand(newLoginCmd())
	Cmd.AddCommand(newLogoutCmd())
	Cmd.AddCommand(newClientCredentialsCmd())
	Cmd.AddCommand(newClearClientCredentialsCmd())
	Cmd.AddCommand(newPrintAccessTokenCmd())
}
