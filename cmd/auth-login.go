package cmd

import (
	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

func auth_codeLogin(cmd *cobra.Command, args []string) {
	fh := cli.NewFlagHelper(cmd)
	clientID := fh.GetOptionalString("client-id")
	tlsNoVerify := fh.GetOptionalBool("tls-no-verify")

	cp := InitProfile(cmd, false)
	printer := cli.NewPrinter(true)

	printer.Println("Initiating login...")
	tok, publicClientID, err := auth.LoginWithPKCE(cmd.Context(), cp.GetEndpoint(), clientID, tlsNoVerify)
	if err != nil {
		cli.ExitWithError("could not authenticate", err)
	}
	printer.Println("ok")

	// Set the auth credentials to profile
	if err := cp.SetAuthCredentials(profiles.AuthCredentials{
		AuthType: profiles.PROFILE_AUTH_TYPE_ACCESS_TOKEN,
		AccessToken: profiles.AuthCredentialsAccessToken{
			PublicClientID: publicClientID,
			AccessToken:    tok.AccessToken,
			Expiration:     tok.Expiry.Unix(),
			RefreshToken:   tok.RefreshToken,
		},
	}); err != nil {
		cli.ExitWithError("failed to set auth credentials", err)
	}

	printer.Println("Storing credentials to profile in keyring...")
	if err := cp.Save(); err != nil {
		printer.Println("failed")
		cli.ExitWithError("An error occurred while storing authentication credentials", err)
	}
	printer.Println("ok")
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
