package cmd

import (
	"fmt"

	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/tructl/issues/73] is addressed

func handleValueSuccess(cmd *cobra.Command, v *policy.Value) {
	rows := [][]string{
		{"Id", v.Id},
		{"FQN", v.Fqn},
		{"Value", v.Value},
	}
	if len(v.Members) > 0 {
		memberIds := make([]string, len(v.Members))
		for i, m := range v.Members {
			memberIds[i] = m.Id
		}
		rows = append(rows, []string{"Members", cli.CommaSeparated(memberIds)})
	}
	t := cli.NewTabular().Rows(rows...)
	HandleSuccess(cmd, v.Id, t, v)
}

var (
	policy_attributeValuesCmd = &cobra.Command{
		Use:   "values",
		Short: "Manage attribute values",
	}

	policy_attributeValuesCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create an attribute value",
		Run: func(cmd *cobra.Command, args []string) {
			flagHelper := cli.NewFlagHelper(cmd)
			attrId := flagHelper.GetRequiredString("attribute-id")
			value := flagHelper.GetRequiredString("value")
			metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

			h := cli.NewHandler(cmd)
			defer h.Close()

			attr, err := h.GetAttribute(attrId)
			if err != nil {
				cli.ExitWithNotFoundError("Could not find attribute", err)
			}

			v, err := h.CreateAttributeValue(attr.Id, value, getMetadataMutable(metadataLabels))
			if err != nil {
				cli.ExitWithError("Could not create attribute value", err)
			}

			handleValueSuccess(cmd, v)
		},
	}

	policy_attributeValuesGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get an attribute value",
		Run: func(cmd *cobra.Command, args []string) {
			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			h := cli.NewHandler(cmd)
			defer h.Close()

			v, err := h.GetAttributeValue(id)
			if err != nil {
				cli.ExitWithNotFoundError("Could not find attribute value", err)
			}

			handleValueSuccess(cmd, v)
		},
	}

	policy_attributeValuesUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an attribute value",
		Run: func(cmd *cobra.Command, args []string) {
			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")
			metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

			h := cli.NewHandler(cmd)
			defer h.Close()

			_, err := h.GetAttributeValue(id)
			if err != nil {
				cli.ExitWithNotFoundError("Could not find attribute value", err)
			}

			v, err := h.UpdateAttributeValue(id, nil, getMetadataMutable(metadataLabels), getMetadataUpdateBehavior())
			if err != nil {
				cli.ExitWithError("Could not update attribute value", err)
			}

			handleValueSuccess(cmd, v)
		},
	}

	policy_attributeValuesDeactivateCmd = &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivate an attribute value",
		Run: func(cmd *cobra.Command, args []string) {
			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			h := cli.NewHandler(cmd)
			defer h.Close()

			value, err := h.GetAttributeValue(id)
			if err != nil {
				cli.ExitWithNotFoundError("Could not find attribute value", err)
			}

			cli.ConfirmAction(cli.ActionDeactivate, "attribute value", value.Value)

			err = h.DeactivateAttributeValue(id)
			if err != nil {
				cli.ExitWithError("Could not deactivate attribute value", err)
			}

			handleValueSuccess(cmd, value)
		},
	}

	///
	/// Attribute Value Members
	///
	attrValueMembers = []string{}

	policy_attributeValueMembersCmd = &cobra.Command{
		Use:   "members",
		Short: "Manage attribute value members",
		Long:  "Manage attribute value members",
	}

	// Add member to attribute value
	policy_attributeValueMembersAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Add members to an attribute value",
		Run: func(cmd *cobra.Command, args []string) {
			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")
			members := flagHelper.GetStringSlice("member", attrValueMembers, cli.FlagHelperStringSliceOptions{})

			// TODO: Implement
			fmt.Println("Not implemented")
			fmt.Printf("id: %s\n", id)
			fmt.Printf("members: %v\n", members)
		},
	}

	// Remove member from attribute value
	policy_attributeValueMembersRemoveCmd = &cobra.Command{
		Use:   "remove",
		Short: "Remove members from an attribute value",
		Run: func(cmd *cobra.Command, args []string) {
			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")
			members := flagHelper.GetStringSlice("members", attrValueMembers, cli.FlagHelperStringSliceOptions{})

			// TODO: Implement
			fmt.Println("Not implemented")
			fmt.Printf("id: %s\n", id)
			fmt.Printf("members: %v\n", members)
		},
	}

	// Replace members of attribute value
	policy_attributeValueMembersReplaceCmd = &cobra.Command{
		Use:   "replace",
		Short: "Replace members from an attribute value",
		Long:  "This command will replace the members of an attribute value with the provided members. ",
		Run: func(cmd *cobra.Command, args []string) {
			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")
			members := flagHelper.GetStringSlice("members", attrValueMembers, cli.FlagHelperStringSliceOptions{})

			// TODO: Implement
			fmt.Println("Not implemented")
			fmt.Printf("id: %s\n", id)
			fmt.Printf("members: %v\n", members)
		},
	}
)

func init() {
	policy_attributesCmd.AddGroup(
		&cobra.Group{
			ID:    "subcommand",
			Title: "Subcommands",
		},
	)
	policy_attributesCmd.AddCommand(policy_attributeValuesCmd)
	policy_attributeValuesCmd.GroupID = "subcommand"
	policy_attributeValuesCmd.AddGroup(
		&cobra.Group{
			ID:    "subcommand",
			Title: "Subcommands",
		},
	)

	policy_attributeValuesCmd.AddCommand(policy_attributeValuesCreateCmd)
	policy_attributeValuesCreateCmd.Flags().StringP("attribute-id", "a", "", "Attribute id")
	policy_attributeValuesCreateCmd.Flags().StringP("value", "v", "", "Value")

	policy_attributeValuesCmd.AddCommand(policy_attributeValuesGetCmd)
	policy_attributeValuesGetCmd.Flags().StringP("id", "i", "", "Attribute value id")

	policy_attributeValuesCmd.AddCommand(policy_attributeValuesUpdateCmd)
	policy_attributeValuesUpdateCmd.Flags().StringP("id", "i", "", "Attribute value id")
	policy_attributeValuesUpdateCmd.Flags().StringP("value", "v", "", "Value")

	policy_attributeValuesCmd.AddCommand(policy_attributeValuesDeactivateCmd)
	policy_attributeValuesDeactivateCmd.Flags().StringP("id", "i", "", "Attribute value id")

	// Attribute value members
	policy_attributeValuesCmd.AddCommand(policy_attributeValueMembersCmd)
	policy_attributeValueMembersCmd.GroupID = "subcommand"

	policy_attributeValueMembersCmd.AddCommand(policy_attributeValueMembersAddCmd)
	policy_attributeValueMembersAddCmd.Flags().StringP("id", "i", "", "Attribute value id")
	policy_attributeValueMembersAddCmd.Flags().StringSliceVar(&attrValueMembers, "members", []string{}, "Members to add")

	policy_attributeValueMembersCmd.AddCommand(policy_attributeValueMembersRemoveCmd)
	policy_attributeValueMembersRemoveCmd.Flags().StringP("id", "i", "", "Attribute value id")
	policy_attributeValueMembersRemoveCmd.Flags().StringSliceVar(&attrValueMembers, "members", []string{}, "Members to add")

	policy_attributeValueMembersCmd.AddCommand(policy_attributeValueMembersReplaceCmd)
	policy_attributeValueMembersReplaceCmd.Flags().StringP("id", "i", "", "Attribute value id")
	policy_attributeValueMembersReplaceCmd.Flags().StringSliceVar(&attrValueMembers, "members", []string{}, "Members to add")
}
