package cmd

import (
	"errors"
	"fmt"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var clientCredentialsCmd = man.Docs.GetCommand("auth/client-credentials",
	man.WithRun(auth_clientCredentials),
)

func auth_clientCredentials(cmd *cobra.Command, args []string) {
	flagHelper := cli.NewFlagHelper(cmd)
	clientID := flagHelper.GetOptionalString("client-id")
	clientSecret := flagHelper.GetOptionalString("client-secret")

	// if not provided by flag, check keyring cache for clientID
	if clientID == "" {
		fmt.Println("No client-id provided. Attempting to retrieve the default from keyring.")
		retrievedClientID, err := handlers.GetClientIDFromCache()
		if err != nil || retrievedClientID == "" {
			cli.ExitWithError("Please provide required flag: (client-id)", errors.New("no client-id found"))
		} else {
			clientID = retrievedClientID
			fmt.Println(cli.SuccessMessage("Retrieved stored client-id from keyring"))
		}
	}

	// check if we have a clientSecret in the keyring, if a null value is passed in
	if clientSecret == "" {
		retrievedSecret, krErr := keyring.Get(handlers.TOKEN_URL, clientID)
		if krErr == nil || retrievedSecret == "" {
			cli.ExitWithError("Please provide required flag: (client-secret)", errors.New("no client-secret found"))
		} else {
			clientSecret = retrievedSecret
			fmt.Println(cli.SuccessMessage("Retrieved stored client-secret from keyring"))
		}
	}

	_, err := handlers.GetTokenWithClientCredentials(cmd.Context(), clientID, clientSecret, handlers.TOKEN_URL, false)
	if err != nil {
		cli.ExitWithError("An error occurred during login. Please check your credentials and try again", err)
	}

	fmt.Println(cli.SuccessMessage("Successfully logged in with client ID and secret"))
}

func init() {
	clientCredentialsCmd := man.Docs.GetCommand("auth/client-credentials",
		man.WithRun(auth_clientCredentials),
	)
	clientCredentialsCmd.Flags().String(
		clientCredentialsCmd.GetDocFlag("client-id").Name,
		clientCredentialsCmd.GetDocFlag("client-id").Default,
		clientCredentialsCmd.GetDocFlag("client-id").Description,
	)
	clientCredentialsCmd.Flags().String(
		clientCredentialsCmd.GetDocFlag("client-secret").Name,
		clientCredentialsCmd.GetDocFlag("client-secret").Default,
		clientCredentialsCmd.GetDocFlag("client-secret").Description,
	)
}
