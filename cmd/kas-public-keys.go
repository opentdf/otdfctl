package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/term"
	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/kasregistry"
	"github.com/spf13/cobra"
)

var policy_kasPublicKeyCmd = man.Docs.GetCommand("policy/kas-registry/public-keys")

func policy_createPublicKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	kas := c.Flags.GetRequiredString("kas")
	publicKey := c.Flags.GetRequiredString("key")
	alg := c.Flags.GetRequiredString("algorithm")
	kid := c.Flags.GetRequiredString("key-id")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	publicKey, err := parseAndFormatKey(publicKey)
	if err != nil {
		cli.ExitWithError("Failed to parse public key", err)
	}

	created, err := h.CreatePublicKey(kas, publicKey, kid, alg, getMetadataMutable(metadataLabels))
	if err != nil {
		cli.ExitWithError("Failed to create public key", err)
	}

	rows := [][]string{
		{"Id", created.GetId()},
		{"Key ID", created.GetPublicKey().GetKid()},
		{"Algorithm", alg},
		{"Public Key", created.GetPublicKey().GetPem()},
		{"Was Mapped", created.GetWasMapped().String()},
		{"Active", created.GetIsActive().String()},
	}

	if mdRows := getMetadataRows(created.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)

	HandleSuccess(cmd, created.GetId(), t, created)
}

func policy_updatePublicKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})
	updated, err := h.UpdatePublicKey(id, getMetadataMutable(metadataLabels), getMetadataUpdateBehavior())
	if err != nil {
		cli.ExitWithError("Failed to update public key", err)
	}

	rows := [][]string{
		{"Id", updated.GetId()},
		{"Key ID", updated.GetPublicKey().GetKid()},
	}

	if mdRows := getMetadataRows(updated.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	t := cli.NewTabular(rows...)

	HandleSuccess(cmd, updated.GetId(), t, updated)
}

func policy_activePublicKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	err := h.ActivatePublicKey(id)
	if err != nil {
		cli.ExitWithError("Failed to active public key", err)
	}

	HandleSuccess(cmd, id, table.Model{}, nil)
}

func policy_deactivePublicKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	err := h.DeactivatePublicKey(id)
	if err != nil {
		cli.ExitWithError("Failed to deactivate public key", err)
	}

	HandleSuccess(cmd, id, table.Model{}, nil)
}

func policy_unsafeDeletePublicKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	force := c.Flags.GetOptionalBool("force")

	pk, err := h.GetPublicKey(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get public key with id %s", id), err)
	}

	if !force {
		cli.ConfirmTextInput(cli.ActionDelete, "public-key", cli.InputNameKeyID, pk.GetPublicKey().GetKid())
	}

	err = h.UnsafeDeletePublicKey(id)
	if err != nil {
		cli.ExitWithError("Failed to delete public key", err)
	}

	rows := [][]string{
		{"Id", pk.GetId()},
		{"Key ID", pk.GetPublicKey().GetKid()},
	}
	if mdRows := getMetadataRows(pk.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, pk.GetId(), t, pk)
}

func policy_getPublicKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")

	key, err := h.GetPublicKey(id)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to get public key with id %s", id), err)
	}

	alg, err := enumToAlg(key.GetPublicKey().GetAlg())
	if err != nil {
		cli.ExitWithError("Failed to get algorithm", err)
	}

	rows := [][]string{
		{"Id", key.GetId()},
		{"Was Mapped", fmt.Sprintf("%t", key.GetWasMapped().GetValue())},
		{"Active", fmt.Sprintf("%t", key.GetIsActive().GetValue())},
		{"KAS Name", key.GetKas().GetName()},
		{"KAS URI", key.GetKas().GetUri()},
		{"Key ID", key.GetPublicKey().GetKid()},
		{"Algorithm", alg},
		{"Public Key", key.GetPublicKey().GetPem()},
	}

	if mdRows := getMetadataRows(key.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}

	cli.NewTable(
		cli.NewUUIDColumn(),
	)

	t := cli.NewTabular(rows...)

	HandleSuccess(cmd, key.GetId(), t, key)
}

