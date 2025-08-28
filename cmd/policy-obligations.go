package cmd

import (
	"fmt"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

//
// Obligations
//

var obligationValues []string

func policyCreateObligation(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()
	name := c.Flags.GetRequiredString("name")
	obligationValues = c.Flags.GetStringSlice("value", obligationValues, cli.FlagsStringSliceOptions{})
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})
	namespace := c.Flags.GetRequiredString("namespace")
	obl, err := h.CreateObligation(cmd.Context(), namespace, name, obligationValues, getMetadataMutable(metadataLabels))
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
	fqn := c.Flags.GetOptionalString("fqn")

	if id == "" && fqn == "" {
		cli.ExitWithError("Either 'id' or 'fqn' must be provided", nil)
	}

	obl, err := h.GetObligation(cmd.Context(), id, fqn)
	if err != nil {
		identifier := fmt.Sprintf("id: %s", id)
		if id == "" {
			identifier = fmt.Sprintf("fqn: %s", fqn)
		}
		errMsg := fmt.Sprintf("Failed to find obligation (%s)", identifier)
		cli.ExitWithError(errMsg, err)
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

func policyListObligations(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	namespace := c.Flags.GetOptionalString("namespace")
	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")

	obls, page, err := h.ListObligations(cmd.Context(), limit, offset, namespace)
	if err != nil {
		cli.ExitWithError("Failed to list obligations", err)
	}

	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("name", "Name", cli.FlexColumnWidthFour),
		table.NewFlexColumn("values", "Values", cli.FlexColumnWidthTwo),
	)
	rows := []table.Row{}
	for _, r := range obls {
		simpleObligationValues := cli.GetSimpleObligationValues(r.GetValues())
		rows = append(rows, table.NewRow(table.RowData{
			"id":     r.GetId(),
			"name":   r.GetName(),
			"values": cli.CommaSeparated(simpleObligationValues),
		}))
	}
	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, obls)
}

func policyUpdateObligation(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	name := c.Flags.GetOptionalString("name")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	updated, err := h.UpdateObligation(
		cmd.Context(),
		id,
		name,
		getMetadataMutable(metadataLabels),
		getMetadataUpdateBehavior(),
	)
	if err != nil {
		cli.ExitWithError("Failed to update obligation", err)
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

	id := c.Flags.GetOptionalID("id")
	fqn := c.Flags.GetOptionalString("fqn")

	if id == "" && fqn == "" {
		cli.ExitWithError("Either 'id' or 'fqn' must be provided", nil)
	}

	force := c.Flags.GetRequiredBool("force")
	ctx := cmd.Context()

	obl, err := h.GetObligation(ctx, id, fqn)
	identifier := id
	if id == "" {
		identifier = fqn
	}
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find obligation (%s)", identifier)
		cli.ExitWithError(errMsg, err)
	}
	id = obl.GetId()
	cli.ConfirmAction(cli.ActionDelete, "obligation", identifier, force)

	err = h.DeleteObligation(ctx, id, fqn)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to delete obligation (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	rows := [][]string{
		{"Id", id},
		{"Name", obl.GetName()},
	}
	if mdRows := getMetadataRows(obl.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, id, t, obl)
}

//
// Obligation Values
//

func policyCreateObligationValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	ctx := cmd.Context()
	obligation := c.Flags.GetRequiredString("obligation")
	value := c.Flags.GetRequiredString("value")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	oblVal, err := h.CreateObligationValue(ctx, obligation, value, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create obligation value", err)
	}

	rows := [][]string{
		{"Id", oblVal.GetId()},
		{"Name", oblVal.GetObligation().GetName()},
		{"Value", oblVal.GetValue()},
	}
	if mdRows := getMetadataRows(oblVal.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, oblVal.GetId(), t, oblVal)
}

func policyGetObligationValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetOptionalID("id")
	fqn := c.Flags.GetOptionalString("fqn")

	if id == "" && fqn == "" {
		cli.ExitWithError("Either 'id' or 'fqn' must be provided", nil)
	}

	value, err := h.GetObligationValue(cmd.Context(), id, fqn)
	if err != nil {
		identifier := fmt.Sprintf("id: %s", id)
		if id == "" {
			identifier = fmt.Sprintf("fqn: %s", fqn)
		}
		errMsg := fmt.Sprintf("Failed to find obligation value (%s)", identifier)
		cli.ExitWithError(errMsg, err)
	}

	rows := [][]string{
		{"Id", value.GetId()},
		{"Name", value.GetObligation().GetName()},
		{"Value", value.GetValue()},
	}
	if mdRows := getMetadataRows(value.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, value.GetId(), t, value)
}

func policyUpdateObligationValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	value := c.Flags.GetOptionalString("value")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	updated, err := h.UpdateObligationValue(
		cmd.Context(),
		id,
		value,
		getMetadataMutable(metadataLabels),
		getMetadataUpdateBehavior(),
	)
	if err != nil {
		cli.ExitWithError("Failed to update obligation value", err)
	}

	rows := [][]string{
		{"Id", id},
		{"Name", updated.GetObligation().GetName()},
		{"Value", updated.GetValue()},
	}
	if mdRows := getMetadataRows(updated.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, id, t, updated)
}

func policyDeleteObligationValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetOptionalID("id")
	fqn := c.Flags.GetOptionalString("fqn")

	if id == "" && fqn == "" {
		cli.ExitWithError("Either 'id' or 'fqn' must be provided", nil)
	}

	force := c.Flags.GetOptionalBool("force")
	ctx := cmd.Context()

	val, err := h.GetObligationValue(ctx, id, fqn)
	identifier := id
	if id == "" {
		identifier = fqn
	}
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find obligation value (%s)", identifier)
		cli.ExitWithError(errMsg, err)
	}
	id = val.GetId()
	cli.ConfirmAction(cli.ActionDelete, "obligation value", identifier, force)

	err = h.DeleteObligationValue(ctx, id, fqn)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to delete obligation value (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	rows := [][]string{
		{"Id", id},
		{"Value", val.GetValue()},
	}
	if mdRows := getMetadataRows(val.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, id, t, val)
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
		getDoc.GetDocFlag("fqn").Name,
		getDoc.GetDocFlag("fqn").Shorthand,
		getDoc.GetDocFlag("fqn").Default,
		getDoc.GetDocFlag("fqn").Description,
	)

	listDoc := man.Docs.GetCommand("policy/obligations/list",
		man.WithRun(policyListObligations),
	)
	listDoc.Flags().StringP(
		listDoc.GetDocFlag("namespace").Name,
		listDoc.GetDocFlag("namespace").Shorthand,
		listDoc.GetDocFlag("namespace").Default,
		listDoc.GetDocFlag("namespace").Description,
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
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("namespace").Name,
		createDoc.GetDocFlag("namespace").Shorthand,
		createDoc.GetDocFlag("namespace").Default,
		createDoc.GetDocFlag("namespace").Description,
	)
	createDoc.Flags().StringSliceVarP(
		&obligationValues,
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
	deleteDoc.Flags().StringP(
		deleteDoc.GetDocFlag("fqn").Name,
		deleteDoc.GetDocFlag("fqn").Shorthand,
		deleteDoc.GetDocFlag("fqn").Default,
		deleteDoc.GetDocFlag("fqn").Description,
	)
	deleteDoc.Flags().Bool(
		deleteDoc.GetDocFlag("force").Name,
		false,
		deleteDoc.GetDocFlag("force").Description,
	)

	// Obligation Values commands

	getValueDoc := man.Docs.GetCommand("policy/obligations/values/get",
		man.WithRun(policyGetObligationValue),
	)
	getValueDoc.Flags().StringP(
		getValueDoc.GetDocFlag("id").Name,
		getValueDoc.GetDocFlag("id").Shorthand,
		getValueDoc.GetDocFlag("id").Default,
		getValueDoc.GetDocFlag("id").Description,
	)
	getValueDoc.Flags().StringP(
		getValueDoc.GetDocFlag("fqn").Name,
		getValueDoc.GetDocFlag("fqn").Shorthand,
		getValueDoc.GetDocFlag("fqn").Default,
		getValueDoc.GetDocFlag("fqn").Description,
	)

	createValueDoc := man.Docs.GetCommand("policy/obligations/values/create",
		man.WithRun(policyCreateObligationValue),
	)
	createValueDoc.Flags().StringP(
		createValueDoc.GetDocFlag("obligation").Name,
		createValueDoc.GetDocFlag("obligation").Shorthand,
		createValueDoc.GetDocFlag("obligation").Default,
		createValueDoc.GetDocFlag("obligation").Description,
	)
	createValueDoc.Flags().StringP(
		createValueDoc.GetDocFlag("value").Name,
		createValueDoc.GetDocFlag("value").Shorthand,
		createValueDoc.GetDocFlag("value").Default,
		createValueDoc.GetDocFlag("value").Description,
	)
	injectLabelFlags(&createValueDoc.Command, false)

	updateValueDoc := man.Docs.GetCommand("policy/obligations/values/update",
		man.WithRun(policyUpdateObligationValue),
	)
	updateValueDoc.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	updateValueDoc.Flags().StringP(
		updateValueDoc.GetDocFlag("value").Name,
		updateValueDoc.GetDocFlag("value").Shorthand,
		updateValueDoc.GetDocFlag("value").Default,
		updateValueDoc.GetDocFlag("value").Description,
	)
	updateValueDoc.Flags().StringSliceVarP(
		&actionAttributeValues,
		updateValueDoc.GetDocFlag("action-attribute-value").Name,
		updateValueDoc.GetDocFlag("action-attribute-value").Shorthand,
		[]string{},
		updateValueDoc.GetDocFlag("action-attribute-value").Description,
	)
	injectLabelFlags(&updateValueDoc.Command, true)
	updateValueDoc.Flags().Bool(
		updateValueDoc.GetDocFlag("force").Name,
		false,
		updateValueDoc.GetDocFlag("force").Description,
	)

	deleteValueDoc := man.Docs.GetCommand("policy/registered-resources/values/delete",
		man.WithRun(policyDeleteRegisteredResourceValue),
	)
	deleteValueDoc.Flags().StringP(
		deleteValueDoc.GetDocFlag("id").Name,
		deleteValueDoc.GetDocFlag("id").Shorthand,
		deleteValueDoc.GetDocFlag("id").Default,
		deleteValueDoc.GetDocFlag("id").Description,
	)
	deleteValueDoc.Flags().Bool(
		deleteValueDoc.GetDocFlag("force").Name,
		false,
		deleteValueDoc.GetDocFlag("force").Description,
	)

	// Add commands to the policy command

	policyObligationsDoc := man.Docs.GetCommand("policy/obligations",
		man.WithSubcommands(
			getDoc,
			listDoc,
			createDoc,
			updateDoc,
			deleteDoc,
		),
	)

	policyCmd.AddCommand(&policyObligationsDoc.Command)
}
