package policy

import (
	"github.com/opentdf/otdfctl/pkg/man"
)

// KeyCmd is the command for managing keys
var KeyManagementCmd = man.Docs.GetCommand("policy/key-management")

// initKeyManagementCommands sets up the key-management command.
func initKeyManagementCommands() {
	Cmd.AddCommand(&KeyManagementCmd.Command)
}
