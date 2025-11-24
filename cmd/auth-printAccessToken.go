package cmd

import (
	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

var auth_printAccessTokenCmd = man.Docs.GetCommand("auth/print-access-token",
	man.WithRun(auth_printAccessToken))

func auth_printAccessToken(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	_, cp := InitProfile(c, false)

	ac := cp.GetAuthCredentials()
	switch ac.AuthType {
	case profiles.AuthTypeClientCredentials:
		c.Printf("Getting access token for %s... ", ac.ClientID)
	case profiles.AuthTypeAccessToken:
		c.Printf("Getting profile's stored access token... ")
	default:
		c.ExitWithError("Invalid auth type", nil)
	}
	tok, err := auth.GetTokenWithProfile(cmd.Context(), cp)
	if err != nil {
		c.Println("failed")
		cli.ExitWithError("Failed to get token", err)
	}
	c.Println("ok")
	c.Printf("Access Token: %s\n", tok.AccessToken)

	c.PrintIfJSON(tok)
}

func init() {
	auth_printAccessTokenCmd.Flags().Bool(
		auth_printAccessTokenCmd.GetDocFlag("json").Name,
		auth_printAccessTokenCmd.GetDocFlag("json").DefaultAsBool(),
		auth_printAccessTokenCmd.GetDocFlag("json").Description,
	)

	authCmd.AddCommand(&auth_printAccessTokenCmd.Command)
}
