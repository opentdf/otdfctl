package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

var auth_printAccessTokenCmd = man.Docs.GetCommand("auth/print-access-token",
	man.WithRun(auth_printAccessToken))

func auth_printAccessToken(cmd *cobra.Command, args []string) {
	flagHelper := cli.NewFlagHelper(cmd)
	jsonOut := flagHelper.GetOptionalBool("json")

	cp := InitProfile(cmd, false)

	printEnabled := !jsonOut
	p := cli.NewPrinter(printEnabled)

	ac := cp.GetAuthCredentials()
	switch ac.AuthType {
	case profiles.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS:
		p.Printf("Getting access token for %s... ", ac.ClientId)
	case profiles.PROFILE_AUTH_TYPE_ACCESS_TOKEN:
		p.Printf("Getting profile's stored access token... ")
	default:
		cli.ExitWithError("Invalid auth type", nil)
	}
	tok, err := auth.GetTokenWithProfile(cmd.Context(), cp)
	if err != nil {
		p.Println("failed")
		cli.ExitWithError("Failed to get token", err)
	}
	p.Println("ok")
	p.Printf("Access Token: %s\n", tok.AccessToken)

	if jsonOut {
		d, err := json.MarshalIndent(tok, "", "  ")
		if err != nil {
			cli.ExitWithError("Failed to marshal token to json", err)
		}

		fmt.Println(string(d))
		return
	}
}

func init() {
	auth_printAccessTokenCmd.Flags().Bool(
		auth_printAccessTokenCmd.GetDocFlag("json").Name,
		auth_printAccessTokenCmd.GetDocFlag("json").DefaultAsBool(),
		auth_printAccessTokenCmd.GetDocFlag("json").Description,
	)

	authCmd.AddCommand(&auth_printAccessTokenCmd.Command)
}
