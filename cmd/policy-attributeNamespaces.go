package cmd

import (
	"fmt"
	"strconv"

	"github.com/charmbracelet/lipgloss"
	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var (
	policy_attributeNamespacesCmd = man.Docs.GetCommand("policy/attributes/namespaces")
	policy_NamespaceKeysCmd       = man.Docs.GetCommand("policy/attributes/namespaces/keys")

	forceUnsafe bool
)

func policy_getAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	ns, err := h.GetNamespace(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get namespace (%s)", id)
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
	HandleSuccess(cmd, ns.GetId(), t, ns)
}

func policy_listAttributeNamespaces(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	state := cli.GetState(cmd)
	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")

	list, page, err := h.ListNamespaces(state, limit, offset)
	if err != nil {
		cli.ExitWithError("Failed to list namespaces", err)
	}
	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("name", "Name", cli.FlexColumnWidthFour),
		table.NewFlexColumn("active", "Active", cli.FlexColumnWidthThree),
		table.NewFlexColumn("labels", "Labels", cli.FlexColumnWidthOne),
		table.NewFlexColumn("created_at", "Created At", cli.FlexColumnWidthOne),
		table.NewFlexColumn("updated_at", "Updated At", cli.FlexColumnWidthOne),
	)
	rows := []table.Row{}
	for _, ns := range list {
		metadata := cli.ConstructMetadata(ns.GetMetadata())
		rows = append(rows,
			table.NewRow(table.RowData{
				"id":         ns.GetId(),
				"name":       ns.GetName(),
				"active":     strconv.FormatBool(ns.GetActive().GetValue()),
				"labels":     metadata["Labels"],
				"created_at": metadata["Created At"],
				"updated_at": metadata["Updated At"],
			}),
		)
	}
	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, list)
}

func policy_createAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	name := c.Flags.GetRequiredString("name")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	created, err := h.CreateNamespace(name, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create namespace", err)
	}
	rows := [][]string{
		{"Name", name},
		{"Id", created.GetId()},
	}
	if mdRows := getMetadataRows(created.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, created.GetId(), t, created)
}

func policy_deactivateAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	force := c.Flags.GetOptionalBool("force")
	id := c.Flags.GetRequiredID("id")

	ns, err := h.GetNamespace(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to find namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	cli.ConfirmAction(cli.ActionDeactivate, "namespace", ns.GetName(), force)

	d, err := h.DeactivateNamespace(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to deactivate namespace (%s)", id)
		cli.ExitWithError(errMsg, err)
	}
	rows := [][]string{
		{"Id", ns.GetId()},
		{"Name", ns.GetName()},
	}
	if mdRows := getMetadataRows(d.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, ns.GetId(), t, d)
}

func policy_updateAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	ns, err := h.UpdateNamespace(
		id,
		getMetadataMutable(metadataLabels),
		getMetadataUpdateBehavior(),
	)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to update namespace (%s)", id), err)
	}
	rows := [][]string{
		{"Id", ns.GetId()},
		{"Name", ns.GetName()},
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

	id := c.Flags.GetRequiredID("id")

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
	HandleSuccess(cmd, ns.GetId(), t, ns)
}

func policy_unsafeReactivateAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

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
	HandleSuccess(cmd, ns.GetId(), t, ns)
}

func policy_unsafeUpdateAttributeNamespace(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
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
	HandleSuccess(cmd, ns.GetId(), t, ns)
}

func policy_NamespaceKeysAdd(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	ns := c.Flags.GetRequiredString("namespace")
	pkID := c.Flags.GetRequiredID("public-key-id")

	_, err := h.AddPublicKeyToNamespace(c.Context(), ns, pkID)
	if err != nil {
		cli.ExitWithError("Failed to add public key to namespace", err)
	}

	rows := [][]string{
		{"Public Key Id", pkID},
		{"Namespace", ns},
	}

	t := cli.NewTabular(rows...)

	HandleSuccess(cmd, "Public key added to namespace", t, nil)
}

func policy_NamespaceKeysRemove(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	ns := c.Flags.GetRequiredString("namespace")
	pkID := c.Flags.GetRequiredID("public-key-id")

	_, err := h.RemovePublicKeyFromNamespace(c.Context(), ns, pkID)
	if err != nil {
		cli.ExitWithError("Failed to remove public key from namespace", err)
	}

	rows := [][]string{
		{"Public Key Id", pkID},
		{"Namespace", ns},
	}

	t := cli.NewTabular(rows...)

	HandleSuccess(cmd, "Public key removed from namespace", t, nil)
}

