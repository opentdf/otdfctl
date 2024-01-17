package cmd

import (
	"strings"

	"github.com/spf13/cobra"
)

var (
	resourceEncodingCmds = []string{
		resourceEncodingsListCmd.Use,
	}

	resourceEncodingsCmd = &cobra.Command{
		Use:   "resource-encodings",
		Short: "Manage resource encodings with subcommands [" + strings.Join(resourceEncodingCmds, ", ") + "]",
		Long: `Resource encodings

Resource encodings are used to encode resources with an....
`,
	}

	resourceEncodingCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create resource encodings",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	resourceEncodingGetCmd = &cobra.Command{
		Use:   "get <id>",
		Short: "Get resource encodings",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	resourceEncodingsListCmd = &cobra.Command{
		Use:   "list",
		Short: "List resource encodings",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	resourceEncodingUpdateCmd = &cobra.Command{
		Use:   "update <id>",
		Short: "Update resource encodings",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	resourceEncodingDeleteCmd = &cobra.Command{
		Use:   "delete <id>",
		Short: "Delete resource encodings",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
)

func init() {
	rootCmd.AddCommand(resourceEncodingsCmd)

	resourceEncodingsCmd.AddCommand(resourceEncodingCreateCmd)

	resourceEncodingsCmd.AddCommand(resourceEncodingGetCmd)

	resourceEncodingsCmd.AddCommand(resourceEncodingsListCmd)

	resourceEncodingsCmd.AddCommand(resourceEncodingUpdateCmd)

	resourceEncodingsCmd.AddCommand(resourceEncodingDeleteCmd)
}
