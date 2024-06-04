package cmd

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var (
	clientCredentialsCmd = man.Docs.GetCommand("auth/client-credentials",
		man.WithRun(auth_clientCredentials),
	)
	noCacheCreds bool
)

func auth_clientCredentials(cmd *cobra.Command, args []string) {
	var err error

	flagHelper := cli.NewFlagHelper(cmd)
	host := flagHelper.GetRequiredString("host")
	tlsNoVerify := flagHelper.GetOptionalBool("tls-no-verify")
	clientID := flagHelper.GetOptionalString("client-id")
	clientSecret := flagHelper.GetOptionalString("client-secret")

	slog.Debug("Checking for client credentials file", slog.String("with-client-creds-file", clientCredsFile))
	if clientCredsFile != "" {
		creds, err := handlers.GetClientCredsFromFile(clientCredsFile)
		if err != nil {
			cli.ExitWithError("Failed to parse client credentials JSON", err)
		}
		clientID = creds.ClientID
		clientSecret = creds.ClientSecret
	}

	// if not provided by flag, check keyring cache for clientID
	if clientID == "" {
		slog.Debug("No client-id provided. Attempting to retrieve the default from keyring.")
		clientID, err = handlers.GetClientIDFromCache(host)
		if err != nil || clientID == "" {
			cli.ExitWithError("Please provide required flag: (client-id)", errors.New("no client-id found"))
		} else {
			slog.Debug(cli.SuccessMessage("Retrieved stored client-id from keyring"))
		}
	}

	// check if we have a clientSecret in the keyring, if a null value is passed in
	if clientSecret == "" {
		clientSecret, err = handlers.GetClientSecretFromCache(host, clientID)
		if err == nil || clientSecret == "" {
			cli.ExitWithError("Please provide required flag: (client-secret)", errors.New("no client-secret found"))
		} else {
			slog.Debug("Retrieved stored client-secret from keyring")
		}
	}

	slog.Debug("Attempting to login with client credentials", slog.String("client-id", clientID))
	if err := handlers.GetTokenWithClientCreds(cmd.Context(), host, clientID, clientSecret, tlsNoVerify); err != nil {
		cli.ExitWithError("An error occurred during login. Please check your credentials and try again", err)
	}

	fmt.Println(cli.SuccessMessage("Successfully logged in with client ID and secret"))
}

func init() {
	clientCredentialsCmd := man.Docs.GetCommand("auth/client-credentials",
		man.WithRun(auth_clientCredentials),
		// use the individual client-id and client-secret flags here instead of the global with-client-creds flag
		man.WithHiddenFlags("with-client-creds"),
	)
	clientCredentialsCmd.Flags().StringP(
		clientCredentialsCmd.GetDocFlag("client-id").Name,
		clientCredentialsCmd.GetDocFlag("client-id").Shorthand,
		clientCredentialsCmd.GetDocFlag("client-id").Default,
		clientCredentialsCmd.GetDocFlag("client-id").Description,
	)
	clientCredentialsCmd.Flags().StringP(
		clientCredentialsCmd.GetDocFlag("client-secret").Name,
		clientCredentialsCmd.GetDocFlag("client-secret").Shorthand,
		clientCredentialsCmd.GetDocFlag("client-secret").Default,
		clientCredentialsCmd.GetDocFlag("client-secret").Description,
	)
	clientCredentialsCmd.Flags().BoolVarP(
		&noCacheCreds,
		clientCredentialsCmd.GetDocFlag("no-cache").Name,
		clientCredentialsCmd.GetDocFlag("no-cache").Shorthand,
		clientCredentialsCmd.GetDocFlag("no-cache").DefaultAsBool(),
		clientCredentialsCmd.GetDocFlag("no-cache").Description,
	)
}
