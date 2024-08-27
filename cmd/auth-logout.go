package cmd

import (
	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

func auth_logout(cmd *cobra.Command, args []string) {
	fh := cli.NewFlagHelper(cmd)
	tlsNoVerify := fh.GetOptionalBool("tls-no-verify")
	cp := InitProfile(cmd, false)
	printer := cli.NewPrinter(true)
	printer.Println("Initiating logout...")

	// we can only revoke access tokens stored for the code login flow, not client credentials
	creds := cp.GetAuthCredentials()
	if creds.AuthType == profiles.PROFILE_AUTH_TYPE_ACCESS_TOKEN {
		printer.Println("Revoking access token...")
		if err := auth.RevokeAccessToken(
			cmd.Context(),
			cp.GetEndpoint(),
			creds.AccessToken.PublicClientID,
			creds.AccessToken.RefreshToken,
			tlsNoVerify,
		); err != nil {
			printer.Println("failed")
			cli.ExitWithError("An error occurred while revoking the access token", err)
		}
	}

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