func policy_listPublicKeys(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	kas := c.Flags.GetOptionalString("kas")
	showPublicKey := c.Flags.GetOptionalBool("show-public-key")
	offset := c.Flags.GetRequiredInt32("offset")
	limit := c.Flags.GetRequiredInt32("limit")

	keys, page, err := h.ListPublicKeys(kas, offset, limit)
	if err != nil {
		cli.ExitWithError("Failed to list public keys", err)
	}

	columns := []table.Column{
		table.NewFlexColumn("id", "ID", cli.FlexColumnWidthThree),
		table.NewFlexColumn("is_active", "Active", cli.FlexColumnWidthTwo),
		table.NewFlexColumn("was_mapped", "Was Mapped", cli.FlexColumnWidthTwo),
		table.NewFlexColumn("kas_name", "KAS Name", cli.FlexColumnWidthThree),
		table.NewFlexColumn("kas_uri", "KAS URI", cli.FlexColumnWidthThree),
		table.NewFlexColumn("key_id", "Key ID", cli.FlexColumnWidthTwo),
		table.NewFlexColumn("algorithm", "Algorithm", cli.FlexColumnWidthTwo),
	}

	if showPublicKey {
		columns = append(columns, table.NewFlexColumn("public_key", "Public Key", cli.FlexColumnWidthFour))
	}

	t := cli.NewTable(columns...)

	rows := []table.Row{}
	for _, key := range keys {
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
			"id":         key.GetId(),
			"is_active":  key.GetIsActive().GetValue(),
			"was_mapped": key.GetWasMapped().GetValue(),
			"kas_id":     key.GetKas().GetId(),
			"kas_name":   key.GetKas().GetName(),
			"kas_uri":    key.GetKas().GetUri(),
			"key_id":     key.GetPublicKey().GetKid(),
			"algorithm":  alg,
			"public_key": key.GetPublicKey().GetPem(),
		}

		rows = append(rows, table.NewRow(rd).WithStyle(rowStyle))
	}

	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)

	HandleSuccess(cmd, "", t, keys)
}

func policy_listPublicKeyMappings(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	kas := c.Flags.GetOptionalString("kas")
	pkID := c.Flags.GetOptionalID("public-key-id")
	offset := c.Flags.GetRequiredInt32("offset")
	limit := c.Flags.GetRequiredInt32("limit")

	mappings, page, err := h.ListPublicKeyMappings(kas, pkID, offset, limit)
	if err != nil {
		cli.ExitWithError("Failed to list public key mappings", err)
	}

	t := cli.NewTable(
		table.NewFlexColumn("kas_name", "KAS Name", cli.FlexColumnWidthTwo),
		table.NewFlexColumn("kas_uri", "KAS URI", cli.FlexColumnWidthThree),
		table.NewFlexColumn("key_count", "Key Count", cli.FlexColumnWidthOne),
		table.NewFlexColumn("publicKeys", "Public Keys", cli.FlexColumnWidthFour),
	)

	rows := []table.Row{}

	termWidth, _, err := term.GetSize(os.Stdout.Fd())
	if err != nil {
		termWidth = 80
	}

	for _, mapping := range mappings {
		rows = append(rows, table.NewRow(table.RowData{
			"kas_name":   mapping.GetKasName(),
			"kas_uri":    mapping.GetKasUri(),
			"key_count":  fmt.Sprintf("%d", len(mapping.GetPublicKeys())),
			"publicKeys": createPublicKeysTable(mapping.GetPublicKeys(), termWidth),
		}).WithStyle(lipgloss.NewStyle().BorderBottom(true).BorderStyle(lipgloss.NormalBorder())))
	}

	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)

	HandleSuccess(cmd, "", t, mappings)
}

