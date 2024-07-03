package cmd

import (
	"fmt"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/otdfctl/issues/73] is addressed

var policy_attributeValuesCmd *cobra.Command

func policy_createAttributeValue(cmd *cobra.Command, args []string) {
	flagHelper := cli.NewFlagHelper(cmd)
	attrId := flagHelper.GetRequiredString("attribute-id")
	value := flagHelper.GetRequiredString("value")
	metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})
	// TODO: support create with members when update is unblocked to remove/alter them after creation [https://github.com/opentdf/platform/issues/476]

	h := NewHandler(cmd)
	defer h.Close()

	attr, err := h.GetAttribute(attrId)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get parent attribute (%s)", attrId), err)
	}

	v, err := h.CreateAttributeValue(attr.Id, value, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create attribute value", err)
	}

	handleValueSuccess(cmd, v)
}

func policy_getAttributeValue(cmd *cobra.Command, args []string) {
	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	h := NewHandler(cmd)
	defer h.Close()

	v, err := h.GetAttributeValue(id)
	if err != nil {
		cli.ExitWithError("Failed to find attribute value", err)
	}

	handleValueSuccess(cmd, v)
}

func policy_listAttributeValue(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()
	flagHelper := cli.NewFlagHelper(cmd)
	attrId := flagHelper.GetRequiredString("attribute-id")
	state := cli.GetState(cmd)
	vals, err := h.ListAttributeValues(attrId, state)
	if err != nil {
		cli.ExitWithError("Failed to list attribute values", err)
	}
	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewColumn("fqn", "Fqn", 16),
		table.NewColumn("members", "Members", 16),
		table.NewColumn("active", "Active", 16),
		table.NewColumn("labels", "Labels", 16),
		table.NewColumn("created_at", "Created At", 16),
		table.NewColumn("updated_at", "Updated At", 16),
	)
	rows := []table.Row{}
	for _, val := range vals {
		v := cli.GetSimpleAttributeValue(val)
		rows = append(rows, table.NewRow(table.RowData{
			"id":         v.Id,
			"fqn":        v.FQN,
			"members":    cli.CommaSeparated(v.Members),
			"active":     v.Active,
			"labels":     v.Metadata["Labels"],
			"created_at": v.Metadata["Created At"],
			"updated_at": v.Metadata["Updated At"],
		}))
	}
	HandleSuccess(cmd, "", t, vals)
}

func policy_updateAttributeValue(cmd *cobra.Command, args []string) {
	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")
	metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

	h := NewHandler(cmd)
	defer h.Close()

	_, err := h.GetAttributeValue(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get attribute value (%s)", id), err)
	}

	v, err := h.UpdateAttributeValue(id, nil, getMetadataMutable(metadataLabels), getMetadataUpdateBehavior())
	if err != nil {
		cli.ExitWithError("Failed to update attribute value", err)
	}

	handleValueSuccess(cmd, v)
}

func policy_deactivateAttributeValue(cmd *cobra.Command, args []string) {
	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	h := NewHandler(cmd)
	defer h.Close()

	value, err := h.GetAttributeValue(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get attribute value (%s)", id), err)
	}

	cli.ConfirmAction(cli.ActionDeactivate, "attribute value", value.Value)

	deactivated, err := h.DeactivateAttributeValue(id)
	if err != nil {
		cli.ExitWithError("Failed to deactivate attribute value", err)
	}

	handleValueSuccess(cmd, deactivated)
}

func policy_unsafeReactivateAttributeValue(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	// TODO: confirm action here!

	if v, err := h.UnsafeReactivateAttributeValue(id); err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to reactivate attribute value (%s)", id), err)
	} else {
		rows := [][]string{
			{"Id", v.GetId()},
			{"Name", v.GetName()},
		}
		if mdRows := getMetadataRows(v.GetMetadata()); mdRows != nil {
			rows = append(rows, mdRows...)
		}
		t := cli.NewTabular(rows...)
		HandleSuccess(cmd, id, t, v)
	}
}

