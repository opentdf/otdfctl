package cmd

import (
	"github.com/opentdf/otdfctl/pkg/man"
)

func init() {
	cmd := man.Docs.GetCommand("auth",
		man.WithSubcommands(clientCredentialsCmd),
		man.WithSubcommands(printAccessToken),
		man.WithSubcommands(clearCachedCredsCmd),
	)
	RootCmd.AddCommand(&cmd.Command)
}
