package cmd

import (
	"fmt"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

var policy_kasRegistryCmd *cobra.Command

func policy_getKeyAccessRegistry(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	kas, err := h.GetKasRegistryEntry(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get Registered KAS entry (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	keyType := "Cached"
	key := &policy.PublicKey{}
	key.PublicKey = &policy.PublicKey_Cached{Cached: kas.GetPublicKey().GetCached()}
	if kas.PublicKey.GetRemote() != "" {
		keyType = "Remote"
		key.PublicKey = &policy.PublicKey_Remote{Remote: kas.GetPublicKey().GetRemote()}
	}

	t := cli.NewTabular(
		[]string{"Id", kas.Id},
		// TODO: render labels [https://github.com/opentdf/otdfctl/issues/73]
		[]string{"URI", kas.Uri},
		[]string{"PublicKey Type", keyType},
		[]string{"PublicKey", kas.GetPublicKey().String()},
	)
	HandleSuccess(cmd, kas.Id, t, kas)
}

func policy_listKeyAccessRegistries(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	list, err := h.ListKasRegistryEntries()
	if err != nil {
		cli.ExitWithError("Failed to list Registered KAS entries", err)
	}

	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("uri", "URI", 4),
		table.NewFlexColumn("pk_loc", "PublicKey Location", 3),
		table.NewFlexColumn("pk", "PublicKey", 3),
	)
	rows := []table.Row{}
	for _, kas := range list {
		keyType := "Cached"
		key := policy.PublicKey{}
		key.PublicKey = &policy.PublicKey_Cached{Cached: kas.GetPublicKey().GetCached()}
		if kas.PublicKey.GetRemote() != "" {
			keyType = "Remote"
			key.PublicKey = &policy.PublicKey_Remote{Remote: kas.GetPublicKey().GetRemote()}
		}

		rows = append(rows, table.NewRow(table.RowData{
			"id":     kas.GetId(),
			"uri":    kas.GetUri(),
			"pk_loc": keyType,
			"pk":     kas.GetPublicKey().String(),
		}))
	}
	t = t.WithRows(rows)
	HandleSuccess(cmd, "", t, list)
}

func policy_createKeyAccessRegistry(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	uri := flagHelper.GetRequiredString("uri")
	cachedJSON := flagHelper.GetOptionalString("public-key-cached")
	remote := flagHelper.GetOptionalString("public-key-remote")
	metadataLabels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

	if cachedJSON == "" && remote == "" {
		e := fmt.Errorf("a public key is required. Please pass either a cached or remote public key")
		cli.ExitWithError("Issue with create flags 'public-key-cached' and 'public-key-remote': ", e)
	}

	key := &policy.PublicKey{}
	keyType := "Cached"
	if cachedJSON != "" {
		if remote != "" {
			e := fmt.Errorf("only one public key is allowed. Please pass either a cached or remote public key but not both")
			cli.ExitWithError("Issue with create flags 'public-key-cached' and 'public-key-remote': ", e)
		}
		cached := new(policy.PublicKey)
		err := protojson.Unmarshal([]byte(cachedJSON), cached)
		if err != nil {
			cli.ExitWithError("Failed to unmarshal cached public key JSON", err)
		}
		key = cached
	} else {
		keyType = "Remote"
		key.PublicKey = &policy.PublicKey_Remote{Remote: remote}
	}

	created, err := h.CreateKasRegistryEntry(
		uri,
		key,
		getMetadataMutable(metadataLabels),
	)
	if err != nil {
		cli.ExitWithError("Failed to create Registered KAS entry", err)
	}

	t := cli.NewTabular(
		[]string{"Id", created.Id},
		[]string{"URI", created.Uri},
		[]string{"PublicKey Type", keyType},
		[]string{"PublicKey", cachedJSON},
		// TODO: render labels [https://github.com/opentdf/otdfctl/issues/73]
	)

	HandleSuccess(cmd, created.Id, t, created)
}

func policy_updateKeyAccessRegistry(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)

	id := flagHelper.GetRequiredString("id")
	uri := flagHelper.GetOptionalString("uri")
	cachedJSON := flagHelper.GetOptionalString("public-key-cached")
	remote := flagHelper.GetOptionalString("public-key-remote")
	labels := flagHelper.GetStringSlice("label", metadataLabels, cli.FlagHelperStringSliceOptions{Min: 0})

	if cachedJSON == "" && remote == "" && len(labels) == 0 && uri == "" {
		cli.ExitWithError("No values were passed to update. Please pass at least one value to update (E.G. 'uri', 'public-key-cached', 'public-key-remote', 'label')", nil)
	}

	var pubKey *policy.PublicKey
	if cachedJSON != "" && remote != "" {
		e := fmt.Errorf("only one public key is allowed. Please pass either a cached or remote public key but not both")
		cli.ExitWithError("Issue with update flags 'public-key-cached' and 'public-key-remote': ", e)
	} else if cachedJSON != "" {
		cached := new(policy.PublicKey)
		err := protojson.Unmarshal([]byte(cachedJSON), cached)
		if err != nil {
			cli.ExitWithError("Failed to unmarshal cached public key JSON", err)
		}
		pubKey = cached
	} else if remote != "" {
		pubKey = &policy.PublicKey{PublicKey: &policy.PublicKey_Remote{Remote: remote}}
	}

	updated, err := h.UpdateKasRegistryEntry(
		id,
		uri,
		pubKey,
		getMetadataMutable(labels),
		getMetadataUpdateBehavior(),
	)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to update Registered KAS entry (%s)", id), err)
	}
	t := cli.NewTabular(
		[]string{"Id", id},
		[]string{"URI", updated.GetUri()},
		// TODO: render labels [https://github.com/opentdf/otdfctl/issues/73]
	)
	HandleSuccess(cmd, id, t, updated)
}

