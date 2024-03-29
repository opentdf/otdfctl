package cmd

import (
	"fmt"
	"strings"

	"github.com/opentdf/tructl/pkg/cli"
	"github.com/opentdf/tructl/pkg/handlers"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

var auth_loginCommands = []string{
	// auth_loginPassword.Use,
	auth_loginClientCredentials.Use,
}

var auth_loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Allows you to login in, in several ways [" + strings.Join(auth_loginCommands, ", ") + "]",
	Long: `
Auth - Login - Allows you to login in via all of the supported OAuth2 methods`,
}

var auth_loginClientCredentials = &cobra.Command{
	Use:   "clientCredentials",
	Short: "Allows the user to login in via clientId and clientSecret. This will subsequently be stored in the OS-specific keychain by default, but can be disabled with the --no-cache flag.",
	Run: func(cmd *cobra.Command, args []string) {
		h := cli.NewHandler(cmd)
		defer h.Close()

		flagHelper := cli.NewFlagHelper(cmd)
		clientId := flagHelper.GetOptionalString("clientId")
		clientSecret := flagHelper.GetOptionalString("clientSecret")
		// noCache := flagHelper.GetOptionalString("noCache")
		errMsg := fmt.Sprintf("Please provide required flag: (%s)", "Param Not Found")

		// h.DEBUG_PrintKeyRingSecrets()

		// check if we have a clientId in the keyring, if a null value is passed in
		if clientId == "" {
			fmt.Println("No clientId provided. Attempting to retrieve the default from keyring.")
			retrievedClientID, errID := keyring.Get(handlers.TOKEN_URL, handlers.TRUCTL_CLIENT_ID_CACHE_KEY)
			if errID == nil {
				clientId = retrievedClientID
				fmt.Println(cli.SuccessMessage("Retrieved stored clientId from keyring"))
			}
		}

		// now lets check if we still don't have it, and if not, throw and error
		if clientId == "" {
			errMsg = fmt.Sprintf("Please provide required flag: (%s)", "clientId")
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
			errMsg = fmt.Sprintf("Please provide required flag: (%s)", "clientSecret")
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

		fmt.Println(cli.SuccessMessage("Successfully logged in with clientId and clientSecret"))
	},
}

func init() {
	auth_loginCmd.AddCommand(auth_loginClientCredentials)
	auth_loginClientCredentials.Flags().StringP("clientId", "i", "", "The client id")
	auth_loginClientCredentials.Flags().StringP("clientSecret", "s", "", "The client secret")
	authCmd.AddCommand(auth_loginCmd)
}
