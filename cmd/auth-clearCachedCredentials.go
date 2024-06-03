package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var clearCachedCredsCmd = man.Docs.GetCommand("auth/clear-cached-credentials",
	man.WithRun(auth_clearCreds),
	man.WithHiddenFlags("with-client-creds", "with-client-creds-file"),
)

func auth_clearCreds(cmd *cobra.Command, args []string) {
	flagHelper := cli.NewFlagHelper(cmd)
	host := flagHelper.GetRequiredString("host")

	if err := handlers.ClearCachedCredentials(host); err != nil {
		cli.ExitWithError("Failed to clear cached client credentials and token", err)
	}

	fmt.Println(cli.SuccessMessage("Cached client credentials and token are clear."))
}
