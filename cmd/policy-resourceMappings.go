package cmd

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var (
	terms                      []string
	policy_resourceMappingsCmd *cobra.Command
)

func policy_createResourceMapping(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	attrId := c.Flags.GetRequiredID("attribute-value-id")
	grpID := c.Flags.GetOptionalID("group-id")
	terms = c.Flags.GetStringSlice("terms", terms, cli.FlagsStringSliceOptions{
		Min: 1,
	})
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	resourceMapping, err := h.CreateResourceMapping(attrId, terms, grpID, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create resource mapping", err)
	}
	rows := [][]string{
		{"Id", resourceMapping.GetId()},
		{"Attribute Value Id", resourceMapping.GetAttributeValue().GetId()},
		{"Attribute Value", resourceMapping.GetAttributeValue().GetValue()},
		{"Terms", strings.Join(resourceMapping.GetTerms(), ", ")},
		{"Group Id", resourceMapping.GetGroup().GetId()},
		{"Group Name", resourceMapping.GetGroup().GetName()},
	}
	if mdRows := getMetadataRows(resourceMapping.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, resourceMapping.GetId(), t, resourceMapping)
}

func policy_getResourceMapping(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	resourceMapping, err := h.GetResourceMapping(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get resource mapping (%s)", id), err)
	}
	rows := [][]string{
		{"Id", resourceMapping.GetId()},
		{"Attribute Value Id", resourceMapping.GetAttributeValue().GetId()},
		{"Attribute Value", resourceMapping.GetAttributeValue().GetValue()},
		{"Terms", strings.Join(resourceMapping.GetTerms(), ", ")},
		{"Group Id", resourceMapping.GetGroup().GetId()},
		{"Group Name", resourceMapping.GetGroup().GetName()},
	}
	if mdRows := getMetadataRows(resourceMapping.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, resourceMapping.GetId(), t, resourceMapping)
}

func policy_listResourceMappings(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")

	rmList, page, err := h.ListResourceMappings(cmd.Context(), limit, offset)
	if err != nil {
		cli.ExitWithError("Failed to list resource mappings", err)
	}

	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("attr_value_id", "Attribute Value Id", cli.FlexColumnWidthFour),
		table.NewFlexColumn("attr_value", "Attribute Value", cli.FlexColumnWidthFour),
		table.NewFlexColumn("terms", "Terms", cli.FlexColumnWidthThree),
		table.NewFlexColumn("group_id", "Group Id", cli.FlexColumnWidthFour),
		table.NewFlexColumn("group_name", "Group Name", cli.FlexColumnWidthFour),
		table.NewFlexColumn("labels", "Labels", cli.FlexColumnWidthOne),
		table.NewFlexColumn("created_at", "Created At", cli.FlexColumnWidthOne),
		table.NewFlexColumn("updated_at", "Updated At", cli.FlexColumnWidthOne),
	)
	rows := []table.Row{}
	for _, resourceMapping := range rmList {
		metadata := cli.ConstructMetadata(resourceMapping.GetMetadata())
		rows = append(rows, table.NewRow(table.RowData{
			"id":            resourceMapping.GetId(),
			"attr_value_id": resourceMapping.GetAttributeValue().GetId(),
			"attr_value":    resourceMapping.GetAttributeValue().GetValue(),
			"group_id":      resourceMapping.GetGroup().GetId(),
			"group_name":    resourceMapping.GetGroup().GetName(),
			"terms":         strings.Join(resourceMapping.GetTerms(), ", "),
			"labels":        metadata["Labels"],
			"created_at":    metadata["Created At"],
			"updated_at":    metadata["Updated At"],
		}))
	}
	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, rmList)
}

func policy_updateResourceMapping(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	attrValueId := c.Flags.GetOptionalID("attribute-value-id")
	grpID := c.Flags.GetOptionalID("group-id")
	terms = c.Flags.GetStringSlice("terms", terms, cli.FlagsStringSliceOptions{})
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	resourceMapping, err := h.UpdateResourceMapping(id, attrValueId, grpID, terms, getMetadataMutable(metadataLabels), getMetadataUpdateBehavior())
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to update resource mapping (%s)", id), err)
	}
	rows := [][]string{
		{"Id", resourceMapping.GetId()},
		{"Attribute Value Id", resourceMapping.GetAttributeValue().GetId()},
		{"Attribute Value", resourceMapping.GetAttributeValue().GetValue()},
		{"Terms", strings.Join(resourceMapping.GetTerms(), ", ")},
		{"Group Id", resourceMapping.GetGroup().GetId()},
		{"Group Name", resourceMapping.GetGroup().GetName()},
	}
	if mdRows := getMetadataRows(resourceMapping.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, resourceMapping.GetId(), t, resourceMapping)
}

func policy_deleteResourceMapping(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	force := c.Flags.GetOptionalBool("force")

	cli.ConfirmAction(cli.ActionDelete, "resource-mapping", id, force)

	resourceMapping, err := h.GetResourceMapping(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get resource mapping for delete (%s)", id), err)
	}

	_, err = h.DeleteResourceMapping(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to delete resource mapping (%s)", id), err)
	}
	rows := [][]string{
		{"Id", resourceMapping.GetId()},
		{"Attribute Value Id", resourceMapping.GetAttributeValue().GetId()},
		{"Attribute Value", resourceMapping.GetAttributeValue().GetValue()},
		{"Terms", strings.Join(resourceMapping.GetTerms(), ", ")},
		{"Group Id", resourceMapping.GetGroup().GetId()},
		{"Group Name", resourceMapping.GetGroup().GetName()},
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, resourceMapping.GetId(), t, resourceMapping)
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
		&terms,
		createDoc.GetDocFlag("terms").Name,
		[]string{},
		createDoc.GetDocFlag("terms").Description,
	)
	createDoc.Flags().String(
		createDoc.GetDocFlag("group-id").Name,
		createDoc.GetDocFlag("group-id").Default,
		createDoc.GetDocFlag("group-id").Description,
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
	injectListPaginationFlags(listDoc)

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
		&terms,
		updateDoc.GetDocFlag("terms").Name,
		[]string{},
		updateDoc.GetDocFlag("terms").Description,
	)
	updateDoc.Flags().String(
		updateDoc.GetDocFlag("group-id").Name,
		updateDoc.GetDocFlag("group-id").Default,
		updateDoc.GetDocFlag("group-id").Description,
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
	deleteDoc.Flags().Bool(
		deleteDoc.GetDocFlag("force").Name,
		false,
		deleteDoc.GetDocFlag("force").Description,
	)

	doc := man.Docs.GetCommand("policy/resource-mappings",
		man.WithSubcommands(createDoc, getDoc, listDoc, updateDoc, deleteDoc),
	)
	policy_resourceMappingsCmd = &doc.Command
	policyCmd.AddCommand(policy_resourceMappingsCmd)
}
