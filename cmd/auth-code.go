package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var noCacheCreds bool

func auth_codeLogin(cmd *cobra.Command, args []string) {
	flagHelper := cli.NewFlagHelper(cmd)
	host := flagHelper.GetRequiredString("host")
	clientID := flagHelper.GetOptionalString("client-id")
	tlsNoVerify := flagHelper.GetOptionalBool("tls-no-verify")

	printer := cli.NewPrinter(!noCacheCreds)

	tok, err := handlers.LoginWithPKCE(host, clientID, tlsNoVerify, noCacheCreds)
	if err != nil {
		cli.ExitWithError("could not authenticate", err)
	}
	if noCacheCreds {
		fmt.Print(tok.AccessToken)
	}
	// TODO: set to the keyring/profile here
	printer.Println(cli.SuccessMessage("Successfully logged in with auth code PKCE flow. Credentials cached on native OS."))
}

var codeLoginCmd *man.Doc

func init() {
	codeLoginCmd = man.Docs.GetCommand("auth/code-login",
		man.WithRun(auth_codeLogin),
	)
	codeLoginCmd.Flags().StringP(
		codeLoginCmd.GetDocFlag("client-id").Name,
		codeLoginCmd.GetDocFlag("client-id").Shorthand,
		codeLoginCmd.GetDocFlag("client-id").Default,
		codeLoginCmd.GetDocFlag("client-id").Description,
	)
	codeLoginCmd.Flags().BoolVarP(
		&noCacheCreds,
		codeLoginCmd.GetDocFlag("no-cache").Name,
		codeLoginCmd.GetDocFlag("no-cache").Shorthand,
		codeLoginCmd.GetDocFlag("no-cache").DefaultAsBool(),
		codeLoginCmd.GetDocFlag("no-cache").Description,
	)
	authCmd.AddCommand(&codeLoginCmd.Command)
}
