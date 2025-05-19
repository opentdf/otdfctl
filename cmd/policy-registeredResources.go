package cmd

import (
	"fmt"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
)

var (
	values []string
)

//
// Registered Resources
//

func policyCreateRegisteredResource(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	name := c.Flags.GetRequiredString("name")
	values = c.Flags.GetStringSlice("value", values, cli.FlagsStringSliceOptions{})
	metadataLabels := c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	resource, err := h.CreateRegisteredResource(cmd.Context(), name, values, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create registered resource", err)
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

func policyGetRegisteredResource(cmd *cobra.Command, args []string) {
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

func policyListRegisteredResources(cmd *cobra.Command, args []string) {
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
	t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, resources)
}

func policyUpdateRegisteredResource(cmd *cobra.Command, args []string) {
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

func policyDeleteRegisteredResource(cmd *cobra.Command, args []string) {
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

//
// Registered Resource Values
//

func policyCreateRegisteredResourceValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	resourceId := c.Flags.GetRequiredID("resource-id")
	value := c.Flags.GetRequiredString("value")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	resourceValue, err := h.CreateRegisteredResourceValue(cmd.Context(), resourceId, value, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create registered resource value", err)
	}

	rows := [][]string{
		{"Id", resourceValue.GetId()},
		{"Value", resourceValue.GetValue()},
	}
	if mdRows := getMetadataRows(resourceValue.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, resourceValue.GetId(), t, resourceValue)
}

func policyGetRegisteredResourceValue(cmd *cobra.Command, args []string) {
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

	rows := [][]string{
		{"Id", value.GetId()},
		{"Value", value.GetValue()},
	}
	if mdRows := getMetadataRows(value.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, value.GetId(), t, value)
}

func policyListRegisteredResourceValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	resourceId := c.Flags.GetOptionalString("resource-id")
	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")

	values, page, err := h.ListRegisteredResourceValues(cmd.Context(), resourceId, limit, offset)
	if err != nil {
		cli.ExitWithError("Failed to list registered resource values", err)
	}

	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("value", "Value", cli.FlexColumnWidthFour),
	)
	rows := []table.Row{}
	for _, v := range values {
		rows = append(rows, table.NewRow(table.RowData{
			"id":    v.GetId(),
			"value": v.GetValue(),
		}))
	}
	list := append([]*policy.RegisteredResourceValue{}, values...)

	t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, list)
}

func policyUpdateRegisteredResourceValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	value := c.Flags.GetOptionalString("value")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	updated, err := h.UpdateRegisteredResourceValue(
		cmd.Context(),
		id,
		value,
		getMetadataMutable(metadataLabels),
		getMetadataUpdateBehavior(),
	)
	if err != nil {
		cli.ExitWithError("Failed to update registered resource value", err)
	}

	rows := [][]string{
		{"Id", id},
		{"Value", updated.GetValue()},
	}
	if mdRows := getMetadataRows(updated.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, id, t, updated)
}

func policyDeleteRegisteredResourceValue(cmd *cobra.Command, args []string) {
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
	// Registered Resources
	// todo

	// Registered Resource Values
	// todo

	policyRegisteredResourcesDoc := man.Docs.GetCommand("policy/registered-resources", man.WithSubcommands(nil)) // todo

	policyRegisteredResourceValuesDoc := man.Docs.GetCommand("policy/registered-resources/values", man.WithSubcommands(nil)) // todo

	policyCmd.AddCommand(&policyRegisteredResourcesDoc.Command)
	policyCmd.AddCommand(&policyRegisteredResourceValuesDoc.Command)
}
