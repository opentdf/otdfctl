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

var clearCachedCredsCmd = man.Docs.GetCommand("auth/clear-cached-credentials",
	man.WithRun(auth_clearCreds),
)

func auth_clearCreds(cmd *cobra.Command, args []string) {
	cachedClientID, err := handlers.GetClientIDFromCache()
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			fmt.Println("No client-id found in the cache to clear.")
		} else {
			cli.ExitWithError("Failed to retrieve client id from keyring", err)
		}
	}

	// clear the client ID and secret from the keyring
	err = keyring.Delete(handlers.TOKEN_URL, cachedClientID)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			fmt.Println("No client secret found in the cache to clear under client-id: ", cachedClientID)
		} else {
			cli.ExitWithError("Failed to clear client secret from keyring", err)
		}
	}

	err = keyring.Delete(handlers.TOKEN_URL, handlers.OTDFCTL_CLIENT_ID_CACHE_KEY)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			fmt.Println("No client id found in the cache to clear.")
		} else {
			cli.ExitWithError("Failed to clear client id from keyring", err)
		}
	}

	err = keyring.Delete(handlers.TOKEN_URL, handlers.OTDFCTL_OIDC_TOKEN_KEY)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			fmt.Println("No token found in the cache to clear.")
		} else {
			cli.ExitWithError("Failed to clear token from keyring", err)
		}
	}

	fmt.Println(cli.SuccessMessage("Cached client credentials and token are clear."))
}
