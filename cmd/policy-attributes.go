package cmd

import (
	"fmt"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
)

var (
	forceReplaceMetadataLabels bool
	attributeValues            []string
	attributeValuesOrder       []string

	policy_attributesCmd = man.Docs.GetCommand("policy/attributes")
)

func policy_createAttribute(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	name := c.Flags.GetRequiredString("name")
	rule := c.Flags.GetRequiredString("rule")
	attributeValues = c.Flags.GetStringSlice("value", attributeValues, cli.FlagsStringSliceOptions{})
	namespace := c.Flags.GetRequiredString("namespace")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	attr, err := h.CreateAttribute(cmd.Context(), name, rule, namespace, attributeValues, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create attribute", err)
	}

	a := cli.GetSimpleAttribute(&policy.Attribute{
		Id:        attr.GetId(),
		Name:      attr.GetName(),
		Rule:      attr.GetRule(),
		Values:    attr.GetValues(),
		Namespace: attr.GetNamespace(),
	})
	rows := [][]string{
		{"Name", a.Name},
		{"Rule", a.Rule},
		{"Values", cli.CommaSeparated(a.Values)},
		{"Namespace", a.Namespace},
	}
	if mdRows := getMetadataRows(attr.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)

	HandleSuccess(cmd, a.Id, t, attr)
}

func policy_getAttribute(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	attr, err := h.GetAttribute(cmd.Context(), id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get attribute (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	a := cli.GetSimpleAttribute(attr)
	rows := [][]string{
		{"Id", a.Id},
		{"Name", a.Name},
		{"Rule", a.Rule},
		{"Values", cli.CommaSeparated(a.Values)},
		{"Namespace", a.Namespace},
	}
	if mdRows := getMetadataRows(attr.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, a.Id, t, attr)
}

func policy_listAttributes(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	state := cli.GetState(cmd)
	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")

	attrs, page, err := h.ListAttributes(cmd.Context(), state, limit, offset)
	if err != nil {
		cli.ExitWithError("Failed to list attributes", err)
	}

	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("namespace", "Namespace", cli.FlexColumnWidthFour),
		table.NewFlexColumn("name", "Name", cli.FlexColumnWidthThree),
		table.NewFlexColumn("rule", "Rule", cli.FlexColumnWidthTwo),
		table.NewFlexColumn("values", "Values", cli.FlexColumnWidthTwo),
		table.NewFlexColumn("active", "Active", cli.FlexColumnWidthTwo),
		table.NewFlexColumn("labels", "Labels", cli.FlexColumnWidthOne),
		table.NewFlexColumn("created_at", "Created At", cli.FlexColumnWidthOne),
		table.NewFlexColumn("updated_at", "Updated At", cli.FlexColumnWidthOne),
	)
	rows := []table.Row{}
	for _, attr := range attrs {
		a := cli.GetSimpleAttribute(attr)
		rows = append(rows, table.NewRow(table.RowData{
			"id":         a.Id,
			"namespace":  a.Namespace,
			"name":       a.Name,
			"rule":       a.Rule,
			"values":     cli.CommaSeparated(a.Values),
			"active":     a.Active,
			"labels":     a.Metadata["Labels"],
			"created_at": a.Metadata["Created At"],
			"updated_at": a.Metadata["Updated At"],
		}))
	}
	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, attrs)
}

