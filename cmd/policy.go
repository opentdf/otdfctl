package cmd

import (
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

// PolicyCmd is the command for managing policies
var policyCmd = &cobra.Command{
	Use:   man.Docs.GetDoc("policy").Use,
	Short: man.Docs.GetDoc("policy").Short,
	Long:  man.Docs.GetDoc("policy").Long,
}

func init() {
	doc := man.Docs.GetDoc("policy")
	policyCmd.PersistentFlags().BoolVar(
		&configFlagOverrides.OutputFormatJSON,
		doc.GetDocFlag("json").Name,
		doc.GetDocFlag("json").DefaultAsBool(),
		doc.GetDocFlag("json").Description,
	)
	RootCmd.AddCommand(policyCmd)
}
