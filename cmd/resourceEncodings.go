package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	resourceEncodingTerms []string

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
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			attrId := flagHelper.GetRequiredString("attribute-id")
			attributeId, err := strconv.Atoi(attrId)
			if err != nil {
				cli.ExitWithError("Invalid attribute ID", err)
			}

			terms := flagHelper.GetStringSlice("terms", resourceEncodingTerms, cli.FlagHelperStringSliceOptions{
				Min: 1,
			})

			_, err = h.CreateResourceEncoding(attributeId, terms)
			if err != nil {
				cli.ExitWithError("Failed to create resource encoding", err)
			}

			fmt.Println(cli.SuccessMessage("Resource encoding created"))
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
	resourceEncodingCreateCmd.Flags().String("attribute-id", "", "Attribute ID")
	resourceEncodingCreateCmd.Flags().StringSliceVar(&resourceEncodingTerms, "terms", []string{}, "Synonym terms")

	resourceEncodingsCmd.AddCommand(resourceEncodingGetCmd)

	resourceEncodingsCmd.AddCommand(resourceEncodingsListCmd)

	resourceEncodingsCmd.AddCommand(resourceEncodingUpdateCmd)

	resourceEncodingsCmd.AddCommand(resourceEncodingDeleteCmd)
}
