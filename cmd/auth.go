package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

func auth_clientCredentials(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	clientId := flagHelper.GetOptionalString("client-id")
	clientSecret := flagHelper.GetOptionalString("client-secret")
	// noCache := flagHelper.GetOptionalString("noCache")
	errMsg := fmt.Sprintf("Please provide required flag: (%s)", "Param Not Found")

	// h.DEBUG_PrintKeyRingSecrets()

	// check if we have a clientId in the keyring, if a null value is passed in
	if clientId == "" {
		fmt.Println("No client-id provided. Attempting to retrieve the default from keyring.")
		retrievedClientID, errID := keyring.Get(handlers.TOKEN_URL, handlers.OTDFCTL_CLIENT_ID_CACHE_KEY)
		if errID == nil {
			clientId = retrievedClientID
			fmt.Println(cli.SuccessMessage("Retrieved stored clientId from keyring"))
		}
	}

	// now lets check if we still don't have it, and if not, throw and error
	if clientId == "" {
		errMsg = fmt.Sprintf("Please provide required flag: (%s)", "client-id")
		fmt.Println(cli.ErrorMessage(errMsg, nil))
		cli.ExitWithError("Failed to create attribute", nil)
		return
	}

	// check if we have a clientSecret in the keyring, if a null value is passed in
	if clientSecret == "" {
		retrievedSecret, krErr := keyring.Get(handlers.TOKEN_URL, clientId)
		if krErr == nil {
			clientSecret = retrievedSecret
			fmt.Println(cli.SuccessMessage("Retrieved stored clientSecret from keyring"))
		}
	}
	// check if we still don't have it, and if not throw an error
	if clientSecret == "" {
		errMsg = fmt.Sprintf("Please provide required flag: (%s)", "client-secret")
		fmt.Println(cli.ErrorMessage(errMsg, nil))
		cli.ExitWithError("Failed to create attribute", nil)
		return
	}

	// for now we're hardcoding the TOKEN_URL as a constant at the top
	_, err := h.GetTokenWithClientCredentials(clientId, clientSecret, handlers.TOKEN_URL, false)
	if err != nil {
		errMsg = cli.ErrorMessage("An error occurred during login. Please check your credentials and try again.", nil)
		fmt.Println(errMsg)
		cli.ExitWithError(errMsg, err)
		return
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

	cmd := man.Docs.GetCommand("auth",
		man.WithSubcommands(clientCredentialsCmd),
	)
	rootCmd.AddCommand(&cmd.Command)
}
