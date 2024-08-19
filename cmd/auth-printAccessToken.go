package cmd

import (
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var auth_printAccessTokenCmd = man.Docs.GetCommand("auth/print-access-token",
	man.WithRun(auth_printAccessToken))

func auth_printAccessToken(cmd *cobra.Command, args []string) {
	flagHelper := cli.NewFlagHelper(cmd)
	host := flagHelper.GetRequiredString("host")
	jsonOut := flagHelper.GetOptionalBool("json")

	printEnabled := !jsonOut
	p := cli.NewPrinter(printEnabled)

	p.Printf("Getting stored client credentials for %s... ", host)
	// clientCredentials, err := handlers.NewKeyring(host).GetClientCredentials()
	// if err != nil {
	// 	p.Println("failed")
	// 	cli.ExitWithError("Client credentials not found. Please use `auth client-credentials` to set them", err)
	// }
	// p.Println("ok")

	// p.Printf("Getting access token for %s... ", clientCredentials.ClientId)
	// tok, err := handlers.GetTokenWithClientCreds(
	// 	context.Background(),
	// 	host,
	// 	clientCredentials,
	// 	flagHelper.GetOptionalBool("tls-no-verify"),
	// )
	// if err != nil {
	// 	p.Println("failed")
	// 	cli.ExitWithError("Failed to get token", err)
	// }
	// p.Println("ok")
	// p.Printf("Access Token: %s\n", tok.AccessToken)

	// if jsonOut {
	// 	d, err := json.MarshalIndent(tok, "", "  ")
	// 	if err != nil {
	// 		cli.ExitWithError("Failed to marshal token to json", err)
	// 	}

	// 	fmt.Println(string(d))
	// 	return
	// }
}

func init() {
	auth_printAccessTokenCmd.Flags().Bool(
		auth_printAccessTokenCmd.GetDocFlag("json").Name,
		auth_printAccessTokenCmd.GetDocFlag("json").DefaultAsBool(),
		auth_printAccessTokenCmd.GetDocFlag("json").Description,
	)

	authCmd.AddCommand(&auth_printAccessTokenCmd.Command)
}
