package cmd

import (
	"github.com/opentdf/otdfctl/internal/auth"
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
	c := cli.New(cmd, args)
	profileMgr, currProfile := InitProfile(c, false)

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
	currProfile.SetAuthCredentials(auth.AuthCredentials{
		AuthType:     auth.AUTH_TYPE_CLIENT_CREDENTIALS,
		ClientID:     clientId,
		ClientSecret: clientSecret,
	})

	// Validate the client credentials
	c.Printf("Validating client credentials for %s... ", currProfile.GetEndpoint())
	if err := profiles.ValidateProfileAuthCredentials(cmd.Context(), currProfile); err != nil {
		c.Println("failed")
		c.ExitWithError("An error occurred during login. Please check your credentials and try again", err)
	}
	c.Println("ok")

	// Save the client credentials
	c.Print("Storing client ID and secret in keyring... ")

	if err := profileMgr.UpdateProfile(currProfile); err != nil {
		c.Println("failed")
		c.ExitWithError("An error occurred while storing client credentials", err)
	}
	c.Println("ok")
}

func init() {
	authCmd.AddCommand(&clientCredentialsCmd.Command)
}