func policy_unsafeUpdateAttributeValue(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")
	value := flagHelper.GetOptionalString("value")

	v, err := h.GetAttributeValue(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get attribute value (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	// TODO: confirm action here!

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
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	v, err := h.GetAttributeValue(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get attribute value (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	// TODO: confirm action here!

	if err := h.UnsafeDeleteAttributeValue(id); err != nil {
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

// TODO: uncomment when update with members is enabled in the platform [https://github.com/opentdf/platform/issues/476]
///
/// Attribute Value Members
///
// attrValueMembers = []string{}

// policy_attributeValueMembersCmd = &cobra.Command{
// 	Use:   "members",
// 	Short: "Manage attribute value members",
// 	Long:  "Manage attribute value members",
// }

// // Add member to attribute value
// policy_attributeValueMembersAddCmd = &cobra.Command{
// 	Use:   "add",
// 	Short: "Add members to an attribute value",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		flagHelper := cli.NewFlagHelper(cmd)
// 		id := flagHelper.GetRequiredString("id")
// 		members := flagHelper.GetStringSlice("member", attrValueMembers, cli.FlagHelperStringSliceOptions{})

// 		h := NewHandler(cmd)
// 		defer h.Close()

// 		prev, err := h.GetAttributeValue(id)
// 		if err != nil {
// 			cli.ExitWithError(fmt.Sprintf("Failed to get attribute value (%s)", id), err)
// 		}

// 		action := fmt.Sprintf("%s [%s] to", cli.ActionMemberAdd, strings.Join(members, ", "))
// 		cli.ConfirmAction(action, "attribute value", id)

// 		prevMemberIds := make([]string, len(prev.Members))
// 		for i, m := range prev.Members {
// 			prevMemberIds[i] = m.GetId()
// 		}
// 		updated := append(prevMemberIds, members...)

// 		v, err := h.UpdateAttributeValue(id, updated, nil, common.MetadataUpdateEnum_METADATA_UPDATE_ENUM_UNSPECIFIED)
// 		if err != nil {
// 			cli.ExitWithError(fmt.Sprintf("Failed to %s [%s] to attribute value (%s)", cli.ActionMemberAdd, strings.Join(members, ", "), id), err)
// 		}

// 		handleValueSuccess(cmd, v)
// 	},
// }

// // Remove member from attribute value
// policy_attributeValueMembersRemoveCmd = &cobra.Command{
// 	Use:   "remove",
// 	Short: "Remove members from an attribute value",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		flagHelper := cli.NewFlagHelper(cmd)
// 		id := flagHelper.GetRequiredString("id")
// 		members := flagHelper.GetStringSlice("members", attrValueMembers, cli.FlagHelperStringSliceOptions{})

// 		h := NewHandler(cmd)
// 		defer h.Close()

// 		prev, err := h.GetAttributeValue(id)
// 		if err != nil {
// 			cli.ExitWithError(fmt.Sprintf("Failed to get attribute value (%s)", id), err)
// 		}

// 		action := fmt.Sprintf("%s [%s] from", cli.ActionMemberRemove, strings.Join(members, ", "))
// 		cli.ConfirmAction(action, "attribute value", id)

// 		// collect the member ids off the members, then make the removals
// 		updatedMemberIds := make([]string, len(prev.Members))
// 		for i, m := range prev.Members {
// 			updatedMemberIds[i] = m.GetId()
// 		}
// 		for _, toBeRemoved := range members {
// 			for i, str := range updatedMemberIds {
// 				if toBeRemoved == str {
// 					updatedMemberIds = append(updatedMemberIds[:i], updatedMemberIds[i+1:]...)
// 					break
// 				}
// 			}
// 		}

// 		v, err := h.UpdateAttributeValue(id, updatedMemberIds, nil, common.MetadataUpdateEnum_METADATA_UPDATE_ENUM_UNSPECIFIED)
// 		if err != nil {
// 			cli.ExitWithError(fmt.Sprintf("Failed to %s [%s] from attribute value (%s)", cli.ActionMemberRemove, strings.Join(members, ", "), id), err)
// 		}

// 		handleValueSuccess(cmd, v)
// 	},
// }

// // Replace members of attribute value
// policy_attributeValueMembersReplaceCmd = &cobra.Command{
// 	Use:   "replace",
// 	Short: "Replace members from an attribute value",
// 	Long:  "This command will replace the members of an attribute value with the provided members. ",
// 	Run: func(cmd *cobra.Command, args []string) {
// 		flagHelper := cli.NewFlagHelper(cmd)
// 		id := flagHelper.GetRequiredString("id")
// 		members := flagHelper.GetStringSlice("members", attrValueMembers, cli.FlagHelperStringSliceOptions{})

// 		h := NewHandler(cmd)
// 		defer h.Close()

// 		prev, err := h.GetAttributeValue(id)
// 		if err != nil {
// 			cli.ExitWithError(fmt.Sprintf("Failed to find attribute value (%s)", id), err)
// 		}

// 		existingMemberIds := make([]string, len(prev.Members))
// 		for i, m := range prev.Members {
// 			existingMemberIds[i] = m.GetId()
// 		}

// 		action := fmt.Sprintf("%s [%s] with [%s] under", cli.ActionMemberReplace, strings.Join(existingMemberIds, ", "), strings.Join(members, ", "))
// 		cli.ConfirmAction(action, "attribute value", id)

// 		v, err := h.UpdateAttributeValue(id, members, nil, common.MetadataUpdateEnum_METADATA_UPDATE_ENUM_UNSPECIFIED)
// 		if err != nil {
// 			cli.ExitWithError(fmt.Sprintf("Failed to %s of attribute value (%s)", cli.ActionMemberReplace, id), err)
// 		}

// 		handleValueSuccess(cmd, v)
// 	},
// }

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
	// unsafe
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

	// Attribute value members
	// policy_attributeValuesCmd.AddCommand(policy_attributeValueMembersCmd)
	// policy_attributeValueMembersCmd.GroupID = "subcommand"

	// policy_attributeValueMembersCmd.AddCommand(policy_attributeValueMembersAddCmd)
	// policy_attributeValueMembersAddCmd.Flags().StringP("id", "i", "", "Attribute value id")
	// policy_attributeValueMembersAddCmd.Flags().StringSliceVar(&attrValueMembers, "member", []string{}, "Each member id to add")

	// policy_attributeValueMembersCmd.AddCommand(policy_attributeValueMembersRemoveCmd)
	// policy_attributeValueMembersRemoveCmd.Flags().StringP("id", "i", "", "Attribute value id")
	// policy_attributeValueMembersRemoveCmd.Flags().StringSliceVar(&attrValueMembers, "member", []string{}, "Each member id to remove")

	// policy_attributeValueMembersCmd.AddCommand(policy_attributeValueMembersReplaceCmd)
	// policy_attributeValueMembersReplaceCmd.Flags().StringP("id", "i", "", "Attribute value id")
	// policy_attributeValueMembersReplaceCmd.Flags().StringSliceVar(&attrValueMembers, "member", []string{}, "Each member id that should exist after replacement")

	unsafeCmd := man.Docs.GetCommand("policy/attributes/values/unsafe")
	unsafeCmd.PersistentFlags().BoolVar(&forceUnsafe,
		unsafeCmd.GetDocFlag("force").Name,
		false,
		unsafeCmd.GetDocFlag("force").Description,
	)

	unsafeCmd.AddSubcommands(unsafeReactivateCmd, unsafeDeleteCmd, unsafeUpdateCmd)
	doc := man.Docs.GetCommand("policy/attributes/values",
		man.WithSubcommands(createCmd, getCmd, listCmd, updateCmd, deactivateCmd, unsafeCmd),
	)
	policy_attributeValuesCmd = &doc.Command
	policy_attributesCmd.AddCommand(policy_attributeValuesCmd)
}

func handleValueSuccess(cmd *cobra.Command, v *policy.Value) {
	rows := [][]string{
		{"Id", v.GetId()},
		{"FQN", v.GetFqn()},
		{"Value", v.GetValue()},
	}
	if mdRows := getMetadataRows(v.Metadata); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	members := v.GetMembers()
	if len(members) > 0 {
		memberIds := make([]string, len(members))
		for i, m := range members {
			memberIds[i] = m.Id
		}
		rows = append(rows, []string{"Members", cli.CommaSeparated(memberIds)})
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, v.Id, t, v)
}
