package cmd

import (
	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

func auth_logout(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	cp := InitProfile(c)
	c.Println("Initiating logout...")

	// we can only revoke access tokens stored for the code login flow, not client credentials
	creds := cp.GetAuthCredentials()
	if creds.AuthType == profiles.PROFILE_AUTH_TYPE_ACCESS_TOKEN {
		c.Println("Revoking access token...")
		if err := auth.RevokeAccessToken(
			cmd.Context(),
			cp.GetEndpoint(),
			creds.AccessToken.ClientID,
			creds.AccessToken.RefreshToken,
			c.FlagHelper.GetOptionalBool("tls-no-verify"),
		); err != nil {
			c.Println("failed")
			c.ExitWithError("An error occurred while revoking the access token", err)
		}
	}

	if err := cp.SetAuthCredentials(profiles.AuthCredentials{}); err != nil {
		c.Println("failed")
		c.ExitWithError("An error occurred while logging out", err)
	}
	c.Println("ok")
}

var codeLogoutCmd *man.Doc

func init() {
	codeLogoutCmd = man.Docs.GetCommand("auth/logout",
		man.WithRun(auth_logout),
	)
	authCmd.AddCommand(&codeLogoutCmd.Command)
}
