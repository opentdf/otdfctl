package cmd

import (
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

// KeyCmd is the command for managing keys
var keyCmd = &cobra.Command{
	Use:   man.Docs.GetDoc("key").Use,
	Short: man.Docs.GetDoc("key").Short,
	Long:  man.Docs.GetDoc("key").Long,
}

func init() {
	doc := man.Docs.GetDoc("key")
	keyCmd.PersistentFlags().BoolVar(
		&configFlagOverrides.OutputFormatJSON,
		doc.GetDocFlag("json").Name,
		doc.GetDocFlag("json").DefaultAsBool(),
		doc.GetDocFlag("json").Description,
	)
	RootCmd.AddCommand(keyCmd)
}
