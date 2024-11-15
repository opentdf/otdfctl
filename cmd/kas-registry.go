package cmd

import (
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	keyRemote = "Remote"
	keyCached = "Cached"
)

var policy_kasRegistryCmd *cobra.Command

func policy_getKeyAccessRegistry(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.FlagHelper.GetRequiredID("id")

	kas, err := h.GetKasRegistryEntry(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get Registered KAS entry (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	keyType := keyCached
	key := &policy.PublicKey{}
	key.PublicKey = &policy.PublicKey_Cached{Cached: kas.GetPublicKey().GetCached()}
	if kas.GetPublicKey().GetRemote() != "" {
		keyType = keyRemote
		key.PublicKey = &policy.PublicKey_Remote{Remote: kas.GetPublicKey().GetRemote()}
	}
	rows := [][]string{
		{"Id", kas.GetId()},
		{"URI", kas.GetUri()},
		{"PublicKey Type", keyType},
		{"PublicKey", kas.GetPublicKey().String()},
	}

	if mdRows := getMetadataRows(kas.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)

	HandleSuccess(cmd, kas.GetId(), t, kas)
}

func policy_listKeyAccessRegistries(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	list, err := h.ListKasRegistryEntries()
	if err != nil {
		cli.ExitWithError("Failed to list Registered KAS entries", err)
	}

	t := cli.NewTable(
		cli.NewUUIDColumn(),
		table.NewFlexColumn("uri", "URI", cli.FlexColumnWidthFour),
		table.NewFlexColumn("pk_loc", "PublicKey Location", cli.FlexColumnWidthThree),
		table.NewFlexColumn("pk", "PublicKey", cli.FlexColumnWidthThree),
	)
	rows := []table.Row{}
	for _, kas := range list {
		keyType := keyCached
		key := policy.PublicKey{}
		key.PublicKey = &policy.PublicKey_Cached{Cached: kas.GetPublicKey().GetCached()}
		if kas.GetPublicKey().GetRemote() != "" {
			keyType = keyRemote
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
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	uri := c.Flags.GetRequiredString("uri")
	cachedJSON := c.Flags.GetOptionalString("public-keys")
	remote := c.Flags.GetOptionalString("public-key-remote")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	if cachedJSON == "" && remote == "" {
		e := fmt.Errorf("a public key is required. Please pass either a cached or remote public key")
		cli.ExitWithError("Issue with create flags 'public-keys' and 'public-key-remote'", e)
	}

	key := &policy.PublicKey{}
	keyType := keyCached
	if cachedJSON != "" {
		if remote != "" {
			e := fmt.Errorf("only one public key is allowed. Please pass either a cached or remote public key but not both")
			cli.ExitWithError("Issue with create flags 'public-keys' and 'public-key-remote'", e)
		}
		var err error
		key, err = parseKASRegistryPublicKey(cachedJSON)
		if err != nil {
			cli.ExitWithError("KAS registry key is invalid, see help for examples", err)
		}
	} else {
		keyType = keyRemote
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

	keyJSON, err := protojson.Marshal(key)
	if err != nil {
		cli.ExitWithError("Failed to marshal public key to JSON", err)
	}

	rows := [][]string{
		{"Id", created.GetId()},
		{"URI", created.GetUri()},
		{"PublicKey Type", keyType},
		{"PublicKey", string(keyJSON)},
	}
	if mdRows := getMetadataRows(created.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)

	HandleSuccess(cmd, created.GetId(), t, created)
}

func policy_updateKeyAccessRegistry(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	uri := c.Flags.GetOptionalString("uri")
	cachedJSON := c.Flags.GetOptionalString("public-keys")
	remote := c.Flags.GetOptionalString("public-key-remote")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	if cachedJSON == "" && remote == "" && len(metadataLabels) == 0 && uri == "" {
		cli.ExitWithError("No values were passed to update. Please pass at least one value to update (E.G. 'uri', 'public-keys', 'public-key-remote', 'label')", nil)
	}

	var pubKey *policy.PublicKey
	//nolint:gocritic // this is more readable than a switch statement
	if cachedJSON != "" && remote != "" {
		e := fmt.Errorf("only one public key is allowed. Please pass either a cached or remote public key but not both")
		cli.ExitWithError("Issue with update flags 'public-keys' and 'public-key-remote': ", e)
	} else if cachedJSON != "" {
		var err error
		pubKey, err = parseKASRegistryPublicKey(cachedJSON)
		if err != nil {
			cli.ExitWithError("KAS registry key is invalid, see help for examples", err)
		}
	} else if remote != "" {
		pubKey = &policy.PublicKey{PublicKey: &policy.PublicKey_Remote{Remote: remote}}
	}

	updated, err := h.UpdateKasRegistryEntry(
		id,
		uri,
		pubKey,
		getMetadataMutable(metadataLabels),
		getMetadataUpdateBehavior(),
	)
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to update Registered KAS entry (%s)", id), err)
	}
	rows := [][]string{
		{"Id", id},
		{"URI", updated.GetUri()},
	}
	if mdRows := getMetadataRows(updated.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, id, t, updated)
}

func policy_deleteKeyAccessRegistry(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	force := c.Flags.GetOptionalBool("force")

	kas, err := h.GetKasRegistryEntry(id)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to get Registered KAS entry (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	if !force {
		cli.ConfirmAction(cli.ActionDelete, "Registered KAS", id, false)
	}

	if _, err := h.DeleteKasRegistryEntry(id); err != nil {
		errMsg := fmt.Sprintf("Failed to delete Registered KAS entry (%s)", id)
		cli.ExitWithError(errMsg, err)
	}

	t := cli.NewTabular(
		[]string{"Id", kas.GetId()},
		[]string{"URI", kas.GetUri()},
	)

	HandleSuccess(cmd, kas.GetId(), t, kas)
}

// TODO remove this when the data is structured
func parseKASRegistryPublicKey(keyStr string) (*policy.PublicKey, error) {
	cachedKeys := new(policy.PublicKey)

	if !json.Valid([]byte(keyStr)) {
		return nil, errors.New("invalid JSON")
	}

	if err := protojson.Unmarshal([]byte(keyStr), cachedKeys); err != nil {
		return nil, errors.New("invalid shape")
	}

	// Validate all PEMs
	keyErrs := []error{}
	for i, k := range cachedKeys.GetCached().GetKeys() {
		block, _ := pem.Decode([]byte(k.GetPem()))
		if block == nil {
			keyErrs = append(keyErrs, fmt.Errorf("error in key[%d] with KID \"%s\": PEM is invalid", i, k.GetKid()))
		}
	}

	if len(keyErrs) != 0 {
		return nil, errors.Join(keyErrs...)
	}

	return cachedKeys, nil
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
		createDoc.GetDocFlag("public-keys").Name,
		createDoc.GetDocFlag("public-keys").Shorthand,
		createDoc.GetDocFlag("public-keys").Default,
		createDoc.GetDocFlag("public-keys").Description,
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
		updateDoc.GetDocFlag("public-keys").Name,
		updateDoc.GetDocFlag("public-keys").Shorthand,
		updateDoc.GetDocFlag("public-keys").Default,
		updateDoc.GetDocFlag("public-keys").Description,
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
	deleteDoc.Flags().Bool(
		deleteDoc.GetDocFlag("force").Name,
		false,
		deleteDoc.GetDocFlag("force").Description,
	)

	doc := man.Docs.GetCommand("policy/kas-registry",
		man.WithSubcommands(createDoc, getDoc, listDoc, updateDoc, deleteDoc),
	)
	policy_kasRegistryCmd = &doc.Command
	policyCmd.AddCommand(policy_kasRegistryCmd)
}
