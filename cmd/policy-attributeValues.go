package cmd

import (
	"github.com/opentdf/tructl/pkg/man"
	"github.com/spf13/cobra"
)

var (
	policy_attributeValuesCmd = &cobra.Command{
		Use:     man.Docs.GetDoc("policy/attributes/values").Use,
		Aliases: man.Docs.GetDoc("policy/attributes/values").Aliases,
		Short:   man.Docs.GetDoc("policy/attributes/values").Short,
		Long:    man.Docs.GetDoc("policy/attributes/values").Long,
	}

	// policy_attributeValuesCreateCmd = &cobra.Command{
	// 	Use:   man.GetDoc("policy-attributeValues-create").Use,
	// 	Short: man.GetDoc("policy-attributeValues-create").Short,
	// 	Long:  man.GetDoc("policy-attributeValues-create").Long,
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		flagHelper := cli.NewFlagHelper(cmd)
	// 		attrId := flagHelper.GetRequiredString("attribute-id")
	// 		value := flagHelper.GetRequiredString("value")

	// 		h := cli.NewHandler(cmd)
	// 		defer h.Close()

	// 		attr, err := h.GetAttribute(attrId)
	// 		if err != nil {
	// 			errMsg := "Could not find attribute"
	// 			cli.ExitWithNotFoundError(errMsg, err)
	// 			cli.ExitWithError(errMsg, err)
	// 		}

	// 		newValue, err := h.CreateAttributeValue(attr.Id, value)
	// 		if err != nil {
	// 			errMsg := "Could not create attribute value"
	// 			cli.ExitWithNotFoundError(errMsg, err)
	// 			cli.ExitWithError(errMsg, err)
	// 		}

	// 		v := cli.GetSimpleAttributeValue(newValue)
	// 		fmt.Println(cli.SuccessMessage("Attribute value created"))
	// 		fmt.Println(
	// 			cli.NewTabular().
	// 				Rows([][]string{
	// 					{"Id", v.Id},
	// 					{"FQN", v.FQN},
	// 					{"Members", cli.CommaSeparated(v.Members)},
	// 				}...).Render(),
	// 		)
	// 	},
	// }

	// policy_attributeValuesGetCmd = &cobra.Command{
	// 	Use:   "get",
	// 	Short: "Get an attribute value",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		flagHelper := cli.NewFlagHelper(cmd)
	// 		id := flagHelper.GetRequiredString("id")

	// 		h := cli.NewHandler(cmd)
	// 		defer h.Close()

	// 		value, err := h.GetAttributeValue(id)
	// 		if err != nil {
	// 			errMsg := "Could not find attribute value"
	// 			cli.ExitWithNotFoundError(errMsg, err)
	// 			cli.ExitWithError(errMsg, err)
	// 		}

	// 		v := cli.GetSimpleAttributeValue(value)
	// 		fmt.Println(
	// 			cli.NewTabular().
	// 				Rows([][]string{
	// 					{"Id", v.Id},
	// 					{"FQN", v.FQN},
	// 					{"Members", cli.CommaSeparated(v.Members)},
	// 				}...).Render(),
	// 		)
	// 	},
	// }

	// policy_attributeValuesUpdateCmd = &cobra.Command{
	// 	Use:   "update",
	// 	Short: "Update an attribute value",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		flagHelper := cli.NewFlagHelper(cmd)
	// 		id := flagHelper.GetRequiredString("id")

	// 		h := cli.NewHandler(cmd)
	// 		defer h.Close()

	// 		_, err := h.GetAttributeValue(id)
	// 		if err != nil {
	// 			errMsg := "Could not find attribute value"
	// 			cli.ExitWithNotFoundError(errMsg, err)
	// 			cli.ExitWithError(errMsg, err)
	// 		}

	// 		newValue := flagHelper.GetRequiredString("value")

	// 		attr, err := h.UpdateAttributeValue(id, newValue)
	// 		if err != nil {
	// 			errMsg := "Could not update attribute value"
	// 			cli.ExitWithNotFoundError(errMsg, err)
	// 			cli.ExitWithError(errMsg, err)
	// 		}

	// 		v := cli.GetSimpleAttributeValue(attr)
	// 		fmt.Println(cli.SuccessMessage("Attribute value updated"))
	// 		fmt.Println(
	// 			cli.NewTabular().
	// 				Rows([][]string{
	// 					{"Id", v.Id},
	// 					{"FQN", v.FQN},
	// 					{"Members", cli.CommaSeparated(v.Members)},
	// 				}...).Render(),
	// 		)
	// 	},
	// }

	// policy_attributeValuesDeleteCmd = &cobra.Command{
	// 	Use:   "delete",
	// 	Short: "Delete an attribute value",
	// 	Run: func(cmd *cobra.Command, args []string) {
	// 		flagHelper := cli.NewFlagHelper(cmd)
	// 		id := flagHelper.GetRequiredString("id")

	// 		h := cli.NewHandler(cmd)
	// 		defer h.Close()

	// 		value, err := h.GetAttributeValue(id)
	// 		if err != nil {
	// 			errMsg := "Could not find attribute value"
	// 			cli.ExitWithNotFoundError(errMsg, err)
	// 			cli.ExitWithError(errMsg, err)
	// 		}

	// 		cli.ConfirmDelete("attribute value", value.Value)

	// 		err = h.DeleteAttributeValue(id)
	// 		if err != nil {
	// 			errMsg := "Could not delete attribute value"
	// 			cli.ExitWithNotFoundError(errMsg, err)
	// 			cli.ExitWithError(errMsg, err)
	// 		}

	// 		v := cli.GetSimpleAttributeValue(value)
	// 		fmt.Println(cli.SuccessMessage("Attribute deleted"))
	// 		fmt.Println(
	// 			cli.NewTabular().
	// 				Rows([][]string{
	// 					{"Id", v.Id},
	// 					{"FQN", v.FQN},
	// 					{"Members", cli.CommaSeparated(v.Members)},
	// 				}...).Render(),
	// 		)
	// 	},
	// }

	// ///
	// /// Attribute Value Members
	// ///
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
	// 		members := flagHelper.GetStringSlice("members", attrValueMembers, cli.FlagHelperStringSliceOptions{})

	// 		// TODO: Implement
	// 		fmt.Println("Not implemented")
	// 		fmt.Printf("id: %s\n", id)
	// 		fmt.Printf("members: %v\n", members)
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

	// 		// TODO: Implement
	// 		fmt.Println("Not implemented")
	// 		fmt.Printf("id: %s\n", id)
	// 		fmt.Printf("members: %v\n", members)
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

	// 		// TODO: Implement
	// 		fmt.Println("Not implemented")
	// 		fmt.Printf("id: %s\n", id)
	// 		fmt.Printf("members: %v\n", members)
	// 	},
	// }
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

	// policy_attributeValuesCmd.AddCommand(policy_attributeValuesCreateCmd)
	// policy_attributeValuesCreateCmd.Flags().StringP("attribute-id", "a", "", "Attribute id")
	// policy_attributeValuesCreateCmd.Flags().StringP("value", "v", "", "Value")

	// policy_attributeValuesCmd.AddCommand(policy_attributeValuesGetCmd)
	// policy_attributeValuesGetCmd.Flags().StringP("id", "i", "", "Attribute value id")

	// policy_attributeValuesCmd.AddCommand(policy_attributeValuesUpdateCmd)
	// policy_attributeValuesUpdateCmd.Flags().StringP("id", "i", "", "Attribute value id")
	// policy_attributeValuesUpdateCmd.Flags().StringP("value", "v", "", "Value")

	// policy_attributeValuesCmd.AddCommand(policy_attributeValuesDeleteCmd)
	// policy_attributeValuesDeleteCmd.Flags().StringP("id", "i", "", "Attribute value id")

	// // Attribute value members
	// policy_attributeValuesCmd.AddCommand(policy_attributeValueMembersCmd)
	// policy_attributeValueMembersCmd.GroupID = "subcommand"

	// policy_attributeValueMembersCmd.AddCommand(policy_attributeValueMembersAddCmd)
	// policy_attributeValueMembersAddCmd.Flags().StringP("id", "i", "", "Attribute value id")
	// policy_attributeValueMembersAddCmd.Flags().StringSliceVar(&attrValueMembers, "members", []string{}, "Members to add")

	// policy_attributeValueMembersCmd.AddCommand(policy_attributeValueMembersRemoveCmd)
	// policy_attributeValueMembersRemoveCmd.Flags().StringP("id", "i", "", "Attribute value id")
	// policy_attributeValueMembersRemoveCmd.Flags().StringSliceVar(&attrValueMembers, "members", []string{}, "Members to add")

	// policy_attributeValueMembersCmd.AddCommand(policy_attributeValueMembersReplaceCmd)
	// policy_attributeValueMembersReplaceCmd.Flags().StringP("id", "i", "", "Attribute value id")
	// policy_attributeValueMembersReplaceCmd.Flags().StringSliceVar(&attrValueMembers, "members", []string{}, "Members to add")

}
