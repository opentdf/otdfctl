package cmd

import (
	"github.com/opentdf/otdfctl/internal/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

func auth_logout(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	profileMgr, currProfile := InitProfile(c, false)
	c.Println("Initiating logout...")

	// we can only revoke access tokens stored for the code login flow, not client credentials
	creds := currProfile.GetAuthCredentials()
	if creds.AuthType == auth.AUTH_TYPE_ACCESS_TOKEN {
		c.Println("Revoking access token...")
		if err := auth.RevokeAccessToken(
			cmd.Context(),
			currProfile.GetEndpoint(),
			creds.AccessToken.PublicClientID,
			creds.AccessToken.RefreshToken,
			c.FlagHelper.GetOptionalBool("tls-no-verify"),
		); err != nil {
			c.Println("failed")
			c.ExitWithError("An error occurred while revoking the access token", err)
		}
	}

	currProfile.SetAuthCredentials(auth.AuthCredentials{})
	if err := profileMgr.UpdateProfile(currProfile); err != nil {
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
