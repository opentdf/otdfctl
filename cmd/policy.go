package cmd

import (
	"github.com/opentdf/tructl/pkg/man"
	"github.com/spf13/cobra"
)

var (
	// PolicyCmd is the command for managing policies
	policyCmd = &cobra.Command{
		Use:   man.Docs.GetDoc("policy").Use,
		Short: man.Docs.GetDoc("policy").Short,
		Long:  man.Docs.GetDoc("policy").Long,
	}
)

func init() {
	rootCmd.AddCommand(policyCmd)
}
