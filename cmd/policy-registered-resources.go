package cmd

import (
	"fmt"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var (
	policy_registeredResourcesCmd = man.Docs.GetCommand("policy/registered-resources")
)

func policy_createRegisteredResource(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	name := c.Flags.GetRequiredString("name")
	metadataLabels := c.Flags.GetStringSlice("label", nil, cli.FlagsStringSliceOptions{})

	res, err := h.CreateRegisteredResource(name, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create registered resource", err)
	}

	t := cli.NewTabular([][]string{
		{"Id", res.GetId()},
		{"Name", res.GetName()},
	})
	HandleSuccess(cmd, res.GetId(), t, res)
}

func policy_getRegisteredResource(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	res, err := h.GetRegisteredResource(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get registered resource (%s)", id), err)
	}

	t := cli.NewTabular([][]string{
		{"Id", res.GetId()},
		{"Name", res.GetName()},
	})
	HandleSuccess(cmd, res.GetId(), t, res)
}

func policy_listRegisteredResources(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")

	resources, page, err := h.ListRegisteredResources(c.Context(), limit, offset)
	if err != nil {
		cli.ExitWithError("Failed to list registered resources", err)
	}

	t := cli.NewTable(
		table.NewUUIDColumn(),
		table.NewFlexColumn("name", "Name", cli.FlexColumnWidthThree),
	)
	rows := []table.Row{}
	for _, res := range resources {
		rows = append(rows, table.NewRow(table.RowData{
			"id":   res.GetId(),
			"name": res.GetName(),
		}))
	}
	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, resources)
}

func policy_updateRegisteredResource(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	name := c.Flags.GetOptionalString("name")
	metadataLabels := c.Flags.GetStringSlice("label", nil, cli.FlagsStringSliceOptions{})

	res, err := h.UpdateRegisteredResource(id, name, getMetadataMutable(metadataLabels), getMetadataUpdateBehavior())
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to update registered resource (%s)", id), err)
	}

	t := cli.NewTabular([][]string{
		{"Id", res.GetId()},
		{"Name", res.GetName()},
	})
	HandleSuccess(cmd, res.GetId(), t, res)
}

func policy_deleteRegisteredResource(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	res, err := h.DeleteRegisteredResource(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to delete registered resource (%s)", id), err)
	}

	t := cli.NewTabular([][]string{
		{"Deleted", "true"},
		{"Id", res.GetId()},
		{"Name", res.GetName()},
	})
	HandleSuccess(cmd, res.GetId(), t, res)
}

func init() {
	// Create a registered resource
	createDoc := man.Docs.GetCommand("policy/registered-resources/create",
		man.WithRun(policy_createRegisteredResource),
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("name").Name,
		createDoc.GetDocFlag("name").Shorthand,
		createDoc.GetDocFlag("name").Default,
		createDoc.GetDocFlag("name").Description,
	)
	injectLabelFlags(&createDoc.Command, false)

	// Get a registered resource
	getDoc := man.Docs.GetCommand("policy/registered-resources/get",
		man.WithRun(policy_getRegisteredResource),
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Shorthand,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)

	// List registered resources
	listDoc := man.Docs.GetCommand("policy/registered-resources/list",
		man.WithRun(policy_listRegisteredResources),
	)
	injectListPaginationFlags(listDoc)

	// Update a registered resource
	updateDoc := man.Docs.GetCommand("policy/registered-resources/update",
		man.WithRun(policy_updateRegisteredResource),
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	injectLabelFlags(&updateDoc.Command, true)

	// Delete a registered resource
	deleteDoc := man.Docs.GetCommand("policy/registered-resources/delete",
		man.WithRun(policy_deleteRegisteredResource),
	)
	deleteDoc.Flags().StringP(
		deleteDoc.GetDocFlag("id").Name,
		deleteDoc.GetDocFlag("id").Shorthand,
		deleteDoc.GetDocFlag("id").Default,
		deleteDoc.GetDocFlag("id").Description,
	)

	policy_registeredResourcesCmd.AddSubcommands(createDoc, getDoc, listDoc, updateDoc, deleteDoc)
	policyCmd.AddCommand(&policy_registeredResourcesCmd.Command)
}