func createPublicKeysTable(keys []*kasregistry.ListPublicKeyMappingResponse_PublicKey, termWidth int) string {
	// Create columns for the nested table
	columns := []table.Column{
		table.NewFlexColumn("kid", "KID", cli.FlexColumnWidthTwo),
		table.NewFlexColumn("algorithm", "Algorithm", cli.FlexColumnWidthTwo),
		table.NewFlexColumn("active", "Active", cli.FlexColumnWidthOne),
		table.NewFlexColumn("namespaces", "Namespaces", cli.FlexColumnWidthTwo),
		table.NewFlexColumn("definitions", "Definitions", cli.FlexColumnWidthThree),
		table.NewFlexColumn("values", "Values", cli.FlexColumnWidthFour),
	}

	var rows []table.Row

	for _, pk := range keys {
		alg, err := enumToAlg(pk.GetKey().GetPublicKey().GetAlg())
		if err != nil {
			cli.ExitWithError("Failed to get algorithm", err)
		}

		rowStyle := lipgloss.NewStyle().BorderBottom(true).BorderStyle(lipgloss.NormalBorder())

		if pk.GetKey().GetIsActive().GetValue() {
			rowStyle = rowStyle.Background(cli.ColorGreen.Background)
		} else {
			rowStyle = rowStyle.Background(cli.ColorRed.Background)
		}

		rows = append(rows, table.NewRow(table.RowData{
			"kid":         pk.GetKey().GetPublicKey().GetKid(),
			"algorithm":   alg,
			"active":      fmt.Sprintf("%v", pk.GetKey().GetIsActive().GetValue()),
			"namespaces":  formatAssociations(pk.GetNamespaces()),
			"definitions": formatAssociations(pk.GetDefinitions()),
			"values":      formatAssociations(pk.GetValues()),
		}).WithStyle(rowStyle))
	}

	minWidth := 80 // Set a minimum width for the nested table
	tableWidthPercentage := 0.75
	tableWidth := int(float64(termWidth) * tableWidthPercentage)
	if tableWidth < minWidth {
		tableWidth = minWidth
	}
	// Create nested table
	nestedTable := table.New(columns).
		WithRows(rows).
		WithTargetWidth(int(float64(tableWidth) * tableWidthPercentage)).
		WithMultiline(true).
		WithNoPagination().
		BorderRounded().
		WithBaseStyle(lipgloss.NewStyle().Align(lipgloss.Left))

	// Convert the table to string and add some indentation
	tableStr := nestedTable.View()

	// Add indentation to each line of the nested table
	indentedLines := strings.Split(tableStr, "\n")
	for i, line := range indentedLines {
		indentedLines[i] = "    " + line
	}
	return strings.Join(indentedLines, "\n")
}

func formatAssociations(assocs []*kasregistry.ListPublicKeyMappingResponse_Association) string {
	if len(assocs) == 0 {
		return "-"
	}
	var fqns []string
	for _, a := range assocs {
		// remove https:// from the beginning of the URI
		fqn, _ := strings.CutPrefix(a.GetFqn(), "https://")
		fqns = append(fqns, fqn)
	}
	return strings.Join(fqns, "\n")
}

func isValidBase64(s string) bool {
	_, err := base64.StdEncoding.DecodeString(s)
	return err == nil
}

func parseAndFormatKey(key string) (string, error) {
	if key == "" {
		return "", errors.New("key is required")
	}

	// If the key contains a newline, replace it with the actual newline character
	if strings.Contains(key, "\\n") {
		return strings.ReplaceAll(key, "\\n", "\n"), nil
	}

	// If the key is base64 encoded, decode it
	if isValidBase64(key) {
		decoded, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return "", err
		}
		return string(decoded), nil
	}

	return key, nil
}

func enumToAlg(enum policy.KasPublicKeyAlgEnum) (string, error) {
	switch enum { //nolint:exhaustive // UNSPECIFIED is not needed here
	case policy.KasPublicKeyAlgEnum_KAS_PUBLIC_KEY_ALG_ENUM_RSA_2048:
		return "rsa:2048", nil
	case policy.KasPublicKeyAlgEnum_KAS_PUBLIC_KEY_ALG_ENUM_RSA_4096:
		return "rsa:4096", nil
	case policy.KasPublicKeyAlgEnum_KAS_PUBLIC_KEY_ALG_ENUM_EC_SECP256R1:
		return "ec:secp256r1", nil
	case policy.KasPublicKeyAlgEnum_KAS_PUBLIC_KEY_ALG_ENUM_EC_SECP384R1:
		return "ec:secp384r1", nil
	case policy.KasPublicKeyAlgEnum_KAS_PUBLIC_KEY_ALG_ENUM_EC_SECP521R1:
		return "ec:secp521r1", nil
	default:
		return "", errors.New("invalid enum algorithm")
	}
}

