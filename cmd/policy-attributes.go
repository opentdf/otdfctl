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

	policy_attributesCmd = man.Docs.GetCommand("policy/attributes")
)

func policy_createAttribute(cmd *cobra.Command, args []string) {
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
		{"Metadata.Labels", a.Metadata["Labels"]},
		{"Metadata.CreatedAt", a.Metadata["Created At"]},
		{"Metadata.UpdatedAt", a.Metadata["Updated At"]},
	}...)

	HandleSuccess(cmd, a.Id, t, attr)
}

func policy_getAttribute(cmd *cobra.Command, args []string) {
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
			{"Metadata.Labels", a.Metadata["Labels"]},
			{"Metadata.CreatedAt", a.Metadata["Created At"]},
			{"Metadata.UpdatedAt", a.Metadata["Updated At"]},
		}...)
	HandleSuccess(cmd, a.Id, t, attr)
}

func policy_listAttributes(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	state := cli.GetState(cmd)
	attrs, err := h.ListAttributes(state)
	if err != nil {
		cli.ExitWithError("Failed to list attributes", err)
	}

	t := cli.NewTable()
	t.Headers("Id", "Namespace", "Name", "Rule", "Values", "Active", "Metadata.Labels", "Metadata.CreatedAt", "Metadata.UpdatedAt")
	for _, attr := range attrs {
		a := cli.GetSimpleAttribute(attr)
		t.Row(
			a.Id,
			a.Namespace,
			a.Name,
			a.Rule,
			cli.CommaSeparated(a.Values),
			a.Active,
			a.Metadata["Labels"],
			a.Metadata["Created At"],
			a.Metadata["Updated At"],
		)
	}
	HandleSuccess(cmd, "", t, attrs)
}

func policy_deactivateAttribute(cmd *cobra.Command, args []string) {
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
			{"Metadata.Labels", a.Metadata["Labels"]},
			{"Metadata.CreatedAt", a.Metadata["Created At"]},
			{"Metadata.UpdatedAt", a.Metadata["Updated At"]},
		}...)
	HandleSuccess(cmd, a.Id, t, a)
}

func policy_updateAttribute(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")
	labels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

	if a, err := h.UpdateAttribute(id, getMetadataMutable(labels), getMetadataUpdateBehavior()); err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to update attribute (%s)", id), err)
	} else {
		metadata := cli.ConstructMetadata(a.Metadata)
		t := cli.NewTabular().
			Rows([][]string{
				{"Id", a.Id},
				{"Name", a.Name},
				{"Metadata.Labels", metadata["Labels"]},
				{"Metadata.CreatedAt", metadata["Created At"]},
				{"Metadata.UpdatedAt", metadata["Updated At"]},
			}...)
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
	createDoc.Flags().StringSliceVarP(
		&attrValues,
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
	//"state", "s", "active", "Filter by state [active, inactive, any]"
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

	policy_attributesCmd.AddSubcommands(createDoc, getDoc, listDoc, updateDoc, deactivateDoc)
	policyCmd.AddCommand(&policy_attributesCmd.Command)
}
