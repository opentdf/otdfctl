package cmd

import (
	"fmt"

	"github.com/evertras/bubble-table/table"
	"github.com/google/uuid"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
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
	resource := c.Flags.GetRequiredString("resource")
	value := c.Flags.GetRequiredString("value")
	actionAttributeValues = c.Flags.GetStringSlice("action-attribute-value", actionAttributeValues, cli.FlagsStringSliceOptions{Min: 0})
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	var resourceID string
	if uuid.Validate(resource) == nil {
		resourceID = resource
	} else {
		resourceByName, err := h.GetRegisteredResource(ctx, "", resource)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to find registered resource (name: %s)", resource), err)
		}
		resourceID = resourceByName.GetId()
	}

	parsedActionAttributeValues := parseActionAttributeValueArgs(actionAttributeValues)

	resourceValue, err := h.CreateRegisteredResourceValue(ctx, resourceID, value, parsedActionAttributeValues, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create registered resource value", err)
	}

	simpleActionAttributeValues := cli.GetSimpleRegisteredResourceActionAttributeValues(resourceValue.GetActionAttributeValues())

	rows := [][]string{
		{"Id", resourceValue.GetId()},
		{"Value", resourceValue.GetValue()},
		{"Action Attribute Values", cli.CommaSeparated(simpleActionAttributeValues)},
	}
	if mdRows := getMetadataRows(resourceValue.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, resourceValue.GetId(), t, resourceValue)
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

	value, err := h.GetRegisteredResourceValue(cmd.Context(), id, fqn)
	if err != nil {
		identifier := fmt.Sprintf("id: %s", id)
		if id == "" {
			identifier = fmt.Sprintf("fqn: %s", fqn)
		}
		errMsg := fmt.Sprintf("Failed to find registered resource value (%s)", identifier)
		cli.ExitWithError(errMsg, err)
	}

	simpleActionAttributeValues := cli.GetSimpleRegisteredResourceActionAttributeValues(value.GetActionAttributeValues())

	rows := [][]string{
		{"Id", value.GetId()},
		{"Value", value.GetValue()},
		{"Action Attribute Values", cli.CommaSeparated(simpleActionAttributeValues)},
	}
	if mdRows := getMetadataRows(value.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, value.GetId(), t, value)
}

func policyListObligationValues(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	ctx := cmd.Context()
	resource := c.Flags.GetRequiredString("resource")
	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")

	var resourceID string
	if uuid.Validate(resource) == nil {
		resourceID = resource
	} else {
		resourceByName, err := h.GetRegisteredResource(ctx, "", resource)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to find registered resource (name: %s)", resource), err)
		}
		resourceID = resourceByName.GetId()
	}

	values, page, err := h.ListRegisteredResourceValues(ctx, resourceID, limit, offset)
	if err != nil {
		cli.ExitWithError("Failed to list registered resource values", err)
	}

	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("value", "Value", cli.FlexColumnWidthFour),
		table.NewFlexColumn("action-attribute-values", "Action Attribute Values", cli.FlexColumnWidthFour),
	)
	rows := []table.Row{}
	for _, v := range values {
		simpleActionAttributeValues := cli.GetSimpleRegisteredResourceActionAttributeValues(v.GetActionAttributeValues())

		rows = append(rows, table.NewRow(table.RowData{
			"id":                      v.GetId(),
			"value":                   v.GetValue(),
			"action-attribute-values": cli.CommaSeparated(simpleActionAttributeValues),
		}))
	}
	list := append([]*policy.RegisteredResourceValue{}, values...)

	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, list)
}

func policyUpdateObligationValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	value := c.Flags.GetOptionalString("value")
	actionAttributeValues = c.Flags.GetStringSlice("action-attribute-value", actionAttributeValues, cli.FlagsStringSliceOptions{Min: 0})
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})
	force := c.Flags.GetOptionalBool("force")

	parsedActionAttributeValues := parseActionAttributeValueArgs(actionAttributeValues)

	// only confirm if new action attribute values provided
	if len(parsedActionAttributeValues) > 0 {
		cli.ConfirmActionSubtext(cli.ActionUpdate, "registered resource value", id,
			"All existing action attribute values will be replaced with the new ones provided.",
			force)
	}

	updated, err := h.UpdateRegisteredResourceValue(
		cmd.Context(),
		id,
		value,
		parsedActionAttributeValues,
		getMetadataMutable(metadataLabels),
		getMetadataUpdateBehavior(),
	)
	if err != nil {
		cli.ExitWithError("Failed to update registered resource value", err)
	}

	simpleActionAttributeValues := cli.GetSimpleRegisteredResourceActionAttributeValues(updated.GetActionAttributeValues())

	rows := [][]string{
		{"Id", id},
		{"Value", updated.GetValue()},
		{"Action Attribute Values", cli.CommaSeparated(simpleActionAttributeValues)},
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

	id := c.Flags.GetRequiredID("id")
	force := c.Flags.GetOptionalBool("force")
	ctx := cmd.Context()

	resource, err := h.GetRegisteredResourceValue(ctx, id, "")
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find registered resource value (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	cli.ConfirmAction(cli.ActionDelete, "registered resource value", id, force)

	err = h.DeleteRegisteredResourceValue(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to delete registered resource value (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	rows := [][]string{
		{"Id", id},
		{"Value", resource.GetValue()},
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