func init() {
	createDoc := man.Docs.GetCommand("policy/kas-registry/public-keys/create",
		man.WithRun(policy_createPublicKey))
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("kas").Name,
		createDoc.GetDocFlag("kas").Shorthand,
		createDoc.GetDocFlag("kas").Default,
		createDoc.GetDocFlag("kas").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("key").Name,
		createDoc.GetDocFlag("key").Shorthand,
		createDoc.GetDocFlag("key").Default,
		createDoc.GetDocFlag("key").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("algorithm").Name,
		createDoc.GetDocFlag("algorithm").Shorthand,
		createDoc.GetDocFlag("algorithm").Default,
		createDoc.GetDocFlag("algorithm").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("key-id").Name,
		createDoc.GetDocFlag("key-id").Shorthand,
		createDoc.GetDocFlag("key-id").Default,
		createDoc.GetDocFlag("key-id").Description,
	)
	injectLabelFlags(&createDoc.Command, false)

	updateDoc := man.Docs.GetCommand("policy/kas-registry/public-keys/update",
		man.WithRun(policy_updatePublicKey))
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	injectLabelFlags(&updateDoc.Command, true)

	activateDoc := man.Docs.GetCommand("policy/kas-registry/public-keys/activate",
		man.WithRun(policy_activePublicKey))
	activateDoc.Flags().StringP(
		activateDoc.GetDocFlag("id").Name,
		activateDoc.GetDocFlag("id").Shorthand,
		activateDoc.GetDocFlag("id").Default,
		activateDoc.GetDocFlag("id").Description,
	)

	deactivateDoc := man.Docs.GetCommand("policy/kas-registry/public-keys/deactivate",
		man.WithRun(policy_deactivePublicKey))
	deactivateDoc.Flags().StringP(
		deactivateDoc.GetDocFlag("id").Name,
		deactivateDoc.GetDocFlag("id").Shorthand,
		deactivateDoc.GetDocFlag("id").Default,
		deactivateDoc.GetDocFlag("id").Description,
	)

	getDoc := man.Docs.GetCommand("policy/kas-registry/public-keys/get",
		man.WithRun(policy_getPublicKey))
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Shorthand,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)

	listDoc := man.Docs.GetCommand("policy/kas-registry/public-keys/list",
		man.WithRun(policy_listPublicKeys))
	listDoc.Flags().StringP(
		listDoc.GetDocFlag("kas").Name,
		listDoc.GetDocFlag("kas").Shorthand,
		listDoc.GetDocFlag("kas").Default,
		listDoc.GetDocFlag("kas").Description,
	)
	listDoc.Flags().BoolP(
		listDoc.GetDocFlag("show-public-key").Name,
		listDoc.GetDocFlag("show-public-key").Shorthand,
		listDoc.GetDocFlag("show-public-key").DefaultAsBool(),
		listDoc.GetDocFlag("show-public-key").Description,
	)
	injectListPaginationFlags(listDoc)

	listMappingsDoc := man.Docs.GetCommand("policy/kas-registry/public-keys/list-mappings",
		man.WithRun(policy_listPublicKeyMappings))
	listMappingsDoc.Flags().StringP(
		listMappingsDoc.GetDocFlag("kas").Name,
		listMappingsDoc.GetDocFlag("kas").Shorthand,
		listMappingsDoc.GetDocFlag("kas").Default,
		listMappingsDoc.GetDocFlag("kas").Description,
	)
	listMappingsDoc.Flags().StringP(
		listMappingsDoc.GetDocFlag("public-key-id").Name,
		listMappingsDoc.GetDocFlag("public-key-id").Shorthand,
		listMappingsDoc.GetDocFlag("public-key-id").Default,
		listMappingsDoc.GetDocFlag("public-key-id").Description,
	)
	injectListPaginationFlags(listMappingsDoc)

	unsafeDeleteDoc := man.Docs.GetCommand("policy/kas-registry/public-keys/unsafe/delete",
		man.WithRun(policy_unsafeDeletePublicKey))
	unsafeDeleteDoc.Flags().StringP(
		unsafeDeleteDoc.GetDocFlag("id").Name,
		unsafeDeleteDoc.GetDocFlag("id").Shorthand,
		unsafeDeleteDoc.GetDocFlag("id").Default,
		unsafeDeleteDoc.GetDocFlag("id").Description,
	)
	unsafeDeleteDoc.Flags().BoolP(
		unsafeDeleteDoc.GetDocFlag("force").Name,
		unsafeDeleteDoc.GetDocFlag("force").Shorthand,
		unsafeDeleteDoc.GetDocFlag("force").DefaultAsBool(),
		unsafeDeleteDoc.GetDocFlag("force").Description,
	)

	policy_kasPublicKeyUnsafeCmd := man.Docs.GetCommand("policy/kas-registry/public-keys/unsafe")
	policy_kasPublicKeyUnsafeCmd.AddSubcommands(unsafeDeleteDoc)

	policy_kasPublicKeyCmd.AddCommand(&policy_kasPublicKeyUnsafeCmd.Command)

	policy_kasPublicKeyCmd.AddSubcommands(createDoc, updateDoc, getDoc, listDoc, listMappingsDoc, activateDoc, deactivateDoc)
	policy_kasRegistryCmd.AddCommand(&policy_kasPublicKeyCmd.Command)
}
