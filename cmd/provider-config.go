package cmd

import (
	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

func key_createProviderConfig(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	name := c.Flags.GetRequiredString("name")
	config := c.Flags.GetRequiredString("config")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	if !isJSON(config) {
		cli.ExitWithError("Invalid JSON format for config ", nil)
	}

	// Do not need to get provider config after, since this endpoint returns the created config.
	pc, err := h.CreateProviderConfig(c.Context(), name, []byte(config), getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create provider config", err)
	}

	rows := [][]string{
		{"ID", pc.GetId()},
		{"Name", pc.GetName()},
		{"Config", string(pc.GetConfigJson())},
	}

	if mdRows := getMetadataRows(pc.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)

	HandleSuccess(cmd, pc.GetId(), t, pc)
}

func key_getProviderConfig(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetOptionalID("id")
	name := c.Flags.GetOptionalString("name")

	pc, err := h.GetProviderConfig(c.Context(), id, name)
	if err != nil {
		cli.ExitWithError("Failed to get provider config", err)
	}

	rows := [][]string{
		{"ID", pc.GetId()},
		{"Name", pc.GetName()},
		{"Config", string(pc.GetConfigJson())},
	}

	if mdRows := getMetadataRows(pc.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)

	HandleSuccess(cmd, pc.GetId(), t, pc)
}

func key_updateProviderConfig(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	name := c.Flags.GetOptionalString("name")
	config := c.Flags.GetOptionalString("config")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	if name == "" && config == "" && len(metadataLabels) == 0 {
		cli.ExitWithError("At least one field (name, config, or metadata labels) must be updated", nil)
	}

	if config != "" && !isJSON(config) {
		cli.ExitWithError("Cannot update provider config with invalid json", nil)
	}

	pc, err := h.UpdateProviderConfig(c.Context(), id, name, []byte(config), getMetadataMutable(metadataLabels), getMetadataUpdateBehavior())
	if err != nil {
		cli.ExitWithError("Failed to update provider config", err)
	}

	// Get updated provider config.
	pc, err = h.GetProviderConfig(c.Context(), id, "")
	if err != nil {
		cli.ExitWithError("Failed to get provider config", err)
	}

	rows := [][]string{
		{"ID", pc.GetId()},
		{"Name", pc.GetName()},
		{"Config", string(pc.GetConfigJson())},
	}

	if mdRows := getMetadataRows(pc.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)

	HandleSuccess(cmd, pc.GetId(), t, pc)
}

func key_listProviderConfigs(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")

	// Get all provider configs
	providerConfigs, page, err := h.ListProviderConfigs(c.Context(), limit, offset)
	if err != nil {
		cli.ExitWithError("Failed to list provider configs", err)
	}

	t := cli.NewTable(
		// columns should be id, name, config, labels, created_at, updated_at
		table.NewFlexColumn("id", "Provider Config ID", cli.FlexColumnWidthThree),
		table.NewFlexColumn("name", "Provider Config Name", cli.FlexColumnWidthThree),
		table.NewFlexColumn("config", "Provider Config", cli.FlexColumnWidthOne),
		table.NewFlexColumn("labels", "Labels", cli.FlexColumnWidthOne),
		table.NewFlexColumn("created_at", "Created At", cli.FlexColumnWidthOne),
		table.NewFlexColumn("updated_at", "Updated At", cli.FlexColumnWidthOne),
	)
	rows := []table.Row{}
	for _, pc := range providerConfigs {
		metadata := cli.ConstructMetadata(pc.GetMetadata())
		rows = append(rows, table.NewRow(table.RowData{
			"id":         pc.GetId(),
			"name":       pc.GetName(),
			"config":     string(pc.GetConfigJson()),
			"labels":     metadata["Labels"],
			"created_at": metadata["Created At"],
			"updated_at": metadata["Updated At"],
		}))
	}
	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, providerConfigs)
}

func key_deleteProviderConfig(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	err := h.DeleteProviderConfig(c.Context(), id)
	if err != nil {
		cli.ExitWithError("Failed to delete provider config", err)
	}

	rows := [][]string{
		{"Deleted", "true"},
		{"Id", id},
	}

	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, id, t, nil)
}

func init() {
	// Create Provider Config
	createDoc := man.Docs.GetCommand("key/provider-config/create",
		man.WithRun(key_createProviderConfig),
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("name").Name,
		createDoc.GetDocFlag("name").Shorthand,
		createDoc.GetDocFlag("name").Default,
		createDoc.GetDocFlag("name").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("config").Name,
		createDoc.GetDocFlag("config").Shorthand,
		createDoc.GetDocFlag("config").Default,
		createDoc.GetDocFlag("config").Description,
	)
	injectLabelFlags(&createDoc.Command, false)

	// Get Provider Config
	getDoc := man.Docs.GetCommand("key/provider-config/get",
		man.WithRun(key_getProviderConfig),
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Shorthand,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("name").Name,
		getDoc.GetDocFlag("name").Shorthand,
		getDoc.GetDocFlag("name").Default,
		getDoc.GetDocFlag("name").Description,
	)
	getDoc.MarkFlagsOneRequired("id", "name")
	getDoc.MarkFlagsMutuallyExclusive("id", "name")

	// Update Provider Config
	updateDoc := man.Docs.GetCommand("key/provider-config/update",
		man.WithRun(key_updateProviderConfig),
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("name").Name,
		updateDoc.GetDocFlag("name").Shorthand,
		updateDoc.GetDocFlag("name").Default,
		updateDoc.GetDocFlag("name").Description,
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("config").Name,
		updateDoc.GetDocFlag("config").Shorthand,
		updateDoc.GetDocFlag("config").Default,
		updateDoc.GetDocFlag("config").Description,
	)
	injectLabelFlags(&updateDoc.Command, true)

	// List Provider Configs
	listDoc := man.Docs.GetCommand("key/provider-config/list",
		man.WithRun(key_listProviderConfigs),
	)
	injectListPaginationFlags(listDoc)

	// Add Delete Provider Config
	deleteDoc := man.Docs.GetCommand("key/provider-config/delete",
		man.WithRun(key_deleteProviderConfig),
	)
	deleteDoc.Flags().StringP(
		deleteDoc.GetDocFlag("id").Name,
		deleteDoc.GetDocFlag("id").Shorthand,
		deleteDoc.GetDocFlag("id").Default,
		deleteDoc.GetDocFlag("id").Description,
	)

	doc := man.Docs.GetCommand("key/provider-config",
		man.WithSubcommands(createDoc, getDoc, updateDoc, listDoc, deleteDoc))

}
