package cmd

import (
	_ "embed"
	"strings"

	"github.com/opentdf/tructl/docs/man"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/tructl/issues/73] is addressed

var (
	policy_resource_mappingsTerms []string

	policy_resource_mappingsCmd = &cobra.Command{
		Use:     man.PolicyResourceMappings["en"].Command,
		Aliases: man.PolicyResourceMappings["en"].Aliases,
		Short: man.PolicyResourceMappings["en"].ShortWithSubCommands([]string{
			policy_resource_mappingsCreateCmd.Use,
			policy_resource_mappingsGetCmd.Use,
			policy_resource_mappingsListCmd.Use,
			policy_resource_mappingsUpdateCmd.Use,
			policy_resource_mappingsDeleteCmd.Use,
		}),
		Long: man.PolicyResourceMappings["en"].Long,
	}

	policy_resource_mappingsCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create resource mappings",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			attrId := flagHelper.GetRequiredString("attribute-value-id")
			terms := flagHelper.GetStringSlice("terms", policy_resource_mappingsTerms, cli.FlagHelperStringSliceOptions{
				Min: 1,
			})
			metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

			resourceMapping, err := h.CreateResourceMapping(attrId, terms, getMetadataMutable(metadataLabels))
			if err != nil {
				cli.ExitWithError("Failed to create resource mapping", err)
			}

			t := cli.NewTabular().Rows([][]string{
				{"Id", resourceMapping.Id},
				{"Attribute Value Id", resourceMapping.AttributeValue.Id},
				{"Attribute Value", resourceMapping.AttributeValue.Value},
				{"Terms", strings.Join(resourceMapping.Terms, ", ")},
			}...)
			HandleSuccess(cmd, resourceMapping.Id, t, resourceMapping)
		},
	}

	policy_resource_mappingsGetCmd = &cobra.Command{
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

			t := cli.NewTabular().Rows([][]string{
				{"Id", resourceMapping.Id},
				{"Attribute Value Id", resourceMapping.AttributeValue.Id},
				{"Attribute Value", resourceMapping.AttributeValue.Value},
				{"Terms", strings.Join(resourceMapping.Terms, ", ")},
			}...)
			HandleSuccess(cmd, resourceMapping.Id, t, resourceMapping)
		},
	}

	policy_resource_mappingsListCmd = &cobra.Command{
		Use:   "list",
		Short: "List resource mappings",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			rmList, err := h.ListResourceMappings()
			if err != nil {
				cli.ExitWithError("Failed to list resource mappings", err)
			}

			t := cli.NewTable()
			t.Headers("Id", "Attribute Value Id", "Attribute Value", "Terms")
			for _, resourceMapping := range rmList {
				t.Row(resourceMapping.Id, resourceMapping.AttributeValue.Id, resourceMapping.AttributeValue.Value, strings.Join(resourceMapping.Terms, ", "))
			}
			HandleSuccess(cmd, "", t, rmList)
		},
	}

	policy_resource_mappingsUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update resource mappings",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")
			attrValueId := flagHelper.GetOptionalString("attribute-value-id")
			terms := flagHelper.GetStringSlice("terms", policy_resource_mappingsTerms, cli.FlagHelperStringSliceOptions{})
			labels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

			resourceMapping, err := h.UpdateResourceMapping(id, attrValueId, terms, getMetadataMutable(labels), getMetadataUpdateBehavior())
			if err != nil {
				cli.ExitWithError("Failed to update resource mapping", err)
			}

			t := cli.NewTabular().Rows([][]string{
				{"Id", resourceMapping.Id},
				{"Attribute Value Id", resourceMapping.AttributeValue.Id},
				{"Attribute Value", resourceMapping.AttributeValue.Value},
				{"Terms", strings.Join(resourceMapping.Terms, ", ")},
			}...)
			HandleSuccess(cmd, resourceMapping.Id, t, resourceMapping)
		},
	}

	policy_resource_mappingsDeleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "Delete resource mappings",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			cli.ConfirmAction(cli.ActionDelete, "resource-mapping", id)

			resourceMapping, err := h.DeleteResourceMapping(id)
			if err != nil {
				cli.ExitWithError("Failed to delete resource mapping", err)
			}

			t := cli.NewTabular().Rows([][]string{
				{"Id", resourceMapping.Id},
				{"Attribute Value Id", resourceMapping.AttributeValue.Id},
				{"Attribute Value", resourceMapping.AttributeValue.Value},
				{"Terms", strings.Join(resourceMapping.Terms, ", ")},
			}...)
			HandleSuccess(cmd, resourceMapping.Id, t, resourceMapping)
		},
	}
)

func init() {
	policyCmd.AddCommand(policy_resource_mappingsCmd)

	policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsCreateCmd)
	policy_resource_mappingsCreateCmd.Flags().String("attribute-value-id", "", "Attribute Value ID")
	policy_resource_mappingsCreateCmd.Flags().StringSliceVar(&policy_resource_mappingsTerms, "terms", []string{}, "Synonym terms")
	injectLabelFlags(policy_resource_mappingsCreateCmd, false)

	policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsGetCmd)
	policy_resource_mappingsGetCmd.Flags().String("id", "", "Resource Mapping ID")

	policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsListCmd)

	policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsUpdateCmd)
	policy_resource_mappingsUpdateCmd.Flags().String("id", "", "Resource Mapping ID")
	policy_resource_mappingsUpdateCmd.Flags().String("attribute-value-id", "", "Attribute Value ID")
	policy_resource_mappingsUpdateCmd.Flags().StringSliceVar(&policy_resource_mappingsTerms, "terms", []string{}, "Synonym terms")
	injectLabelFlags(policy_resource_mappingsUpdateCmd, true)

	policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsDeleteCmd)
	policy_resource_mappingsDeleteCmd.Flags().String("id", "", "Resource Mapping ID")
}