func policy_NamespaceKeysListcmd(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	ns := c.Flags.GetRequiredString("namespace")
	showPublicKey := c.Flags.GetOptionalBool("show-public-key")

	list, err := h.GetNamespace(ns)
	if err != nil {
		cli.ExitWithError("Failed to list namespace keys", err)
	}

	columns := []table.Column{
		table.NewFlexColumn("kas_name", "KAS Name", cli.FlexColumnWidthThree),
		table.NewFlexColumn("kas_uri", "KAS URI", cli.FlexColumnWidthThree),
		table.NewFlexColumn("kid", "Key ID", cli.FlexColumnWidthThree),
		table.NewFlexColumn("alg", "Algorithm", cli.FlexColumnWidthThree),
	}

	if showPublicKey {
		columns = append(columns, table.NewFlexColumn("public_key", "Public Key", cli.FlexColumnWidthFour))
	}

	t := cli.NewTable(columns...)
	rows := []table.Row{}
	for _, key := range list.GetKeys() {

		alg, err := enumToAlg(key.GetPublicKey().GetAlg())
		if err != nil {
			cli.ExitWithError("Failed to get algorithm", err)
		}

		rowStyle := lipgloss.NewStyle().BorderBottom(true).BorderStyle(lipgloss.NormalBorder())

		if key.GetIsActive().GetValue() {
			rowStyle = rowStyle.Background(cli.ColorGreen.Background)
		} else {
			rowStyle = rowStyle.Background(cli.ColorRed.Background)
		}

		rd := table.RowData{
			"key_id":     key.GetPublicKey().GetKid(),
			"algorithm":  alg,
			"is_active":  key.GetIsActive().GetValue(),
			"kas_id":     key.GetKas().GetId(),
			"kas_name":   key.GetKas().GetName(),
			"kas_uri":    key.GetKas().GetUri(),
			"public_key": key.GetPublicKey().GetPem(),
		}
		rows = append(rows, table.NewRow(rd).WithStyle(rowStyle))
	}
	t = t.WithRows(rows)
	HandleSuccess(cmd, "", t, list)
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
	injectListPaginationFlags(listCmd)

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

	namespaceKeysAddDoc := man.Docs.GetCommand("policy/attributes/namespaces/keys/add",
		man.WithRun(policy_NamespaceKeysAdd),
	)
	namespaceKeysAddDoc.Flags().StringP(
		namespaceKeysAddDoc.GetDocFlag("namespace").Name,
		namespaceKeysAddDoc.GetDocFlag("namespace").Shorthand,
		namespaceKeysAddDoc.GetDocFlag("namespace").Default,
		namespaceKeysAddDoc.GetDocFlag("namespace").Description,
	)
	namespaceKeysAddDoc.Flags().StringP(
		namespaceKeysAddDoc.GetDocFlag("public-key-id").Name,
		namespaceKeysAddDoc.GetDocFlag("public-key-id").Shorthand,
		namespaceKeysAddDoc.GetDocFlag("public-key-id").Default,
		namespaceKeysAddDoc.GetDocFlag("public-key-id").Description,
	)

	namespaceKeysRemoveDoc := man.Docs.GetCommand("policy/attributes/namespaces/keys/remove",
		man.WithRun(policy_NamespaceKeysRemove),
	)
	namespaceKeysRemoveDoc.Flags().StringP(
		namespaceKeysRemoveDoc.GetDocFlag("namespace").Name,
		namespaceKeysRemoveDoc.GetDocFlag("namespace").Shorthand,
		namespaceKeysRemoveDoc.GetDocFlag("namespace").Default,
		namespaceKeysRemoveDoc.GetDocFlag("namespace").Description,
	)
	namespaceKeysRemoveDoc.Flags().StringP(
		namespaceKeysRemoveDoc.GetDocFlag("public-key-id").Name,
		namespaceKeysRemoveDoc.GetDocFlag("public-key-id").Shorthand,
		namespaceKeysRemoveDoc.GetDocFlag("public-key-id").Default,
		namespaceKeysRemoveDoc.GetDocFlag("public-key-id").Description,
	)

	namespaceKeysListDoc := man.Docs.GetCommand("policy/attributes/namespaces/keys/list",
		man.WithRun(policy_NamespaceKeysListcmd),
	)
	namespaceKeysListDoc.Flags().StringP(
		namespaceKeysListDoc.GetDocFlag("namespace").Name,
		namespaceKeysListDoc.GetDocFlag("namespace").Shorthand,
		namespaceKeysListDoc.GetDocFlag("namespace").Default,
		namespaceKeysListDoc.GetDocFlag("namespace").Description,
	)
	namespaceKeysListDoc.Flags().BoolP(
		namespaceKeysListDoc.GetDocFlag("show-public-key").Name,
		namespaceKeysListDoc.GetDocFlag("show-public-key").Shorthand,
		namespaceKeysListDoc.GetDocFlag("show-public-key").DefaultAsBool(),
		namespaceKeysListDoc.GetDocFlag("show-public-key").Description,
	)

	policy_NamespaceKeysCmd.AddSubcommands(namespaceKeysAddDoc, namespaceKeysRemoveDoc, namespaceKeysListDoc)
	policy_attributeNamespacesCmd.AddSubcommands(getCmd, listCmd, createDoc, updateCmd, deactivateCmd, unsafeCmd, policy_NamespaceKeysCmd)
	policy_attributesCmd.AddCommand(&policy_attributeNamespacesCmd.Command)
}
