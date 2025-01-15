package cmd

import (
	"github.com/opentdf/otdfctl/internal/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

func auth_codeLogin(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	profileMgr, currProfile := InitProfile(c, false)

	c.Print("Initiating login...")
	tok, publicClientID, err := auth.LoginWithPKCE(
		cmd.Context(),
		currProfile.GetEndpoint(),
		c.FlagHelper.GetOptionalString("client-id"),
		c.FlagHelper.GetOptionalBool("tls-no-verify"),
	)
	if err != nil {
		c.Println("failed")
		c.ExitWithError("could not authenticate", err)
	}
	c.Println("ok")

	// Set the auth credentials to profile
	currProfile.SetAuthCredentials(&auth.AuthCredentials{
		AuthType: auth.AUTH_TYPE_ACCESS_TOKEN,
		AccessToken: &auth.AuthCredentialsAccessToken{
			PublicClientID: publicClientID,
			AccessToken:    tok.AccessToken,
			Expiration:     tok.Expiry.Unix(),
			RefreshToken:   tok.RefreshToken,
		},
	})

	c.Print("Storing credentials to profile in keyring...")
	if err := profileMgr.UpdateProfile(currProfile); err != nil {
		c.Println("failed")
		c.ExitWithError("An error occurred while storing authentication credentials", err)
	}
	c.Println("ok")
}

var codeLoginCmd *man.Doc

func init() {
	codeLoginCmd = man.Docs.GetCommand("auth/login",
		man.WithRun(auth_codeLogin),
	)
	codeLoginCmd.Flags().StringP(
		codeLoginCmd.GetDocFlag("client-id").Name,
		codeLoginCmd.GetDocFlag("client-id").Shorthand,
		codeLoginCmd.GetDocFlag("client-id").Default,
		codeLoginCmd.GetDocFlag("client-id").Description,
	)
	authCmd.AddCommand(&codeLoginCmd.Command)
}