func policy_deactivateAttribute(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	ctx := cmd.Context()
	id := c.Flags.GetRequiredID("id")
	force := c.Flags.GetOptionalBool("force")

	attr, err := h.GetAttribute(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get attribute (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	cli.ConfirmAction(cli.ActionDeactivate, "attribute", attr.GetName(), force)

	attr, err = h.DeactivateAttribute(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to deactivate attribute (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	a := cli.GetSimpleAttribute(attr)
	rows := [][]string{
		{"Name", a.Name},
		{"Rule", a.Rule},
		{"Values", cli.CommaSeparated(a.Values)},
		{"Namespace", a.Namespace},
	}
	if mdRows := getMetadataRows(attr.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, a.Id, t, a)
}

func policy_updateAttribute(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	if a, err := h.UpdateAttribute(cmd.Context(), id, getMetadataMutable(metadataLabels), getMetadataUpdateBehavior()); err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to update attribute (%s)", id), err)
	} else {
		rows := [][]string{
			{"Id", a.GetId()},
			{"Name", a.GetName()},
		}
		if mdRows := getMetadataRows(a.GetMetadata()); mdRows != nil {
			rows = append(rows, mdRows...)
		}
		t := cli.NewTabular(rows...)
		HandleSuccess(cmd, id, t, a)
	}
}

func policy_unsafeReactivateAttribute(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	ctx := cmd.Context()
	id := c.Flags.GetRequiredID("id")

	a, err := h.GetAttribute(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get attribute (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	if !forceUnsafe {
		cli.ConfirmTextInput(cli.ActionReactivate, "attribute", cli.InputNameFQN, a.GetFqn())
	}

	if reactivatedAttr, err := h.UnsafeReactivateAttribute(ctx, id); err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to reactivate attribute (%s)", id), err)
	} else {
		rows := [][]string{
			{"Id", reactivatedAttr.GetId()},
			{"Name", reactivatedAttr.GetName()},
		}
		if mdRows := getMetadataRows(a.GetMetadata()); mdRows != nil {
			rows = append(rows, mdRows...)
		}
		t := cli.NewTabular(rows...)
		HandleSuccess(cmd, id, t, a)
	}
}

func policy_unsafeUpdateAttribute(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	ctx := cmd.Context()
	id := c.Flags.GetRequiredID("id")
	name := c.Flags.GetOptionalString("name")
	rule := c.Flags.GetOptionalString("rule")
	attributeValuesOrder = c.Flags.GetStringSlice("values-order", attributeValuesOrder, cli.FlagsStringSliceOptions{})

	a, err := h.GetAttribute(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get attribute (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	if !forceUnsafe {
		cli.ConfirmTextInput(cli.ActionUpdateUnsafe, "attribute", cli.InputNameFQN, a.GetFqn())
	}

	if err := h.UnsafeUpdateAttribute(ctx, id, name, rule, attributeValuesOrder); err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to update attribute (%s)", id), err)
	} else {
		var (
			retrievedVals []string
			valueIDs      []string
		)
		for _, v := range a.GetValues() {
			retrievedVals = append(retrievedVals, v.GetValue())
			valueIDs = append(valueIDs, v.GetId())
		}
		rows := [][]string{
			{"Id", a.GetId()},
			{"Name", a.GetName()},
			{"Rule", handlers.GetAttributeRuleFromAttributeType(a.GetRule())},
			{"Values", cli.CommaSeparated(retrievedVals)},
			{"Value IDs", cli.CommaSeparated(valueIDs)},
		}
		if mdRows := getMetadataRows(a.GetMetadata()); mdRows != nil {
			rows = append(rows, mdRows...)
		}
		t := cli.NewTabular(rows...)
		HandleSuccess(cmd, id, t, a)
	}
}

func policy_unsafeDeleteAttribute(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	ctx := cmd.Context()
	id := c.Flags.GetRequiredID("id")

	a, err := h.GetAttribute(ctx, id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get attribute (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	if !forceUnsafe {
		cli.ConfirmTextInput(cli.ActionDelete, "attribute", cli.InputNameFQN, a.GetFqn())
	}

	if err := h.UnsafeDeleteAttribute(ctx, id, a.GetFqn()); err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to delete attribute (%s)", id), err)
	} else {
		rows := [][]string{
			{"Deleted", "true"},
			{"Id", a.GetId()},
			{"Name", a.GetName()},
		}
		if mdRows := getMetadataRows(a.GetMetadata()); mdRows != nil {
			rows = append(rows, mdRows...)
		}
		t := cli.NewTabular(rows...)
		HandleSuccess(cmd, id, t, a)
	}
}

func policyAssignKeyToAttribute(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	attribute := c.Flags.GetRequiredString("attribute")
	keyID := c.Flags.GetRequiredID("key-id")

	// Get the attribute to show meaningful information in case of error
	attrKey, err := h.AssignKeyToAttribute(c.Context(), attribute, keyID)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to assign key: (%s) to attribute: (%s)", keyID, attribute)
		cli.ExitWithError(errMsg, err)
	}

	// Prepare and display the result
	rows := [][]string{
		{"Attribute ID", attrKey.GetAttributeId()},
		{"Key ID", attrKey.GetKeyId()},
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, attribute, t, attrKey)
}

func policyRemoveKeyFromAttribute(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	attribute := c.Flags.GetRequiredString("attribute")
	keyID := c.Flags.GetRequiredID("key-id")

	err := h.RemoveKeyFromAttribute(c.Context(), attribute, keyID)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to remove key (%s) from attribute (%s)", keyID, attribute)
		cli.ExitWithError(errMsg, err)
	}

	// Prepare and display the result
	rows := [][]string{
		{"Removed", "true"},
		{"Attribute", attribute},
		{"Key ID", keyID},
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, attribute, t, nil)
}

func init() {
	// Create an attribute
	createDoc := man.Docs.GetCommand("policy/attributes/create",
		man.WithRun(policy_createAttribute),
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("name").Name,
		createDoc.GetDocFlag("name").Shorthand,
		createDoc.GetDocFlag("name").Default,
		createDoc.GetDocFlag("name").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("rule").Name,
		createDoc.GetDocFlag("rule").Shorthand,
		createDoc.GetDocFlag("rule").Default,
		createDoc.GetDocFlag("rule").Description,
	)
	createDoc.Flags().StringSliceVarP(
		&attributeValues,
		createDoc.GetDocFlag("value").Name,
		createDoc.GetDocFlag("value").Shorthand,
		[]string{},
		createDoc.GetDocFlag("value").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("namespace").Name,
		createDoc.GetDocFlag("namespace").Shorthand,
		createDoc.GetDocFlag("namespace").Default,
		createDoc.GetDocFlag("namespace").Description,
	)
	injectLabelFlags(&createDoc.Command, false)

	// Get an attribute
	getDoc := man.Docs.GetCommand("policy/attributes/get",
		man.WithRun(policy_getAttribute),
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Shorthand,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)

	// List attributes
	listDoc := man.Docs.GetCommand("policy/attributes/list",
		man.WithRun(policy_listAttributes),
	)
	listDoc.Flags().StringP(
		listDoc.GetDocFlag("state").Name,
		listDoc.GetDocFlag("state").Shorthand,
		listDoc.GetDocFlag("state").Default,
		listDoc.GetDocFlag("state").Description,
	)
	injectListPaginationFlags(listDoc)

	// Update an attribute
	updateDoc := man.Docs.GetCommand("policy/attributes/update",
		man.WithRun(policy_updateAttribute),
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	injectLabelFlags(&updateDoc.Command, true)

	// Deactivate an attribute
	deactivateDoc := man.Docs.GetCommand("policy/attributes/deactivate",
		man.WithRun(policy_deactivateAttribute),
	)
	deactivateDoc.Flags().StringP(
		deactivateDoc.GetDocFlag("id").Name,
		deactivateDoc.GetDocFlag("id").Shorthand,
		deactivateDoc.GetDocFlag("id").Default,
		deactivateDoc.GetDocFlag("id").Description,
	)
	deactivateDoc.Flags().Bool(
		deactivateDoc.GetDocFlag("force").Name,
		false,
		deactivateDoc.GetDocFlag("force").Description,
	)

	// unsafe actions on attributes
	unsafeCmd := man.Docs.GetCommand("policy/attributes/unsafe")
	unsafeCmd.PersistentFlags().BoolVar(&forceUnsafe,
		unsafeCmd.GetDocFlag("force").Name,
		false,
		unsafeCmd.GetDocFlag("force").Description,
	)

	reactivateCmd := man.Docs.GetCommand("policy/attributes/unsafe/reactivate",
		man.WithRun(policy_unsafeReactivateAttribute),
	)
	reactivateCmd.Flags().StringP(
		reactivateCmd.GetDocFlag("id").Name,
		reactivateCmd.GetDocFlag("id").Shorthand,
		reactivateCmd.GetDocFlag("id").Default,
		reactivateCmd.GetDocFlag("id").Description,
	)
	deleteCmd := man.Docs.GetCommand("policy/attributes/unsafe/delete",
		man.WithRun(policy_unsafeDeleteAttribute),
	)
	deleteCmd.Flags().StringP(
		deleteCmd.GetDocFlag("id").Name,
		deleteCmd.GetDocFlag("id").Shorthand,
		deleteCmd.GetDocFlag("id").Default,
		deleteCmd.GetDocFlag("id").Description,
	)
	unsafeUpdateCmd := man.Docs.GetCommand("policy/attributes/unsafe/update",
		man.WithRun(policy_unsafeUpdateAttribute),
	)
	unsafeUpdateCmd.Flags().StringP(
		unsafeUpdateCmd.GetDocFlag("id").Name,
		unsafeUpdateCmd.GetDocFlag("id").Shorthand,
		unsafeUpdateCmd.GetDocFlag("id").Default,
		unsafeUpdateCmd.GetDocFlag("id").Description,
	)
	unsafeUpdateCmd.Flags().StringP(
		unsafeUpdateCmd.GetDocFlag("name").Name,
		unsafeUpdateCmd.GetDocFlag("name").Shorthand,
		unsafeUpdateCmd.GetDocFlag("name").Default,
		unsafeUpdateCmd.GetDocFlag("name").Description,
	)
	unsafeUpdateCmd.Flags().StringP(
		unsafeUpdateCmd.GetDocFlag("rule").Name,
		unsafeUpdateCmd.GetDocFlag("rule").Shorthand,
		unsafeUpdateCmd.GetDocFlag("rule").Default,
		unsafeUpdateCmd.GetDocFlag("rule").Description,
	)
	unsafeUpdateCmd.Flags().StringSliceVarP(
		&attributeValuesOrder,
		unsafeUpdateCmd.GetDocFlag("values-order").Name,
		unsafeUpdateCmd.GetDocFlag("values-order").Shorthand,
		[]string{},
		unsafeUpdateCmd.GetDocFlag("values-order").Description,
	)

	keyCmd := man.Docs.GetCommand("policy/attributes/key")

	// Assign KAS key to attribute
	assignKasKeyCmd := man.Docs.GetCommand("policy/attributes/key/assign",
		man.WithRun(policyAssignKeyToAttribute),
	)
	assignKasKeyCmd.Flags().StringP(
		assignKasKeyCmd.GetDocFlag("attribute").Name,
		assignKasKeyCmd.GetDocFlag("attribute").Shorthand,
		assignKasKeyCmd.GetDocFlag("attribute").Default,
		assignKasKeyCmd.GetDocFlag("attribute").Description,
	)
	assignKasKeyCmd.Flags().StringP(
		assignKasKeyCmd.GetDocFlag("key-id").Name,
		assignKasKeyCmd.GetDocFlag("key-id").Shorthand,
		assignKasKeyCmd.GetDocFlag("key-id").Default,
		assignKasKeyCmd.GetDocFlag("key-id").Description,
	)

	removeKasKeyCmd := man.Docs.GetCommand("policy/attributes/key/remove",
		man.WithRun(policyRemoveKeyFromAttribute),
	)
	removeKasKeyCmd.Flags().StringP(
		removeKasKeyCmd.GetDocFlag("attribute").Name,
		removeKasKeyCmd.GetDocFlag("attribute").Shorthand,
		removeKasKeyCmd.GetDocFlag("attribute").Default,
		removeKasKeyCmd.GetDocFlag("attribute").Description,
	)
	removeKasKeyCmd.Flags().StringP(
		removeKasKeyCmd.GetDocFlag("key-id").Name,
		removeKasKeyCmd.GetDocFlag("key-id").Shorthand,
		removeKasKeyCmd.GetDocFlag("key-id").Default,
		removeKasKeyCmd.GetDocFlag("key-id").Description,
	)

	keyCmd.AddSubcommands(assignKasKeyCmd, removeKasKeyCmd)
	unsafeCmd.AddSubcommands(reactivateCmd, deleteCmd, unsafeUpdateCmd)
	policy_attributesCmd.AddSubcommands(createDoc, getDoc, listDoc, updateDoc, deactivateDoc, unsafeCmd, keyCmd)
	policyCmd.AddCommand(&policy_attributesCmd.Command)
}
