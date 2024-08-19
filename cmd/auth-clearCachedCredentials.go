package cmd

import (
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var auth_clearClientCredentialsCmd = man.Docs.GetCommand("auth/clear-client-credentials",
	man.WithRun(auth_clearCreds),
	man.WithHiddenFlags("with-client-creds", "with-client-creds-file"),
)

func auth_clearCreds(cmd *cobra.Command, args []string) {
	flagHelper := cli.NewFlagHelper(cmd)
	host := flagHelper.GetRequiredString("host")

	p := cli.NewPrinter(true)

	p.Printf("Clearing cached client credentials for %s... ", host)
	// if err := handlers.NewKeyring(host).DeleteClientCredentials(); err != nil {
	// 	fmt.Println("failed")
	// 	cli.ExitWithError("Failed to clear cached client credentials", err)
	// }
	p.Println("ok")
}

func init() {
	auth_clearClientCredentialsCmd.Flags().String(
		auth_clearClientCredentialsCmd.GetDocFlag("all").Name,
		auth_clearClientCredentialsCmd.GetDocFlag("all").Description,
		auth_clearClientCredentialsCmd.GetDocFlag("all").Default,
	)

	authCmd.AddCommand(&auth_clearClientCredentialsCmd.Command)
}
