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
	policy_resourceMappingsCmd    *cobra.Command
)

func policy_createResourceMapping(cmd *cobra.Command, args []string) {
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
	rows := [][]string{
		{"Id", resourceMapping.Id},
		{"Attribute Value Id", resourceMapping.AttributeValue.Id},
		{"Attribute Value", resourceMapping.AttributeValue.Value},
		{"Terms", strings.Join(resourceMapping.Terms, ", ")},
	}
	if mdRows := getMetadataRows(resourceMapping.Metadata); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular().Rows(rows...)
	HandleSuccess(cmd, resourceMapping.Id, t, resourceMapping)
}

func policy_getResourceMapping(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	resourceMapping, err := h.GetResourceMapping(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get resource mapping (%s)", id), err)
	}
	rows := [][]string{
		{"Id", resourceMapping.Id},
		{"Attribute Value Id", resourceMapping.AttributeValue.Id},
		{"Attribute Value", resourceMapping.AttributeValue.Value},
		{"Terms", strings.Join(resourceMapping.Terms, ", ")},
	}
	if mdRows := getMetadataRows(resourceMapping.Metadata); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular().Rows(rows...)
	HandleSuccess(cmd, resourceMapping.Id, t, resourceMapping)
}

func policy_listResourceMappings(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	rmList, err := h.ListResourceMappings()
	if err != nil {
		cli.ExitWithError("Failed to list resource mappings", err)
	}

	t := cli.NewTable()
	t.Headers("Id", "Attribute Value Id", "Attribute Value", "Terms", "Metadata.Labels", "Metadata.CreatedAt", "Metadata.UpdatedAt")
	for _, resourceMapping := range rmList {
		metadata := cli.ConstructMetadata(resourceMapping.Metadata)
		t.Row(resourceMapping.Id, resourceMapping.AttributeValue.Id, resourceMapping.AttributeValue.Value, strings.Join(resourceMapping.Terms, ", "), metadata["Labels"], metadata["Created At"], metadata["Updated At"])
	}
	HandleSuccess(cmd, "", t, rmList)
}

func policy_updateResourceMapping(cmd *cobra.Command, args []string) {
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
	rows := [][]string{
		{"Id", resourceMapping.Id},
		{"Attribute Value Id", resourceMapping.AttributeValue.Id},
		{"Attribute Value", resourceMapping.AttributeValue.Value},
		{"Terms", strings.Join(resourceMapping.Terms, ", ")},
	}
	if mdRows := getMetadataRows(resourceMapping.Metadata); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular().Rows(rows...)
	HandleSuccess(cmd, resourceMapping.Id, t, resourceMapping)
}

func policy_deleteResourceMapping(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	cli.ConfirmAction(cli.ActionDelete, "resource-mapping", id)

	resourceMapping, err := h.DeleteResourceMapping(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to delete resource mapping (%s)", id), err)
	}
	rows := [][]string{
		{"Id", resourceMapping.Id},
		{"Attribute Value Id", resourceMapping.AttributeValue.Id},
		{"Attribute Value", resourceMapping.AttributeValue.Value},
		{"Terms", strings.Join(resourceMapping.Terms, ", ")},
	}
	t := cli.NewTabular().Rows(rows...)
	HandleSuccess(cmd, resourceMapping.Id, t, resourceMapping)
}

func init() {
	createDoc := man.Docs.GetCommand("policy/resource-mappings/create",
		man.WithRun(policy_createResourceMapping),
	)
	createDoc.Flags().String(
		createDoc.GetDocFlag("attribute-value-id").Name,
		createDoc.GetDocFlag("attribute-value-id").Default,
		createDoc.GetDocFlag("attribute-value-id").Description,
	)
	createDoc.Flags().StringSliceVar(
		&policy_resource_mappingsTerms,
		createDoc.GetDocFlag("terms").Name,
		[]string{},
		createDoc.GetDocFlag("terms").Description,
	)
	injectLabelFlags(&createDoc.Command, false)

	getDoc := man.Docs.GetCommand("policy/resource-mappings/get",
		man.WithRun(policy_getResourceMapping),
	)
	getDoc.Flags().String(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)

	listDoc := man.Docs.GetCommand("policy/resource-mappings/list",
		man.WithRun(policy_listResourceMappings),
	)

	updateDoc := man.Docs.GetCommand("policy/resource-mappings/update",
		man.WithRun(policy_updateResourceMapping),
	)
	updateDoc.Flags().String(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	updateDoc.Flags().String(
		updateDoc.GetDocFlag("attribute-value-id").Name,
		updateDoc.GetDocFlag("attribute-value-id").Default,
		updateDoc.GetDocFlag("attribute-value-id").Description,
	)
	updateDoc.Flags().StringSliceVar(
		&policy_resource_mappingsTerms,
		updateDoc.GetDocFlag("terms").Name,
		[]string{},
		updateDoc.GetDocFlag("terms").Description,
	)
	injectLabelFlags(&updateDoc.Command, true)

	deleteDoc := man.Docs.GetCommand("policy/resource-mappings/delete",
		man.WithRun(policy_deleteResourceMapping),
	)
	deleteDoc.Flags().String(
		deleteDoc.GetDocFlag("id").Name,
		deleteDoc.GetDocFlag("id").Default,
		deleteDoc.GetDocFlag("id").Description,
	)

	doc := man.Docs.GetCommand("policy/resource-mappings",
		man.WithSubcommands(createDoc, getDoc, listDoc, updateDoc, deleteDoc),
	)
	policy_resourceMappingsCmd = &doc.Command
	policyCmd.AddCommand(policy_resourceMappingsCmd)
}
