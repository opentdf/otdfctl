package cmd

import (
	"fmt"
	"strings"

	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/tructl/issues/73] is addressed

var (
	attrValues                 []string
	metadataLabels             []string
	forceReplaceMetadataLabels bool

	policy_attributeCommands = []string{
		policy_attributesCreateCmd.Use,
		policy_attributeGetCmd.Use,
		policy_attributesListCmd.Use,
		policy_attributeUpdateCmd.Use,
		policy_attributesDeactivateCmd.Use,
	}

	policy_attributesCmd = &cobra.Command{
		Use:   "attributes",
		Short: "Manage attributes [" + strings.Join(policy_attributeCommands, ", ") + "]",
		Long: `
Attributes - commands to manage attributes within the platform.

Attributes are used to to define the properties of a piece of data. These attributes will then be
used to define the access controls based on subject encodings and entity entitlements.
`,
	}

	// Create an attribute
	policy_attributesCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "Create an attribute",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			name := flagHelper.GetRequiredString("name")
			rule := flagHelper.GetRequiredString("rule")
			values := flagHelper.GetStringSlice("value", attrValues, cli.FlagHelperStringSliceOptions{})
			namespace := flagHelper.GetRequiredString("namespace")
			metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

			attr, err := h.CreateAttribute(name, rule, namespace, values, getMetadataMutable(metadataLabels))
			if err != nil {
				cli.ExitWithError("Failed to create attribute", err)
			}

			a := cli.GetSimpleAttribute(&policy.Attribute{
				Id:        attr.Id,
				Name:      attr.Name,
				Rule:      attr.Rule,
				Values:    attr.Values,
				Namespace: attr.Namespace,
			})

			t := cli.NewTabular().Rows([][]string{
				{"Name", a.Name},
				{"Rule", a.Rule},
				{"Values", cli.CommaSeparated(a.Values)},
				{"Namespace", a.Namespace},
			}...)

			HandleSuccess(cmd, a.Id, t, attr)
		},
	}

	// Get an attribute
	policy_attributeGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get an attribute",
		Run: func(cmd *cobra.Command, args []string) {
			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			h := cli.NewHandler(cmd)
			defer h.Close()

			attr, err := h.GetAttribute(id)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to get attribute (%s)", id)
				cli.ExitWithError(errMsg, err)
			}

			a := cli.GetSimpleAttribute(attr)
			t := cli.NewTabular().
				Rows([][]string{
					{"Id", a.Id},
					{"Name", a.Name},
					{"Rule", a.Rule},
					{"Values", cli.CommaSeparated(a.Values)},
					{"Namespace", a.Namespace},
				}...)
			HandleSuccess(cmd, a.Id, t, attr)
		},
	}

	// List attributes
	policy_attributesListCmd = &cobra.Command{
		Use:   "list",
		Short: "List attributes",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			attrs, err := h.ListAttributes()
			if err != nil {
				cli.ExitWithError("Failed to list attributes", err)
			}

			t := cli.NewTable()
			t.Headers("Id", "Namespace", "Name", "Rule", "Values")
			for _, attr := range attrs {
				a := cli.GetSimpleAttribute(attr)
				t.Row(
					a.Id,
					a.Namespace,
					a.Name,
					a.Rule,
					cli.CommaSeparated(a.Values),
				)
			}
			HandleSuccess(cmd, "", t, attrs)
		},
	}

	policy_attributesDeactivateCmd = &cobra.Command{
		Use:   "deactivate",
		Short: "Deactivate an attribute",
		Run: func(cmd *cobra.Command, args []string) {
			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")

			h := cli.NewHandler(cmd)
			defer h.Close()

			attr, err := h.GetAttribute(id)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to get attribute (%s)", id)
				cli.ExitWithError(errMsg, err)
			}

			cli.ConfirmAction(cli.ActionDeactivate, "attribute", attr.Name)

			attr, err = h.DeactivateAttribute(id)
			if err != nil {
				errMsg := fmt.Sprintf("Failed to deactivate attribute (%s)", id)
				cli.ExitWithError(errMsg, err)
			}

			a := cli.GetSimpleAttribute(attr)
			t := cli.NewTabular().
				Rows([][]string{
					{"Name", a.Name},
					{"Rule", a.Rule},
					{"Values", cli.CommaSeparated(a.Values)},
					{"Namespace", a.Namespace},
				}...)
			HandleSuccess(cmd, a.Id, t, a)
		},
	}

	// Update one attribute
	policy_attributeUpdateCmd = &cobra.Command{
		Use:   "update",
		Short: "Update an attribute",
		Run: func(cmd *cobra.Command, args []string) {
			h := cli.NewHandler(cmd)
			defer h.Close()

			flagHelper := cli.NewFlagHelper(cmd)
			id := flagHelper.GetRequiredString("id")
			labels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

			if a, err := h.UpdateAttribute(id, getMetadataMutable(labels), getMetadataUpdateBehavior()); err != nil {
				cli.ExitWithError(fmt.Sprintf("Failed to update attribute (%s)", id), err)
			} else {
				HandleSuccess(cmd, id, nil, a)
			}
		},
	}
)

func init() {
	policyCmd.AddCommand(policy_attributesCmd)

	// Create an attribute
	policy_attributesCmd.AddCommand(policy_attributesCreateCmd)
	policy_attributesCreateCmd.Flags().StringP("name", "n", "", "Name of the attribute")
	policy_attributesCreateCmd.Flags().StringP("rule", "r", "", "Rule of the attribute")
	policy_attributesCreateCmd.Flags().StringSliceVarP(&attrValues, "value", "v", []string{}, "Values of the attribute")
	policy_attributesCreateCmd.Flags().StringP("namespace", "s", "", "Namespace of the attribute")
	injectLabelFlags(policy_attributesCreateCmd, false)

	// Get an attribute
	policy_attributesCmd.AddCommand(policy_attributeGetCmd)
	policy_attributeGetCmd.Flags().StringP("id", "i", "", "Id of the attribute")

	// List attributes
	policy_attributesCmd.AddCommand(policy_attributesListCmd)

	// Update an attribute
	policy_attributesCmd.AddCommand(policy_attributeUpdateCmd)
	policy_attributeUpdateCmd.Flags().StringP("id", "i", "", "Id of the attribute")
	injectLabelFlags(policy_attributeUpdateCmd, true)

	// Deactivate an attribute
	policy_attributesCmd.AddCommand(policy_attributesDeactivateCmd)
	policy_attributesDeactivateCmd.Flags().StringP("id", "i", "", "Id of the attribute")
}
