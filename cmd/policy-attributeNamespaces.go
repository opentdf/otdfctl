package cmd

import (
	"fmt"
	"strconv"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var (
	policy_attributeNamespacesCmd = man.Docs.GetCommand("policy/attributes/namespaces")

	forceUnsafe bool
)

func policy_getAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredString("id")

	ns, err := h.GetNamespace(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}
	rows := [][]string{
		{"Id", ns.Id},
		{"Name", ns.Name},
	}
	if mdRows := getMetadataRows(ns.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, ns.Id, t, ns)
}

func policy_listAttributeNamespaces(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	state := cli.GetState(cmd)
	list, err := h.ListNamespaces(state)
	if err != nil {
		cli.ExitWithError("Failed to list namespaces", err)
	}
	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("name", "Name", 4),
		table.NewFlexColumn("active", "Active", 3),
		table.NewFlexColumn("labels", "Labels", 1),
		table.NewFlexColumn("created_at", "Created At", 1),
		table.NewFlexColumn("updated_at", "Updated At", 1),
	)
	rows := []table.Row{}
	for _, ns := range list {
		metadata := cli.ConstructMetadata(ns.Metadata)
		rows = append(rows,
			table.NewRow(table.RowData{
				"id":         ns.Id,
				"name":       ns.Name,
				"active":     strconv.FormatBool(ns.Active.GetValue()),
				"labels":     metadata["Labels"],
				"created_at": metadata["Created At"],
				"updated_at": metadata["Updated At"],
			}),
		)
	}
	t = t.WithRows(rows)
	HandleSuccess(cmd, "", t, list)
}

func policy_createAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	name := c.Flags.GetRequiredString("name")
	metadataLabels := c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	created, err := h.CreateNamespace(name, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create namespace", err)
	}
	rows := [][]string{
		{"Name", name},
		{"Id", created.Id},
	}
	if mdRows := getMetadataRows(created.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, created.Id, t, created)
}

func policy_deactivateAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	force := c.Flags.GetOptionalBool("force")
	id := c.Flags.GetRequiredString("id")

	ns, err := h.GetNamespace(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	if !force {
		cli.ConfirmAction(cli.ActionDeactivate, "namespace", ns.Name, false)
	}

	d, err := h.DeactivateNamespace(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to deactivate namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}
	rows := [][]string{
		{"Id", ns.Id},
		{"Name", ns.Name},
	}
	if mdRows := getMetadataRows(d.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, ns.Id, t, d)
}

func policy_updateAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredString("id")
	labels := c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	ns, err := h.UpdateNamespace(
		id,
		getMetadataMutable(labels),
		getMetadataUpdateBehavior(),
	)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to update namespace (%s)", id), err)
	}
	rows := [][]string{
		{"Id", ns.Id},
		{"Name", ns.Name},
	}
	if mdRows := getMetadataRows(ns.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, id, t, ns)
}

func policy_unsafeDeleteAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredString("id")

	ns, err := h.GetNamespace(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	if !forceUnsafe {
		cli.ConfirmTextInput(cli.ActionDelete, "namespace", cli.InputNameFQN, ns.GetFqn())
	}

	if err := h.UnsafeDeleteNamespace(id, ns.GetFqn()); err != nil {
		errMsg := fmt.Sprintf("Failed to delete namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	rows := [][]string{
		{"Id", ns.GetId()},
		{"Name", ns.GetName()},
	}
	if mdRows := getMetadataRows(ns.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, ns.Id, t, ns)
}

func policy_unsafeReactivateAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredString("id")

	ns, err := h.GetNamespace(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	if !forceUnsafe {
		cli.ConfirmTextInput(cli.ActionReactivate, "namespace", cli.InputNameFQN, ns.GetFqn())
	}

	ns, err = h.UnsafeReactivateNamespace(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to reactivate namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	rows := [][]string{
		{"Id", ns.GetId()},
		{"Name", ns.GetName()},
	}
	if mdRows := getMetadataRows(ns.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, ns.Id, t, ns)
}

func policy_unsafeUpdateAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredString("id")
	name := c.Flags.GetRequiredString("name")

	ns, err := h.GetNamespace(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	if !forceUnsafe {
		cli.ConfirmTextInput(cli.ActionUpdateUnsafe, "namespace", cli.InputNameFQNUpdated, ns.GetFqn())
	}

	ns, err = h.UnsafeUpdateNamespace(id, name)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to reactivate namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	rows := [][]string{
		{"Id", ns.GetId()},
		{"Name", ns.GetName()},
	}
	if mdRows := getMetadataRows(ns.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, ns.Id, t, ns)
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
	deactivateCmd.Flags().Bool(
		deactivateCmd.GetDocFlag("force").Name,
		false,
		deactivateCmd.GetDocFlag("force").Description,
	)

	// unsafe
	unsafeCmd := man.Docs.GetCommand("policy/attributes/namespaces/unsafe")
	unsafeCmd.PersistentFlags().BoolVar(
		&forceUnsafe,
		unsafeCmd.GetDocFlag("force").Name,
		false,
		unsafeCmd.GetDocFlag("force").Description,
	)
	deleteCmd := man.Docs.GetCommand("policy/attributes/namespaces/unsafe/delete",
		man.WithRun(policy_unsafeDeleteAttributeNamespace),
	)
	deleteCmd.Flags().StringP(
		deactivateCmd.GetDocFlag("id").Name,
		deactivateCmd.GetDocFlag("id").Shorthand,
		deactivateCmd.GetDocFlag("id").Default,
		deactivateCmd.GetDocFlag("id").Description,
	)
	reactivateCmd := man.Docs.GetCommand("policy/attributes/namespaces/unsafe/reactivate",
		man.WithRun(policy_unsafeReactivateAttributeNamespace),
	)
	reactivateCmd.Flags().StringP(
		deactivateCmd.GetDocFlag("id").Name,
		deactivateCmd.GetDocFlag("id").Shorthand,
		deactivateCmd.GetDocFlag("id").Default,
		deactivateCmd.GetDocFlag("id").Description,
	)
	unsafeUpdateCmd := man.Docs.GetCommand("policy/attributes/namespaces/unsafe/update",
		man.WithRun(policy_unsafeUpdateAttributeNamespace),
	)
	unsafeUpdateCmd.Flags().StringP(
		deactivateCmd.GetDocFlag("id").Name,
		deactivateCmd.GetDocFlag("id").Shorthand,
		deactivateCmd.GetDocFlag("id").Default,
		deactivateCmd.GetDocFlag("id").Description,
	)
	unsafeUpdateCmd.Flags().StringP(
		unsafeUpdateCmd.GetDocFlag("name").Name,
		unsafeUpdateCmd.GetDocFlag("name").Shorthand,
		unsafeUpdateCmd.GetDocFlag("name").Default,
		unsafeUpdateCmd.GetDocFlag("name").Description,
	)
	unsafeCmd.AddSubcommands(deleteCmd, reactivateCmd, unsafeUpdateCmd)

	policy_attributeNamespacesCmd.AddSubcommands(getCmd, listCmd, createDoc, updateCmd, deactivateCmd, unsafeCmd)
	policy_attributesCmd.AddCommand(&policy_attributeNamespacesCmd.Command)
}
