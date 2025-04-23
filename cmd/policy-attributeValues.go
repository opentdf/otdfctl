package cmd

import (
	"fmt"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
)

var policy_attributeValuesCmd *cobra.Command

func policy_createAttributeValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	attrId := c.FlagHelper.GetRequiredID("attribute-id")
	value := c.FlagHelper.GetRequiredString("value")
	metadataLabels = c.FlagHelper.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	attr, err := h.GetAttribute(attrId)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get parent attribute (%s)", attrId), err)
	}

	v, err := h.CreateAttributeValue(attr.GetId(), value, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create attribute value", err)
	}

	handleValueSuccess(cmd, v)
}

func policy_getAttributeValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.FlagHelper.GetRequiredID("id")

	v, err := h.GetAttributeValue(id)
	if err != nil {
		cli.ExitWithError("Failed to find attribute value", err)
	}

	handleValueSuccess(cmd, v)
}

func policy_listAttributeValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	attrId := c.FlagHelper.GetRequiredID("attribute-id")
	state := cli.GetState(cmd)
	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")

	vals, page, err := h.ListAttributeValues(attrId, state, limit, offset)
	if err != nil {
		cli.ExitWithError("Failed to list attribute values", err)
	}
	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("fqn", "Fqn", cli.FlexColumnWidthFour),
		table.NewFlexColumn("active", "Active", cli.FlexColumnWidthThree),
		table.NewFlexColumn("labels", "Labels", cli.FlexColumnWidthOne),
		table.NewFlexColumn("created_at", "Created At", cli.FlexColumnWidthOne),
		table.NewFlexColumn("updated_at", "Updated At", cli.FlexColumnWidthOne),
	)
	rows := []table.Row{}
	for _, val := range vals {
		v := cli.GetSimpleAttributeValue(val)
		rows = append(rows, table.NewRow(table.RowData{
			"id":         v.Id,
			"fqn":        v.FQN,
			"active":     v.Active,
			"labels":     v.Metadata["Labels"],
			"created_at": v.Metadata["Created At"],
			"updated_at": v.Metadata["Updated At"],
		}))
	}
	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, vals)
}

func policy_updateAttributeValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	_, err := h.GetAttributeValue(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get attribute value (%s)", id), err)
	}

	v, err := h.UpdateAttributeValue(id, getMetadataMutable(metadataLabels), getMetadataUpdateBehavior())
	if err != nil {
		cli.ExitWithError("Failed to update attribute value", err)
	}

	handleValueSuccess(cmd, v)
}

func policy_deactivateAttributeValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	value, err := h.GetAttributeValue(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get attribute value (%s)", id), err)
	}

	cli.ConfirmAction(cli.ActionDeactivate, "attribute value", value.GetValue(), false)

	deactivated, err := h.DeactivateAttributeValue(id)
	if err != nil {
		cli.ExitWithError("Failed to deactivate attribute value", err)
	}

	handleValueSuccess(cmd, deactivated)
}

func policy_unsafeReactivateAttributeValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	v, err := h.GetAttributeValue(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get attribute value (%s)", id), err)
	}

	if !forceUnsafe {
		cli.ConfirmTextInput(cli.ActionReactivate, "attribute value", cli.InputNameFQN, v.GetFqn())
	}

	if reactivated, err := h.UnsafeReactivateAttributeValue(id); err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to reactivate attribute value (%s)", id), err)
	} else {
		rows := [][]string{
			{"Id", reactivated.GetId()},
			{"Value", reactivated.GetValue()},
		}
		if mdRows := getMetadataRows(v.GetMetadata()); mdRows != nil {
			rows = append(rows, mdRows...)
		}
		t := cli.NewTabular(rows...)
		HandleSuccess(cmd, id, t, v)
	}
}

func policy_unsafeUpdateAttributeValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	value := c.Flags.GetOptionalString("value")

	v, err := h.GetAttributeValue(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get attribute value (%s)", id), err)
	}

	if !forceUnsafe {
		cli.ConfirmTextInput(cli.ActionUpdateUnsafe, "attribute value", cli.InputNameFQN, v.GetFqn())
	}

	if err := h.UnsafeUpdateAttributeValue(id, value); err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to update attribute value (%s)", id), err)
	} else {
		rows := [][]string{
			{"Id", v.GetId()},
			{"Value", value},
		}
		if mdRows := getMetadataRows(v.GetMetadata()); mdRows != nil {
			rows = append(rows, mdRows...)
		}
		t := cli.NewTabular(rows...)
		HandleSuccess(cmd, id, t, v)
	}
}

func policy_unsafeDeleteAttributeValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	v, err := h.GetAttributeValue(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get attribute value (%s)", id), err)
	}

	if !forceUnsafe {
		cli.ConfirmTextInput(cli.ActionDelete, "attribute value", cli.InputNameFQN, v.GetFqn())
	}

	if err := h.UnsafeDeleteAttributeValue(id, v.GetFqn()); err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to delete attribute (%s)", id), err)
	} else {
		rows := [][]string{
			{"Id", v.GetId()},
			{"Value", v.GetValue()},
			{"Deleted", "true"},
		}
		if mdRows := getMetadataRows(v.GetMetadata()); mdRows != nil {
			rows = append(rows, mdRows...)
		}
		t := cli.NewTabular(rows...)
		HandleSuccess(cmd, id, t, v)
	}
}

func policy_assignKASKeyToAttributeValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	value := c.Flags.GetRequiredString("value")
	keyId := c.Flags.GetRequiredID("keyId")

	attrKey, err := h.AssignKeyToAttributeValue(c.Context(), value, keyId)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to assign key: (%s) to attribute value: (%s)", keyId, value)
		cli.ExitWithError(errMsg, err)
	}

	rows := [][]string{
		{"Value ID", attrKey.GetValueId()},
		{"Key ID", attrKey.GetKeyId()},
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, value, t, attrKey)
}

func policy_removeKASKeyFromAttributeValue(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	value := c.Flags.GetRequiredString("value")
	keyId := c.Flags.GetRequiredID("keyId")

	err := h.RemoveKeyFromAttributeValue(c.Context(), value, keyId)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to remove key (%s) from attribute value (%s)", keyId, value)
		cli.ExitWithError(errMsg, err)
	}

	rows := [][]string{
		{"Removed", "true"},
		{"Value", value},
		{"Key ID", keyId},
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, value, t, nil)
}

