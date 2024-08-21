package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

var clientCredentialsCmd = man.Docs.GetCommand("auth/client-credentials",
	man.WithRun(auth_clientCredentials),
	man.WithHiddenFlags("with-client-creds", "with-client-creds-file"),
)

func auth_clientCredentials(cmd *cobra.Command, args []string) {
	cp := InitProfile(cmd, false)

	p := cli.NewPrinter(true)

	var clientId string
	var clientSecret string

	if len(args) > 0 {
		clientId = args[0]
	}
	if len(args) > 1 {
		clientSecret = args[1]
	}

	if clientId == "" {
		clientId = cli.AskForInput("Enter client id: ")
	}
	if clientSecret == "" {
		clientSecret = cli.AskForSecret("Enter client secret: ")
	}

	// Set the client credentials
	cp.SetAuthCredentials(profiles.AuthCredentials{
		AuthType:     profiles.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS,
		ClientId:     clientId,
		ClientSecret: clientSecret,
	})

	// Validate the client credentials
	p.Printf("Validating client credentials for %s... ", cp.GetEndpoint())
	if err := auth.ValidateProfileAuthCredentials(cmd.Context(), cp); err != nil {
		fmt.Println("failed")
		cli.ExitWithError("An error occurred during login. Please check your credentials and try again", err)
	}
	p.Println("ok")

	// Save the client credentials
	p.Print("Storing client ID and secret in keyring... ")
	if err := cp.Save(); err != nil {
		p.Println("failed")
		cli.ExitWithError("An error occurred while storing client credentials", err)
	}
	p.Println("ok")
}

func init() {
	authCmd.AddCommand(&clientCredentialsCmd.Command)
}
