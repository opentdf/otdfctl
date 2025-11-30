package cmd

import (
	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

func auth_codeLogin(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	cp := InitProfile(c)

	c.Print("Initiating login...")
	clientID := c.FlagHelper.GetRequiredString("client-id")
	port := c.FlagHelper.GetOptionalString("port")
	tok, err := auth.LoginWithPKCE(
		cmd.Context(),
		cp.GetEndpoint(),
		clientID,
		c.FlagHelper.GetOptionalBool("tls-no-verify"),
		port,
	)
	if err != nil {
		c.Println("failed")
		c.ExitWithError("could not authenticate", err)
	}
	c.Println("ok")

	// Set the auth credentials to profile
	if err := cp.SetAuthCredentials(profiles.AuthCredentials{
		AuthType: profiles.PROFILE_AUTH_TYPE_ACCESS_TOKEN,
		AccessToken: profiles.AuthCredentialsAccessToken{
			ClientID:     clientID,
			AccessToken:  tok.AccessToken,
			Expiration:   tok.Expiry.Unix(),
			RefreshToken: tok.RefreshToken,
		},
	}); err != nil {
		c.ExitWithError("failed to set auth credentials", err)
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

	// intentionally a string flag to support an empty port which represents a dynamic port
	codeLoginCmd.Flags().StringP(
		codeLoginCmd.GetDocFlag("port").Name,
		codeLoginCmd.GetDocFlag("port").Shorthand,
		codeLoginCmd.GetDocFlag("port").Default,
		codeLoginCmd.GetDocFlag("port").Description,
	)
	authCmd.AddCommand(&codeLoginCmd.Command)
}
