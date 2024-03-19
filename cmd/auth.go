package cmd

import "github.com/spf13/cobra"

var (
	// authCmd is the command for managing local authentication session (login, logout, and token caching)
	authCmd = &cobra.Command{
		Use:   "auth",
		Short: "Manage local authentication session",
		Long:  `This command will allow you to manage your local authentication session in regards to the DSP platform.`,
	}
)

func init() {
	rootCmd.AddCommand(authCmd)
}
