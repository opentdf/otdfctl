package cmd

import (
	"github.com/opentdf/otdfctl/pkg/man"
)

func init() {
	cmd := man.Docs.GetCommand("auth",
		man.WithSubcommands(clientCredentialsCmd),
		man.WithSubcommands(printAccessTokenCmd),
		man.WithSubcommands(clearCachedCredsCmd),
		man.WithSubcommands(codeLoginCmd),
	)
	RootCmd.AddCommand(&cmd.Command)
}
