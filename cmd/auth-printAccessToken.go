package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var printAccessToken = man.Docs.GetCommand("auth/print-access-token",
	man.WithRun(auth_printAccessToken),
)

func auth_printAccessToken(cmd *cobra.Command, args []string) {
	tok, err := handlers.GetOIDCTokenFromCache()
	if err != nil {
		cli.ExitWithError("Failed to get OIDC token from cache", err)
	}
	fmt.Print(tok)
}
