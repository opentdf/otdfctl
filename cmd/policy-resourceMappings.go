package cmd

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/otdfctl/issues/73] is addressed

var (
	policy_resource_mappingsTerms []string

	//#region Resource Mappings
	policy_resource_mappingsCmd = &cobra.Command{
		Use:     man.Docs.GetDoc("policy/resource-mappings").Use,
		Aliases: man.Docs.GetDoc("policy/resource-mappings").Aliases,
		Short: man.Docs.GetDoc("policy/resource-mappings").GetShort([]string{
			policy_resource_mappingsCreateCmd.Use,
			policy_resource_mappingsGetCmd.Use,
			policy_resource_mappingsListCmd.Use,
			policy_resource_mappingsUpdateCmd.Use,
			policy_resource_mappingsDeleteCmd.Use,
		}),
		Long: man.Docs.GetDoc("policy/resource-mappings").Long,
	}
	//#endregion

	//#region Resource Mapping Create
	policy_resource_mappingsCreateCmd = &cobra.Command{
		Use:   man.Docs.GetDoc("policy/resource-mappings/create").Use,
		Short: man.Docs.GetDoc("policy/resource-mappings/create").Short,
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
	//#endregion

	policy_resource_mappingsGetCmd = &cobra.Command{
		Use:   man.Docs.GetDoc("policy/resource-mappings/get").Use,
		Short: man.Docs.GetDoc("policy/resource-mappings/get").Short,
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			resourceMapping, err := h.GetResourceMapping(id)
			if err != nil {
				cli.ExitWithError(fmt.Sprintf("Failed to get resource mapping (%s)", id), err)
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
		Use:   man.Docs.GetDoc("policy/resource-mappings/list").Use,
		Short: man.Docs.GetDoc("policy/resource-mappings/list").Short,
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
		Use:   man.Docs.GetDoc("policy/resource-mappings/update").Use,
		Short: man.Docs.GetDoc("policy/resource-mappings/update").Short,
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
				cli.ExitWithError(fmt.Sprintf("Failed to update resource mapping (%s)", id), err)
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
		Use:   man.Docs.GetDoc("policy/resource-mappings/delete").Use,
		Short: man.Docs.GetDoc("policy/resource-mappings/delete").Short,
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			cli.ConfirmAction(cli.ActionDelete, "resource-mapping", id)

			resourceMapping, err := h.DeleteResourceMapping(id)
			if err != nil {
				cli.ExitWithError(fmt.Sprintf("Failed to delete resource mapping (%s)", id), err)
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

	createDoc := man.Docs.GetDoc("policy/resource-mappings/create")
	policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsCreateCmd)
	policy_resource_mappingsCreateCmd.Flags().String(
		createDoc.GetDocFlag("attribute-value-id").Name,
		createDoc.GetDocFlag("attribute-value-id").Default,
		createDoc.GetDocFlag("attribute-value-id").Description,
	)
	policy_resource_mappingsCreateCmd.Flags().StringSliceVar(
		&policy_resource_mappingsTerms,
		createDoc.GetDocFlag("terms").Name,
		[]string{},
		createDoc.GetDocFlag("terms").Description,
	)
	injectLabelFlags(policy_resource_mappingsCreateCmd, false)

	getDoc := man.Docs.GetDoc("policy/resource-mappings/get")
	policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsGetCmd)
	policy_resource_mappingsGetCmd.Flags().String(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)

	policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsListCmd)

	updateDoc := man.Docs.GetDoc("policy/resource-mappings/update")
	policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsUpdateCmd)
	policy_resource_mappingsUpdateCmd.Flags().String(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	policy_resource_mappingsUpdateCmd.Flags().String(
		updateDoc.GetDocFlag("attribute-value-id").Name,
		updateDoc.GetDocFlag("attribute-value-id").Default,
		updateDoc.GetDocFlag("attribute-value-id").Description,
	)
	policy_resource_mappingsUpdateCmd.Flags().StringSliceVar(
		&policy_resource_mappingsTerms,
		updateDoc.GetDocFlag("terms").Name,
		[]string{},
		updateDoc.GetDocFlag("terms").Description,
	)
	injectLabelFlags(policy_resource_mappingsUpdateCmd, true)

	deleteDoc := man.Docs.GetDoc("policy/resource-mappings/delete")
	policy_resource_mappingsCmd.AddCommand(policy_resource_mappingsDeleteCmd)
	policy_resource_mappingsDeleteCmd.Flags().String(
		deleteDoc.GetDocFlag("id").Name,
		deleteDoc.GetDocFlag("id").Default,
		deleteDoc.GetDocFlag("id").Description,
	)
}
