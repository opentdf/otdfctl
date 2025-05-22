package cmd

import (
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

// KeyCmd is the command for managing keys
var keyMngmtCmd = &cobra.Command{
	Use:   man.Docs.GetDoc("key-management").Use,
	Short: man.Docs.GetDoc("key-management").Short,
	Long:  man.Docs.GetDoc("key-management").Long,
}

func init() {
	doc := man.Docs.GetDoc("key-management")
	keyMngmtCmd.PersistentFlags().BoolVar(
		&configFlagOverrides.OutputFormatJSON,
		doc.GetDocFlag("json").Name,
		doc.GetDocFlag("json").DefaultAsBool(),
		doc.GetDocFlag("json").Description,
	)
	RootCmd.AddCommand(keyMngmtCmd)
}
