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

	policy_attributesCmd = man.Docs.GetCommand("policy/attributes")
)

func policy_createAttribute(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	name := c.Flags.GetRequiredString("name")
	rule := c.Flags.GetRequiredString("rule")
	values := c.Flags.GetStringSlice("value", []string{}, cli.FlagsStringSliceOptions{})
	namespace := c.Flags.GetRequiredString("namespace")
	labels := c.Flags.GetStringSlice("label", []string{}, cli.FlagsStringSliceOptions{Min: 0})

	attr, err := h.CreateAttribute(name, rule, namespace, values, getMetadataMutable(labels))
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

	id := c.Flags.GetRequiredString("id")

	attr, err := h.GetAttribute(id)
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
	attrs, err := h.ListAttributes(state)
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
	HandleSuccess(cmd, "", t, attrs)
}

func policy_deactivateAttribute(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredString("id")

	attr, err := h.GetAttribute(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get attribute (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	cli.ConfirmAction(cli.ActionDeactivate, "attribute", attr.GetName(), false)

	attr, err = h.DeactivateAttribute(id)
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

	id := c.Flags.GetRequiredString("id")
	labels := c.Flags.GetStringSlice("label", []string{}, cli.FlagsStringSliceOptions{Min: 0})

	if a, err := h.UpdateAttribute(id, getMetadataMutable(labels), getMetadataUpdateBehavior()); err != nil {
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

	id := c.Flags.GetRequiredString("id")

	a, err := h.GetAttribute(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get attribute (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	if !forceUnsafe {
		cli.ConfirmTextInput(cli.ActionReactivate, "attribute", cli.InputNameFQN, a.GetFqn())
	}

	if reactivatedAttr, err := h.UnsafeReactivateAttribute(id); err != nil {
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

	id := c.Flags.GetRequiredString("id")
	name := c.Flags.GetOptionalString("name")
	rule := c.Flags.GetOptionalString("rule")
	valuesOrder := c.Flags.GetStringSlice("values-order", []string{}, cli.FlagsStringSliceOptions{})

	a, err := h.GetAttribute(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get attribute (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	if !forceUnsafe {
		cli.ConfirmTextInput(cli.ActionUpdateUnsafe, "attribute", cli.InputNameFQN, a.GetFqn())
	}

	if err := h.UnsafeUpdateAttribute(id, name, rule, valuesOrder); err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to update attribute (%s)", id), err)
	} else {
		var (
			values   []string
			valueIDs []string
		)
		for _, v := range a.GetValues() {
			values = append(values, v.GetValue())
			valueIDs = append(valueIDs, v.GetId())
		}
		rows := [][]string{
			{"Id", a.GetId()},
			{"Name", a.GetName()},
			{"Rule", handlers.GetAttributeRuleFromAttributeType(a.GetRule())},
			{"Values", cli.CommaSeparated(values)},
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

	id := c.Flags.GetRequiredString("id")

	a, err := h.GetAttribute(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get attribute (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	if !forceUnsafe {
		cli.ConfirmTextInput(cli.ActionDelete, "attribute", cli.InputNameFQN, a.GetFqn())
	}

	if err := h.UnsafeDeleteAttribute(id, a.GetFqn()); err != nil {
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
	createDoc.Flags().StringSliceP(
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
	unsafeUpdateCmd.Flags().StringSliceP(
		unsafeUpdateCmd.GetDocFlag("values-order").Name,
		unsafeUpdateCmd.GetDocFlag("values-order").Shorthand,
		[]string{},
		unsafeUpdateCmd.GetDocFlag("values-order").Description,
	)

	unsafeCmd.AddSubcommands(reactivateCmd, deleteCmd, unsafeUpdateCmd)
	policy_attributesCmd.AddSubcommands(createDoc, getDoc, listDoc, updateDoc, deactivateDoc, unsafeCmd)
	policyCmd.AddCommand(&policy_attributesCmd.Command)
}
