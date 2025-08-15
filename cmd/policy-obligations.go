package cmd

import (
	"fmt"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var (
	obligationValues []string
)

//
// Obligations
//

func policyCreateObligation(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	name := c.Flags.GetRequiredString("name")
	obligationValues = c.Flags.GetStringSlice("value", obligationValues, cli.FlagsStringSliceOptions{})
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	obl, err := h.CreateObligation(cmd.Context(), name, obligationValues, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create obligation", err)
	}

	simpleObligationValues := cli.GetSimpleObligationValues(obl.GetValues())

	rows := [][]string{
		{"Id", obl.GetId()},
		{"Name", obl.GetName()},
		{"Values", cli.CommaSeparated(simpleObligationValues)},
	}

	if mdRows := getMetadataRows(obl.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, obl.GetId(), t, obl)
}

func policyGetObligation(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetOptionalID("id")
	name := c.Flags.GetOptionalString("name")

	if id == "" && name == "" {
		cli.ExitWithError("Either 'id' or 'name' must be provided", nil)
	}

	resource, err := h.GetRegisteredResource(cmd.Context(), id, name)
	if err != nil {
		identifier := fmt.Sprintf("id: %s", id)
		if id == "" {
			identifier = fmt.Sprintf("name: %s", name)
		}
		errMsg := fmt.Sprintf("Failed to find registered resource (%s)", identifier)
		cli.ExitWithError(errMsg, err)
	}

	simpleRegResValues := cli.GetSimpleRegisteredResourceValues(resource.GetValues())

	rows := [][]string{
		{"Id", resource.GetId()},
		{"Name", resource.GetName()},
		{"Values", cli.CommaSeparated(simpleRegResValues)},
	}
	if mdRows := getMetadataRows(resource.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, resource.GetId(), t, resource)
}

func policyListObligations(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")

	resources, page, err := h.ListRegisteredResources(cmd.Context(), limit, offset)
	if err != nil {
		cli.ExitWithError("Failed to list registered resources", err)
	}

	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("name", "Name", cli.FlexColumnWidthFour),
		table.NewFlexColumn("values", "Values", cli.FlexColumnWidthTwo),
		// todo: do we need to show metadata labels and created/updated at?
	)
	rows := []table.Row{}
	for _, r := range resources {
		simpleRegResValues := cli.GetSimpleRegisteredResourceValues(r.GetValues())
		rows = append(rows, table.NewRow(table.RowData{
			"id":     r.GetId(),
			"name":   r.GetName(),
			"values": cli.CommaSeparated(simpleRegResValues),
			// todo: do we need to show metadata labels and created/updated at?
		}))
	}
	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, resources)
}

func policyUpdateObligation(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	name := c.Flags.GetOptionalString("name")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	updated, err := h.UpdateRegisteredResource(
		cmd.Context(),
		id,
		name,
		getMetadataMutable(metadataLabels),
		getMetadataUpdateBehavior(),
	)
	if err != nil {
		cli.ExitWithError("Failed to update registered resource", err)
	}

	rows := [][]string{
		{"Id", id},
		{"Name", updated.GetName()},
	}
	if mdRows := getMetadataRows(updated.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, id, t, updated)
}

func policyDeleteObligation(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	force := c.Flags.GetRequiredBool("force")
	ctx := cmd.Context()

	resource, err := h.GetRegisteredResource(ctx, id, "")
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find registered resource (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	cli.ConfirmAction(cli.ActionDelete, "registered resource", id, force)

	err = h.DeleteRegisteredResource(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to delete registered resource (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	rows := [][]string{
		{"Id", id},
		{"Name", resource.GetName()},
	}
	if mdRows := getMetadataRows(resource.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, id, t, resource)
}

func init() {
	// Obligations commands

	getDoc := man.Docs.GetCommand("policy/obligations/get",
		man.WithRun(policyGetObligation),
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Shorthand,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("name").Name,
		getDoc.GetDocFlag("name").Shorthand,
		getDoc.GetDocFlag("name").Default,
		getDoc.GetDocFlag("name").Description,
	)

	listDoc := man.Docs.GetCommand("policy/obligations/list",
		man.WithRun(policyListObligations),
	)
	injectListPaginationFlags(listDoc)

	createDoc := man.Docs.GetCommand("policy/obligations/create",
		man.WithRun(policyCreateObligation),
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("name").Name,
		createDoc.GetDocFlag("name").Shorthand,
		createDoc.GetDocFlag("name").Default,
		createDoc.GetDocFlag("name").Description,
	)
	createDoc.Flags().StringSliceVarP(
		&registeredResourceValues,
		createDoc.GetDocFlag("value").Name,
		createDoc.GetDocFlag("value").Shorthand,
		[]string{},
		createDoc.GetDocFlag("value").Description,
	)
	injectLabelFlags(&createDoc.Command, false)

	updateDoc := man.Docs.GetCommand("policy/obligations/update",
		man.WithRun(policyUpdateObligation),
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("name").Name,
		updateDoc.GetDocFlag("name").Shorthand,
		updateDoc.GetDocFlag("name").Default,
		updateDoc.GetDocFlag("name").Description,
	)
	injectLabelFlags(&updateDoc.Command, true)

	deleteDoc := man.Docs.GetCommand("policy/obligations/delete",
		man.WithRun(policyDeleteObligation),
	)
	deleteDoc.Flags().StringP(
		deleteDoc.GetDocFlag("id").Name,
		deleteDoc.GetDocFlag("id").Shorthand,
		deleteDoc.GetDocFlag("id").Default,
		deleteDoc.GetDocFlag("id").Description,
	)
	deleteDoc.Flags().Bool(
		deleteDoc.GetDocFlag("force").Name,
		false,
		deleteDoc.GetDocFlag("force").Description,
	)

	// Add commands to the policy command

	policyRegisteredResourcesDoc := man.Docs.GetCommand("policy/registered-resources",
		man.WithSubcommands(
			getDoc,
			listDoc,
			createDoc,
			updateDoc,
			deleteDoc,
		),
	)

	policyCmd.AddCommand(&policyRegisteredResourcesDoc.Command)
}