func policy_deleteKeyAccessRegistry(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	id := flagHelper.GetRequiredString("id")

	kas, err := h.GetKasRegistryEntry(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get Registered KAS entry (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	cli.ConfirmAction(cli.ActionDelete, "Registered KAS", id, false)

	if _, err := h.DeleteKasRegistryEntry(id); err != nil {
		errMsg := fmt.Sprintf("Failed to delete Registered KAS entry (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	t := cli.NewTabular(
		[]string{"Id", "URI"},
		[]string{kas.Id, kas.Uri},
	)

	HandleSuccess(cmd, kas.Id, t, kas)
}

func init() {
	getDoc := man.Docs.GetCommand("policy/kas-registry/get",
		man.WithRun(policy_getKeyAccessRegistry),
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Shorthand,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)

	listDoc := man.Docs.GetCommand("policy/kas-registry/list",
		man.WithRun(policy_listKeyAccessRegistries),
	)
	// TODO: active, inactive, any state querying [https://github.com/opentdf/otdfctl/issues/68]

	createDoc := man.Docs.GetCommand("policy/kas-registry/create",
		man.WithRun(policy_createKeyAccessRegistry),
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("uri").Name,
		createDoc.GetDocFlag("uri").Shorthand,
		createDoc.GetDocFlag("uri").Default,
		createDoc.GetDocFlag("uri").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("public-key-cached").Name,
		createDoc.GetDocFlag("public-key-cached").Shorthand,
		createDoc.GetDocFlag("public-key-cached").Default,
		createDoc.GetDocFlag("public-key-cached").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("public-key-remote").Name,
		createDoc.GetDocFlag("public-key-remote").Shorthand,
		createDoc.GetDocFlag("public-key-remote").Default,
		createDoc.GetDocFlag("public-key-remote").Description,
	)
	injectLabelFlags(&createDoc.Command, false)

	updateDoc := man.Docs.GetCommand("policy/kas-registry/update",
		man.WithRun(policy_updateKeyAccessRegistry),
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("uri").Name,
		updateDoc.GetDocFlag("uri").Shorthand,
		updateDoc.GetDocFlag("uri").Default,
		updateDoc.GetDocFlag("uri").Description,
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("public-key-cached").Name,
		updateDoc.GetDocFlag("public-key-cached").Shorthand,
		updateDoc.GetDocFlag("public-key-cached").Default,
		updateDoc.GetDocFlag("public-key-cached").Description,
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("public-key-remote").Name,
		updateDoc.GetDocFlag("public-key-remote").Shorthand,
		updateDoc.GetDocFlag("public-key-remote").Default,
		updateDoc.GetDocFlag("public-key-remote").Description,
	)
	injectLabelFlags(&updateDoc.Command, true)

	deleteDoc := man.Docs.GetCommand("policy/kas-registry/delete",
		man.WithRun(policy_deleteKeyAccessRegistry),
	)
	deleteDoc.Flags().StringP(
		deleteDoc.GetDocFlag("id").Name,
		deleteDoc.GetDocFlag("id").Shorthand,
		deleteDoc.GetDocFlag("id").Default,
		deleteDoc.GetDocFlag("id").Description,
	)

	doc := man.Docs.GetCommand("policy/kas-registry",
		man.WithSubcommands(createDoc, getDoc, listDoc, updateDoc, deleteDoc),
	)
	policy_kasRegistryCmd = &doc.Command
	policyCmd.AddCommand(policy_kasRegistryCmd)
}
