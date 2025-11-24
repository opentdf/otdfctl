package policy

import (
	"github.com/opentdf/otdfctl/pkg/man"
)

// KeyCmd is the command for managing keys
var keyMngmtCmd = man.Docs.GetCommand("policy/key-management")

func init() {
	Cmd.AddCommand(&keyMngmtCmd.Command)
}
