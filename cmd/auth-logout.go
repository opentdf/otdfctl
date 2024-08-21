package cmd

import (
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

func auth_logout(cmd *cobra.Command, args []string) {
	cp := InitProfile(cmd, false)
	printer := cli.NewPrinter(true)

	printer.Println("Initiating logout...")
	if err := cp.SetAuthCredentials(profiles.AuthCredentials{}); err != nil {
		printer.Println("failed")
		cli.ExitWithError("An error occurred while logging out", err)
	}
	printer.Println("ok")
}

var codeLogoutCmd *man.Doc

func init() {
	codeLogoutCmd = man.Docs.GetCommand("auth/logout",
		man.WithRun(auth_logout),
	)
	authCmd.AddCommand(&codeLogoutCmd.Command)
}
