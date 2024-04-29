package cmd

import (
	"github.com/opentdf/otdfctl/pkg/man"
)

func init() {
	cmd := man.Docs.GetCommand("auth",
		man.WithSubcommands(clientCredentialsCmd),
		man.WithSubcommands(printAccessToken),
	)
	RootCmd.AddCommand(&cmd.Command)
}
