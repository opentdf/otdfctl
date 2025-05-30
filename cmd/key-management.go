package cmd

import (
	"github.com/opentdf/otdfctl/pkg/man"
)

// KeyCmd is the command for managing keys
var keyMngmtCmd = man.Docs.GetCommand("key-management")

func init() {
	keyMngmtCmd.PersistentFlags().BoolVar(
		&configFlagOverrides.OutputFormatJSON,
		keyMngmtCmd.GetDocFlag("json").Name,
		keyMngmtCmd.GetDocFlag("json").DefaultAsBool(),
		keyMngmtCmd.GetDocFlag("json").Description,
	)
	RootCmd.AddCommand(&keyMngmtCmd.Command)
}
