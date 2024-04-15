package cmd

import (
	"fmt"
	"strconv"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

// TODO: add metadata to outputs once [https://github.com/opentdf/otdfctl/issues/73] is addressed

var (
	policy_attributeNamespacesCmd = man.Docs.GetCommand("policy/attributes/namespaces")
)

func policy_getAttributeNamespace(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	ns, err := h.GetNamespace(id)
	metadata := cli.ConstructMetadata(ns.Metadata)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	t := cli.NewTabular().
		Rows([][]string{
			{"Id", ns.Id},
			{"Name", ns.Name},
			{"Metadata.Labels", metadata["Labels"]},
			{"Metadata.CreatedAt", metadata["Created At"]},
			{"Metadata.UpdatedAt", metadata["Updated At"]},
		}...)
	HandleSuccess(cmd, ns.Id, t, ns)
}

func policy_listAttributeNamespaces(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	state := cli.GetState(cmd)
	list, err := h.ListNamespaces(state)
	if err != nil {
		cli.ExitWithError("Failed to list namespaces", err)
	}

	t := cli.NewTable()
	t.Headers("Id", "Name", "Active")
	for _, ns := range list {
		t.Row(
			ns.Id,
			ns.Name,
			strconv.FormatBool(ns.Active.GetValue()),
		)
	}
	HandleSuccess(cmd, "", t, list)
}

func policy_createAttributeNamespace(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	name := flagHelper.GetRequiredString("name")
	metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

	created, err := h.CreateNamespace(name, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create namespace", err)
	}
	metadata := cli.ConstructMetadata(created.Metadata)

	t := cli.NewTabular().Rows([][]string{
		{"Name", name},
		{"Id", created.Id},
		{"Metadata.Labels", metadata["Labels"]},
		{"Metadata.CreatedAt", metadata["Created At"]},
		{"Metadata.UpdatedAt", metadata["Updated At"]},
	}...)
	HandleSuccess(cmd, created.Id, t, created)
}

func policy_deactivateAttributeNamespace(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	ns, err := h.GetNamespace(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	cli.ConfirmAction(cli.ActionDeactivate, "namespace", ns.Name)

	d, err := h.DeactivateNamespace(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to deactivate namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}
	metadata := cli.ConstructMetadata(d.Metadata)

	t := cli.NewTabular().
		Rows([][]string{
			{"Id", ns.Id},
			{"Name", ns.Name},
			{"Metadata.Labels", metadata["Labels"]},
			{"Metadata.CreatedAt", metadata["Created At"]},
			{"Metadata.UpdatedAt", metadata["Updated At"]},
		}...)
	HandleSuccess(cmd, ns.Id, t, d)
}

func policy_updateAttributeNamespace(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")
	labels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

	ns, err := h.UpdateNamespace(
		id,
		getMetadataMutable(labels),
		getMetadataUpdateBehavior(),
	)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to update namespace (%s)", id), err)
	}
	metadata := cli.ConstructMetadata(ns.Metadata)

	t := cli.NewTabular().Rows([][]string{
		{"Id", ns.Id},
		{"Name", ns.Name},
		{"Metadata.Labels", metadata["Labels"]},
		{"Metadata.CreatedAt", metadata["Created At"]},
		{"Metadata.UpdatedAt", metadata["Updated At"]},
	}...)
	HandleSuccess(cmd, id, t, ns)
}

func init() {
	getCmd := man.Docs.GetCommand("policy/attributes/namespaces/get",
		man.WithRun(policy_getAttributeNamespace),
	)
	getCmd.Flags().StringP(
		getCmd.GetDocFlag("id").Name,
		getCmd.GetDocFlag("id").Shorthand,
		getCmd.GetDocFlag("id").Default,
		getCmd.GetDocFlag("id").Description,
	)

	listCmd := man.Docs.GetCommand("policy/attributes/namespaces/list",
		man.WithRun(policy_listAttributeNamespaces),
	)
	listCmd.Flags().StringP(
		listCmd.GetDocFlag("state").Name,
		listCmd.GetDocFlag("state").Shorthand,
		listCmd.GetDocFlag("state").Default,
		listCmd.GetDocFlag("state").Description,
	)

	createDoc := man.Docs.GetCommand("policy/attributes/namespaces/create",
		man.WithRun(policy_createAttributeNamespace),
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("name").Name,
		createDoc.GetDocFlag("name").Shorthand,
		createDoc.GetDocFlag("name").Default,
		createDoc.GetDocFlag("name").Description,
	)
	injectLabelFlags(&createDoc.Command, false)

	updateCmd := man.Docs.GetCommand("policy/attributes/namespaces/update",
		man.WithRun(policy_updateAttributeNamespace),
	)
	updateCmd.Flags().StringP(
		updateCmd.GetDocFlag("id").Name,
		updateCmd.GetDocFlag("id").Shorthand,
		updateCmd.GetDocFlag("id").Default,
		updateCmd.GetDocFlag("id").Description,
	)
	injectLabelFlags(&updateCmd.Command, true)

	deactivateCmd := man.Docs.GetCommand("policy/attributes/namespaces/deactivate",
		man.WithRun(policy_deactivateAttributeNamespace),
	)
	deactivateCmd.Flags().StringP(
		deactivateCmd.GetDocFlag("id").Name,
		deactivateCmd.GetDocFlag("id").Shorthand,
		deactivateCmd.GetDocFlag("id").Default,
		deactivateCmd.GetDocFlag("id").Description,
	)

	policy_attributeNamespacesCmd.AddSubcommands(getCmd, listCmd, createDoc, updateCmd, deactivateCmd)
	policy_attributesCmd.AddCommand(&policy_attributeNamespacesCmd.Command)
}
