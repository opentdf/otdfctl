package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/otdfctl/issues/73] is addressed

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
		Use:   man.Docs.GetDoc("policy/attributes").Use,
		Short: man.Docs.GetDoc("policy/attributes").GetShort(policy_attributeCommands),
		Long:  man.Docs.GetDoc("policy/attributes").Long,
	}

	// Create an attribute
	policy_attributesCreateCmd = &cobra.Command{
		Use:   man.Docs.GetDoc("policy/attributes/create").Use,
		Short: man.Docs.GetDoc("policy/attributes/create").Short,
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
		Use:   man.Docs.GetDoc("policy/attributes/get").Use,
		Short: man.Docs.GetDoc("policy/attributes/get").Short,
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

			state := cli.GetState(cmd)
			attrs, err := h.ListAttributes(state)
			if err != nil {
				cli.ExitWithError("Failed to list attributes", err)
			}

			t := cli.NewTable()
			t.Headers("Id", "Namespace", "Name", "Rule", "Values", "Active")
			for _, attr := range attrs {
				a := cli.GetSimpleAttribute(attr)
				t.Row(
					a.Id,
					a.Namespace,
					a.Name,
					a.Rule,
					cli.CommaSeparated(a.Values),
					a.Active,
				)
			}
			HandleSuccess(cmd, "", t, attrs)
		},
	}

	policy_attributesDeactivateCmd = &cobra.Command{
		Use:   man.Docs.GetDoc("policy/attributes/deactivate").Use,
		Short: man.Docs.GetDoc("policy/attributes/deactivate").Short,
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
		Use:   man.Docs.GetDoc("policy/attributes/update").Use,
		Short: man.Docs.GetDoc("policy/attributes/update").Short,
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
	createDoc := man.Docs.GetDoc("policy/attributes/create")
	policy_attributesCmd.AddCommand(policy_attributesCreateCmd)
	policy_attributesCreateCmd.Flags().StringP(
		createDoc.GetDocFlag("name").Name,
		createDoc.GetDocFlag("name").Shorthand,
		createDoc.GetDocFlag("name").Default,
		createDoc.GetDocFlag("name").Description,
	)
	policy_attributesCreateCmd.Flags().StringP(
		createDoc.GetDocFlag("rule").Name,
		createDoc.GetDocFlag("rule").Shorthand,
		createDoc.GetDocFlag("rule").Default,
		createDoc.GetDocFlag("rule").Description,
	)
	policy_attributesCreateCmd.Flags().StringSliceVarP(
		&attrValues,
		createDoc.GetDocFlag("value").Name,
		createDoc.GetDocFlag("value").Shorthand,
		[]string{},
		createDoc.GetDocFlag("value").Description,
	)
	policy_attributesCreateCmd.Flags().StringP(
		createDoc.GetDocFlag("namespace").Name,
		createDoc.GetDocFlag("namespace").Shorthand,
		createDoc.GetDocFlag("namespace").Default,
		createDoc.GetDocFlag("namespace").Description,
	)
	injectLabelFlags(policy_attributesCreateCmd, false)

	// Get an attribute
	getDoc := man.Docs.GetDoc("policy/attributes/get")
	policy_attributesCmd.AddCommand(policy_attributeGetCmd)
	policy_attributeGetCmd.Flags().StringP(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Shorthand,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)

	// List attributes
	policy_attributesCmd.AddCommand(policy_attributesListCmd)
	policy_attributesListCmd.Flags().StringP("state", "s", "active", "Filter by state [active, inactive, any]")

	// Update an attribute
	updateDoc := man.Docs.GetDoc("policy/attributes/update")
	policy_attributesCmd.AddCommand(policy_attributeUpdateCmd)
	policy_attributeUpdateCmd.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	injectLabelFlags(policy_attributeUpdateCmd, true)

	// Deactivate an attribute
	deactivateDoc := man.Docs.GetDoc("policy/attributes/deactivate")
	policy_attributesCmd.AddCommand(policy_attributesDeactivateCmd)
	policy_attributesDeactivateCmd.Flags().StringP(
		deactivateDoc.GetDocFlag("id").Name,
		deactivateDoc.GetDocFlag("id").Shorthand,
		deactivateDoc.GetDocFlag("id").Default,
		deactivateDoc.GetDocFlag("id").Description,
	)
}
