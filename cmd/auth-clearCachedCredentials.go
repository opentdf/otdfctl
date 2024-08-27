package cmd

import (
	"github.com/opentdf/otdfctl/pkg/man"
)

var auth_clearClientCredentialsCmd = man.Docs.GetCommand("auth/clear-client-credentials")

func init() {
	authCmd.AddCommand(&auth_clearClientCredentialsCmd.Command)
}
