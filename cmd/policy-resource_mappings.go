package cmd

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/opentdf/tructl/docs/man"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

var (
	policy_resourceMappingsTerms []string

	policy_resourceMappingsCmd = &cobra.Command{
		Use:     man.PolicyResourceMappings["en"].Command,
		Aliases: man.PolicyResourceMappings["en"].Aliases,
		Short: man.PolicyResourceMappings["en"].ShortWithSubCommands([]string{
			policy_resourceMappingsCreateCmd.Use,
			policy_resourceMappingsGetCmd.Use,
			policy_resourceMappingsListCmd.Use,
			policy_resourceMappingsUpdateCmd.Use,
			policy_resourceMappingsDeleteCmd.Use,
		}),
		Long: man.PolicyResourceMappings["en"].Long,
	}

	policy_resourceMappingsCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create resource mappings",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			attrId := flagHelper.GetRequiredString("attribute-value-id")
			terms := flagHelper.GetStringSlice("terms", policy_resourceMappingsTerms, cli.FlagHelperStringSliceOptions{
				Min: 1,
			})

			resourceMapping, err := h.CreateResourceMapping(attrId, terms)
			if err != nil {
				cli.ExitWithError("Failed to create resource mapping", err)
			}

			fmt.Println(cli.SuccessMessage("Resource mapping created"))
			fmt.Println(cli.NewTabular().Rows([][]string{
				{"Id", resourceMapping.Id},
				{"Attribute Value Id", resourceMapping.AttributeValue.Id},
				{"Attribute Value", resourceMapping.AttributeValue.Value},
				{"Terms", strings.Join(resourceMapping.Terms, ", ")},
			}...).Render())
		},
	}

	policy_resourceMappingsGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get resource mappings",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			resourceMapping, err := h.GetResourceMapping(id)
			if err != nil {
				cli.ExitWithError("Failed to get resource mapping", err)
			}

			fmt.Println(cli.NewTabular().Rows([][]string{
				{"Id", resourceMapping.Id},
				{"Attribute Value Id", resourceMapping.AttributeValue.Id},
				{"Attribute Value", resourceMapping.AttributeValue.Value},
				{"Terms", strings.Join(resourceMapping.Terms, ", ")},
			}...).Render())
		},
	}

	policy_resourceMappingsListCmd = &cobra.Command{
		Use:   "list",
		Short: "List resource mappings",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			r, err := h.ListResourceMappings()
			if err != nil {
				cli.ExitWithError("Failed to list resource mappings", err)
			}

			t := cli.NewTable()
			t.Headers("Id", "Attribute Value Id", "Attribute Value", "Terms")
			for _, resourceMapping := range r {
				t.Row(resourceMapping.Id, resourceMapping.AttributeValue.Id, resourceMapping.AttributeValue.Value, strings.Join(resourceMapping.Terms, ", "))
			}
			fmt.Println(t.Render())
		},
	}

	policy_resourceMappingsUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update resource mappings",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")
			attrValueId := flagHelper.GetOptionalString("attribute-value-id")
			terms := flagHelper.GetStringSlice("terms", policy_resourceMappingsTerms, cli.FlagHelperStringSliceOptions{})

			resourceMapping, err := h.UpdateResourceMapping(id, attrValueId, terms)
			if err != nil {
				cli.ExitWithError("Failed to update resource mapping", err)
			}

			fmt.Println(cli.SuccessMessage("Resource mapping updated"))
			fmt.Println(cli.NewTabular().Rows([][]string{
				{"Id", resourceMapping.Id},
				{"Attribute Value Id", resourceMapping.AttributeValue.Id},
				{"Attribute Value", resourceMapping.AttributeValue.Value},
				{"Terms", strings.Join(resourceMapping.Terms, ", ")},
			}...).Render())
		},
	}

	policy_resourceMappingsDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete resource mappings",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			cli.ConfirmDelete("resource-mapping", id)

			resourceMapping, err := h.DeleteResourceMapping(id)
			if err != nil {
				cli.ExitWithError("Failed to delete resource mapping", err)
			}

			fmt.Println(cli.SuccessMessage("Resource mapping deleted"))
			fmt.Println(cli.NewTabular().Rows([][]string{
				{"Id", resourceMapping.Id},
				{"Attribute Value Id", resourceMapping.AttributeValue.Id},
				{"Attribute Value", resourceMapping.AttributeValue.Value},
				{"Terms", strings.Join(resourceMapping.Terms, ", ")},
			}...).Render())
		},
	}
)

func init() {
	policyCmd.AddCommand(policy_resourceMappingsCmd)

	policy_resourceMappingsCmd.AddCommand(policy_resourceMappingsCreateCmd)
	policy_resourceMappingsCreateCmd.Flags().String("attribute-value-id", "", "Attribute Value ID")
	policy_resourceMappingsCreateCmd.Flags().StringSliceVar(&policy_resourceMappingsTerms, "terms", []string{}, "Synonym terms")

	policy_resourceMappingsCmd.AddCommand(policy_resourceMappingsGetCmd)
	policy_resourceMappingsGetCmd.Flags().String("id", "", "Resource Mapping ID")

	policy_resourceMappingsCmd.AddCommand(policy_resourceMappingsListCmd)

	policy_resourceMappingsCmd.AddCommand(policy_resourceMappingsUpdateCmd)
	policy_resourceMappingsUpdateCmd.Flags().String("id", "", "Resource Mapping ID")
	policy_resourceMappingsUpdateCmd.Flags().String("attribute-value-id", "", "Attribute Value ID")
	policy_resourceMappingsUpdateCmd.Flags().StringSliceVar(&policy_resourceMappingsTerms, "terms", []string{}, "Synonym terms")

	policy_resourceMappingsCmd.AddCommand(policy_resourceMappingsDeleteCmd)
	policy_resourceMappingsDeleteCmd.Flags().String("id", "", "Resource Mapping ID")
}
