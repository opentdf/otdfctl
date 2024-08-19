package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var clientCredentialsCmd = man.Docs.GetCommand("auth/client-credentials",
	man.WithRun(auth_clientCredentials),
	man.WithHiddenFlags("with-client-creds", "with-client-creds-file"),
)

func auth_clientCredentials(cmd *cobra.Command, args []string) {
	var c handlers.ClientCredentials

	flagHelper := cli.NewFlagHelper(cmd)
	host := flagHelper.GetRequiredString("host")
	tlsNoVerify := flagHelper.GetOptionalBool("tls-no-verify")

	p := cli.NewPrinter(true)

	if len(args) > 0 {
		c.ClientId = args[0]
	}
	if len(args) > 1 {
		c.ClientSecret = args[1]
	}

	if c.ClientId == "" {
		c.ClientId = cli.AskForInput("Enter client id: ")
	}
	if c.ClientSecret == "" {
		c.ClientSecret = cli.AskForSecret("Enter client secret: ")
	}

	p.Printf("Logging in with client ID and secret for %s... ", host)
	if _, err := handlers.GetTokenWithClientCreds(cmd.Context(), host, c, tlsNoVerify); err != nil {
		fmt.Println("failed")
		cli.ExitWithError("An error occurred during login. Please check your credentials and try again", err)
	}
	p.Println("ok")

	p.Print("Storing client ID and secret in keyring... ")
	// if err := handlers.NewKeyring(host).SetClientCredentials(c); err != nil {
	// 	fmt.Println("failed")
	// 	cli.ExitWithError("Failed to cache client credentials", err)
	// }
	p.Println("ok")
}

func init() {
	authCmd.AddCommand(&clientCredentialsCmd.Command)
}