func init() {
	createCmd := man.Docs.GetCommand("policy/attributes/values/create",
		man.WithRun(policy_createAttributeValue),
	)
	createCmd.Flags().StringP(
		createCmd.GetDocFlag("attribute-id").Name,
		createCmd.GetDocFlag("attribute-id").Shorthand,
		createCmd.GetDocFlag("attribute-id").Default,
		createCmd.GetDocFlag("attribute-id").Description,
	)
	createCmd.Flags().StringP(
		createCmd.GetDocFlag("value").Name,
		createCmd.GetDocFlag("value").Shorthand,
		createCmd.GetDocFlag("value").Default,
		createCmd.GetDocFlag("value").Description,
	)
	injectLabelFlags(&createCmd.Command, false)

	getCmd := man.Docs.GetCommand("policy/attributes/values/get",
		man.WithRun(policy_getAttributeValue),
	)
	getCmd.Flags().StringP(
		getCmd.GetDocFlag("id").Name,
		getCmd.GetDocFlag("id").Shorthand,
		getCmd.GetDocFlag("id").Default,
		getCmd.GetDocFlag("id").Description,
	)

	listCmd := man.Docs.GetCommand("policy/attributes/values/list",
		man.WithRun(policy_listAttributeValue),
	)
	listCmd.Flags().StringP(
		listCmd.GetDocFlag("attribute-id").Name,
		listCmd.GetDocFlag("attribute-id").Shorthand,
		listCmd.GetDocFlag("attribute-id").Default,
		listCmd.GetDocFlag("attribute-id").Description,
	)
	listCmd.Flags().StringP(
		listCmd.GetDocFlag("state").Name,
		listCmd.GetDocFlag("state").Shorthand,
		listCmd.GetDocFlag("state").Default,
		listCmd.GetDocFlag("state").Description,
	)
	injectListPaginationFlags(listCmd)

	updateCmd := man.Docs.GetCommand("policy/attributes/values/update",
		man.WithRun(policy_updateAttributeValue),
	)
	updateCmd.Flags().StringP(
		updateCmd.GetDocFlag("id").Name,
		updateCmd.GetDocFlag("id").Shorthand,
		updateCmd.GetDocFlag("id").Default,
		updateCmd.GetDocFlag("id").Description,
	)
	injectLabelFlags(&updateCmd.Command, true)

	deactivateCmd := man.Docs.GetCommand("policy/attributes/values/deactivate",
		man.WithRun(policy_deactivateAttributeValue),
	)
	deactivateCmd.Flags().StringP(
		deactivateCmd.GetDocFlag("id").Name,
		deactivateCmd.GetDocFlag("id").Shorthand,
		deactivateCmd.GetDocFlag("id").Default,
		deactivateCmd.GetDocFlag("id").Description,
	)

	keyCmd := man.Docs.GetCommand("policy/attributes/values/key")

	assignKasKeyCmd := man.Docs.GetCommand("policy/attributes/values/key/assign",
		man.WithRun(policy_assignKASKeyToAttributeValue),
	)
	assignKasKeyCmd.Flags().StringP(
		assignKasKeyCmd.GetDocFlag("value").Name,
		assignKasKeyCmd.GetDocFlag("value").Shorthand,
		assignKasKeyCmd.GetDocFlag("value").Default,
		assignKasKeyCmd.GetDocFlag("value").Description,
	)
	assignKasKeyCmd.Flags().StringP(
		assignKasKeyCmd.GetDocFlag("keyId").Name,
		assignKasKeyCmd.GetDocFlag("keyId").Shorthand,
		assignKasKeyCmd.GetDocFlag("keyId").Default,
		assignKasKeyCmd.GetDocFlag("keyId").Description,
	)

	removeKasKeyCmd := man.Docs.GetCommand("policy/attributes/values/key/remove",
		man.WithRun(policy_removeKASKeyFromAttributeValue),
	)
	removeKasKeyCmd.Flags().StringP(
		removeKasKeyCmd.GetDocFlag("value").Name,
		removeKasKeyCmd.GetDocFlag("value").Shorthand,
		removeKasKeyCmd.GetDocFlag("value").Default,
		removeKasKeyCmd.GetDocFlag("value").Description,
	)
	removeKasKeyCmd.Flags().StringP(
		removeKasKeyCmd.GetDocFlag("keyId").Name,
		removeKasKeyCmd.GetDocFlag("keyId").Shorthand,
		removeKasKeyCmd.GetDocFlag("keyId").Default,
		removeKasKeyCmd.GetDocFlag("keyId").Description,
	)

	unsafeReactivateCmd := man.Docs.GetCommand("policy/attributes/values/unsafe/reactivate",
		man.WithRun(policy_unsafeReactivateAttributeValue),
	)
	unsafeReactivateCmd.Flags().StringP(
		unsafeReactivateCmd.GetDocFlag("id").Name,
		unsafeReactivateCmd.GetDocFlag("id").Shorthand,
		unsafeReactivateCmd.GetDocFlag("id").Default,
		unsafeReactivateCmd.GetDocFlag("id").Description,
	)

	unsafeDeleteCmd := man.Docs.GetCommand("policy/attributes/values/unsafe/delete",
		man.WithRun(policy_unsafeDeleteAttributeValue),
	)
	unsafeDeleteCmd.Flags().StringP(
		unsafeDeleteCmd.GetDocFlag("id").Name,
		unsafeDeleteCmd.GetDocFlag("id").Shorthand,
		unsafeDeleteCmd.GetDocFlag("id").Default,
		unsafeDeleteCmd.GetDocFlag("id").Description,
	)

	unsafeUpdateCmd := man.Docs.GetCommand("policy/attributes/values/unsafe/update",
		man.WithRun(policy_unsafeUpdateAttributeValue),
	)
	unsafeUpdateCmd.Flags().StringP(
		unsafeUpdateCmd.GetDocFlag("id").Name,
		unsafeUpdateCmd.GetDocFlag("id").Shorthand,
		unsafeUpdateCmd.GetDocFlag("id").Default,
		unsafeUpdateCmd.GetDocFlag("id").Description,
	)
	unsafeUpdateCmd.Flags().StringP(
		unsafeUpdateCmd.GetDocFlag("value").Name,
		unsafeUpdateCmd.GetDocFlag("value").Shorthand,
		unsafeUpdateCmd.GetDocFlag("value").Default,
		unsafeUpdateCmd.GetDocFlag("value").Description,
	)

	unsafeCmd := man.Docs.GetCommand("policy/attributes/values/unsafe")
	unsafeCmd.PersistentFlags().BoolVar(&forceUnsafe,
		unsafeCmd.GetDocFlag("force").Name,
		false,
		unsafeCmd.GetDocFlag("force").Description,
	)

	keyCmd.AddSubcommands(assignKasKeyCmd, removeKasKeyCmd)
	unsafeCmd.AddSubcommands(unsafeReactivateCmd, unsafeDeleteCmd, unsafeUpdateCmd)
	doc := man.Docs.GetCommand("policy/attributes/values",
		man.WithSubcommands(createCmd, getCmd, listCmd, updateCmd, deactivateCmd, unsafeCmd, keyCmd),
	)
	policy_attributeValuesCmd = &doc.Command
	policy_attributesCmd.AddCommand(policy_attributeValuesCmd)
}

func handleValueSuccess(cmd *cobra.Command, v *policy.Value) {
	keyIds := make([]string, len(v.GetKeys()))
	for i, k := range v.GetKeys() {
		keyIds[i] = k.GetId()
	}

	rows := [][]string{
		{"Id", v.GetId()},
		{"FQN", v.GetFqn()},
		{"Value", v.GetValue()},
		{"Associated Keys", cli.CommaSeparated(keyIds)},
	}
	if mdRows := getMetadataRows(v.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, v.GetId(), t, v)
}
