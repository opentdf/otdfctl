package cmd

import (
	"fmt"
	"strings"

	"github.com/evertras/bubble-table/table"
	"github.com/google/uuid"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/registeredresources"
	"github.com/spf13/cobra"
)

var (
	registeredResourceValues []string
	actionAttributeValues    []string
)

const actionAttributeValueFlagSplitCount = 2

//
// Registered Resources
//

func policyCreateRegisteredResource(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	name := c.Flags.GetRequiredString("name")
	registeredResourceValues = c.Flags.GetStringSlice("value", registeredResourceValues, cli.FlagsStringSliceOptions{})
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	resource, err := h.CreateRegisteredResource(cmd.Context(), name, registeredResourceValues, getMetadataMutable(metadataLabels))
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

	ctx := cmd.Context()
	resourceID := c.Flags.GetRequiredID("resource-id")
	value := c.Flags.GetRequiredString("value")
	actionAttributeValues = c.Flags.GetStringSlice("action-attribute-value", actionAttributeValues, cli.FlagsStringSliceOptions{Min: 0})
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	parsedActionAttributeValues := make([]*registeredresources.ActionAttributeValue, len(actionAttributeValues))

	for i, aav := range actionAttributeValues {
		// split on semicolon
		split := strings.Split(aav, ";")
		if len(split) != actionAttributeValueFlagSplitCount {
			cli.ExitWithError("Invalid action attribute value format", nil)
		}

		actionIdentifier := split[0]
		attrValIdentifier := split[1]

		newActionAttrVal := &registeredresources.ActionAttributeValue{}

		if uuid.Validate(actionIdentifier) == nil {
			newActionAttrVal.ActionIdentifier = &registeredresources.ActionAttributeValue_ActionId{
				ActionId: actionIdentifier,
			}
		} else {
			newActionAttrVal.ActionIdentifier = &registeredresources.ActionAttributeValue_ActionName{
				ActionName: actionIdentifier,
			}
		}

		if uuid.Validate(attrValIdentifier) == nil {
			newActionAttrVal.AttributeValueIdentifier = &registeredresources.ActionAttributeValue_AttributeValueId{
				AttributeValueId: attrValIdentifier,
			}
		} else {
			newActionAttrVal.AttributeValueIdentifier = &registeredresources.ActionAttributeValue_AttributeValueFqn{
				AttributeValueFqn: attrValIdentifier,
			}
		}

		parsedActionAttributeValues[i] = newActionAttrVal
	}

	resourceValue, err := h.CreateRegisteredResourceValue(ctx, resourceID, value, parsedActionAttributeValues, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create registered resource value", err)
	}

	simpleActionAttributeValues := cli.GetSimpleRegisteredResourceActionAttributeValues(resourceValue.GetActionAttributeValues())

	rows := [][]string{
		{"Id", resourceValue.GetId()},
		{"Value", resourceValue.GetValue()},
		// todo: figure out how to stop this output from being cut off
		{"Action Attribute Values", cli.CommaSeparated(simpleActionAttributeValues)},
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

func policyListRegisteredResourceValues(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	resourceID := c.Flags.GetOptionalString("resource-id")
	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")

	values, page, err := h.ListRegisteredResourceValues(cmd.Context(), resourceID, limit, offset)
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
	// todo: figure out how to get action attribute values from command line
	actionAttributeValues := []*registeredresources.ActionAttributeValue{}
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	updated, err := h.UpdateRegisteredResourceValue(
		cmd.Context(),
		id,
		value,
		actionAttributeValues,
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
	// Registered Resources commands

	getDoc := man.Docs.GetCommand("policy/registered-resources/get",
		man.WithRun(policyGetRegisteredResource),
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

	listDoc := man.Docs.GetCommand("policy/registered-resources/list",
		man.WithRun(policyListRegisteredResources),
	)
	injectListPaginationFlags(listDoc)

	createDoc := man.Docs.GetCommand("policy/registered-resources/create",
		man.WithRun(policyCreateRegisteredResource),
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

	updateDoc := man.Docs.GetCommand("policy/registered-resources/update",
		man.WithRun(policyUpdateRegisteredResource),
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

	deleteDoc := man.Docs.GetCommand("policy/registered-resources/delete",
		man.WithRun(policyDeleteRegisteredResource),
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

	// Registered Resource Values commands

	getValueDoc := man.Docs.GetCommand("policy/registered-resources/values/get",
		man.WithRun(policyGetRegisteredResourceValue),
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

	listValuesDoc := man.Docs.GetCommand("policy/registered-resources/values/list",
		man.WithRun(policyListRegisteredResourceValues),
	)
	listValuesDoc.Flags().StringP(
		listValuesDoc.GetDocFlag("resource-id").Name,
		listValuesDoc.GetDocFlag("resource-id").Shorthand,
		listValuesDoc.GetDocFlag("resource-id").Default,
		listValuesDoc.GetDocFlag("resource-id").Description,
	)
	injectListPaginationFlags(listValuesDoc)

	createValueDoc := man.Docs.GetCommand("policy/registered-resources/values/create",
		man.WithRun(policyCreateRegisteredResourceValue),
	)
	createValueDoc.Flags().StringP(
		createValueDoc.GetDocFlag("resource-id").Name,
		createValueDoc.GetDocFlag("resource-id").Shorthand,
		createValueDoc.GetDocFlag("resource-id").Default,
		createValueDoc.GetDocFlag("resource-id").Description,
	)
	createValueDoc.Flags().StringP(
		createValueDoc.GetDocFlag("value").Name,
		createValueDoc.GetDocFlag("value").Shorthand,
		createValueDoc.GetDocFlag("value").Default,
		createValueDoc.GetDocFlag("value").Description,
	)
	createValueDoc.Flags().StringSliceVarP(
		&actionAttributeValues,
		createValueDoc.GetDocFlag("action-attribute-value").Name,
		createValueDoc.GetDocFlag("action-attribute-value").Shorthand,
		[]string{},
		createValueDoc.GetDocFlag("action-attribute-value").Description,
	)
	injectLabelFlags(&createValueDoc.Command, false)

	updateValueDoc := man.Docs.GetCommand("policy/registered-resources/values/update",
		man.WithRun(policyUpdateRegisteredResourceValue),
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
	injectLabelFlags(&updateValueDoc.Command, true)

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

	policyRegisteredResourcesDoc := man.Docs.GetCommand("policy/registered-resources",
		man.WithSubcommands(
			getDoc,
			listDoc,
			createDoc,
			updateDoc,
			deleteDoc,
		),
	)

	policyRegisteredResourceValuesDoc := man.Docs.GetCommand("policy/registered-resources/values",
		man.WithSubcommands(
			getValueDoc,
			listValuesDoc,
			createValueDoc,
			updateValueDoc,
			deleteValueDoc,
		),
	)

	policyRegisteredResourcesDoc.AddCommand(&policyRegisteredResourceValuesDoc.Command)
	policyCmd.AddCommand(&policyRegisteredResourcesDoc.Command)
}
