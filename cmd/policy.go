package cmd

import "github.com/spf13/cobra"

var (
	// PolicyCmd is the command for managing policies
	policyCmd = &cobra.Command{
		Use:   "policy",
		Short: "Manage policies",
		Long: `
Manage policies within the platform.

Policy is a set of rules that are enforced by the platform. Specific to the the data centric
security, policy revolves around data attributes (referred to as attributes). Within the context
of attributes are namespaces, values, subject-mappings, resource-mappings, key-access-server grants,
and other key elements.
`,
	}
)

func init() {
	rootCmd.AddCommand(policyCmd)
}
