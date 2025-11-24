package auth

import (
	"github.com/opentdf/otdfctl/pkg/man"
)

var clearClientCredentialsCmd = man.Docs.GetCommand("auth/clear-client-credentials")

func init() {
	Cmd.AddCommand(&clearClientCredentialsCmd.Command)
}
