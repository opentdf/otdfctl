package cmd

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/lib/ocrypto"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/proto"
)

const (
	defaultAlg    = 0
	defaultMode   = 0
	defaultStatus = 0
)

var policy_kasRegistryKeysCmd *cobra.Command

func generateKeys(wrappingKey string, alg policy.Algorithm) ([]byte, []byte, error) {
	wrappingKeyBytes, err := ocrypto.Base64Decode([]byte(wrappingKey))
	if err != nil {
		return nil, nil, errors.Join(errors.New("Failed to decode wrapping key"), err)
	}

	aesKey, err := ocrypto.NewAESGcm(wrappingKeyBytes)
	if err != nil {
		return nil, nil, errors.Join(errors.New("Failed to create AES key"), err)
	}

	kek, err := generateKeyPair(alg)
	if err != nil {
		return nil, nil, errors.Join(errors.New("Failed to generate key pair"), err)
	}

	kekPrivPem, err := kek.PrivateKeyInPemFormat()
	if err != nil {
		return nil, nil, errors.Join(errors.New("Failed to get private key in pem format"), err)
	}

	kekPubPem, err := kek.PublicKeyInPemFormat()
	if err != nil {
		return nil, nil, errors.Join(errors.New("Failed to get public key in pem format"), err)
	}

	wrappedKek, err := aesKey.Encrypt([]byte(kekPrivPem))
	if err != nil {
		return nil, nil, errors.Join(errors.New("Failed to wrap key"), err)
	}

	return wrappedKek, []byte(kekPubPem), nil
}

func generateKeyPair(alg policy.Algorithm) (ocrypto.KeyPair, error) {
	var key ocrypto.KeyPair
	var err error
	switch alg {
	case policy.Algorithm_ALGORITHM_RSA_2048:
		key, err = generateRSAKey(2048)
	case policy.Algorithm_ALGORITHM_RSA_4096:
		key, err = generateRSAKey(4096)
	case policy.Algorithm_ALGORITHM_EC_P256:
		key, err = generateECCKey(256)
	case policy.Algorithm_ALGORITHM_EC_P384:
		key, err = generateECCKey(384)
	case policy.Algorithm_ALGORITHM_EC_P521:
		key, err = generateECCKey(521)
	default:
		return nil, errors.New("unsupported algorithm")
	}

	return key, err
}

func generateRSAKey(size int) (ocrypto.RsaKeyPair, error) {
	return ocrypto.NewRSAKeyPair(size)
}

func generateECCKey(size int) (ocrypto.ECKeyPair, error) {
	mode, err := ocrypto.ECSizeToMode(size)
	if err != nil {
		return ocrypto.ECKeyPair{}, err
	}

	return ocrypto.NewECKeyPair(mode)
}

func enumToStatus(enum policy.KeyStatus) (string, error) {
	switch enum { //nolint:exhaustive // UNSPECIFIED is not needed here
	case policy.KeyStatus_KEY_STATUS_ACTIVE:
		return "active", nil
	case policy.KeyStatus_KEY_STATUS_INACTIVE:
		return "inactive", nil
	case policy.KeyStatus_KEY_STATUS_COMPROMISED:
		return "compromised", nil
	default:
		return "", errors.New("invalid enum status")
	}
}

func statusToEnum(status string) (policy.KeyStatus, error) {
	switch strings.ToLower(status) {
	case "active":
		return policy.KeyStatus_KEY_STATUS_ACTIVE, nil
	case "inactive":
		return policy.KeyStatus_KEY_STATUS_INACTIVE, nil
	case "compromised":
		return policy.KeyStatus_KEY_STATUS_COMPROMISED, nil
	default:
		return policy.KeyStatus_KEY_STATUS_UNSPECIFIED, errors.New("invalid status")
	}
}

func enumToMode(enum policy.KeyMode) (string, error) {
	switch enum { //nolint:exhaustive // UNSPECIFIED is not needed here
	case policy.KeyMode_KEY_MODE_LOCAL:
		return "local", nil
	case policy.KeyMode_KEY_MODE_REMOTE:
		return "remote", nil
	default:
		return "", errors.New("invalid enum mode")
	}
}

func modeToEnum(mode string) (policy.KeyMode, error) {
	switch strings.ToLower(mode) {
	case "local":
		return policy.KeyMode_KEY_MODE_LOCAL, nil
	case "remote":
		return policy.KeyMode_KEY_MODE_REMOTE, nil
	default:
		return policy.KeyMode_KEY_MODE_UNSPECIFIED, errors.New("invalid mode")
	}
}

func enumToAlg(enum policy.Algorithm) (string, error) {
	switch enum { //nolint:exhaustive // UNSPECIFIED is not needed here
	case policy.Algorithm_ALGORITHM_RSA_2048:
		return "rsa:2048", nil
	case policy.Algorithm_ALGORITHM_RSA_4096:
		return "rsa:4096", nil
	case policy.Algorithm_ALGORITHM_EC_P256:
		return "ec:secp256r1", nil
	case policy.Algorithm_ALGORITHM_EC_P384:
		return "ec:secp384r1", nil
	case policy.Algorithm_ALGORITHM_EC_P521:
		return "ec:secp521r1", nil
	default:
		return "", errors.New("invalid enum algorithm")
	}
}

func algToEnum(alg string) (policy.Algorithm, error) {
	switch strings.ToLower(alg) {
	case "rsa:2048":
		return policy.Algorithm_ALGORITHM_RSA_2048, nil
	case "rsa:4096":
		return policy.Algorithm_ALGORITHM_RSA_4096, nil
	case "ec:secp256r1":
		return policy.Algorithm_ALGORITHM_EC_P256, nil
	case "ec:secp384r1":
		return policy.Algorithm_ALGORITHM_EC_P384, nil
	case "ec:secp521r1":
		return policy.Algorithm_ALGORITHM_EC_P521, nil
	default:
		return policy.Algorithm_ALGORITHM_UNSPECIFIED, errors.New("invalid algorithm")
	}
}

func getTableRows(asymkey *policy.AsymmetricKey) [][]string {
	var providerConfig []byte
	var err error
	if asymkey.GetProviderConfig() != nil {
		providerConfig, err = proto.Marshal(asymkey.GetProviderConfig())
		if err != nil {
			cli.ExitWithError("Failed to marshal provider config", err)
		}
	}

	statusStr, err := enumToStatus(asymkey.GetKeyStatus())
	if err != nil {
		cli.ExitWithError("Failed to convert status", err)
	}
	modeStr, err := enumToMode(asymkey.GetKeyMode())
	if err != nil {
		cli.ExitWithError("Failed to convert mode", err)
	}
	algStr, err := enumToAlg(asymkey.GetKeyAlgorithm())
	if err != nil {
		cli.ExitWithError("Failed to convert algorithm", err)
	}

	rows := [][]string{
		{"ID", asymkey.GetId()},
		{"KeyId", asymkey.GetKeyId()},
		{"Algorithm", algStr},
		{"Status", statusStr},
		{"Mode", modeStr},
		{"PubKeyCtx", string(asymkey.GetPublicKeyCtx())},
		{"PrivateKeyCtx", string(asymkey.GetPrivateKeyCtx())},
		{"ProviderConfig", string(providerConfig)},
	}
	return rows
}

// TODO: Handle wrapping the generated key with provider config.
func policy_createKASKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	kasId := c.Flags.GetRequiredString("kasId")
	keyId := c.Flags.GetRequiredString("keyId")
	alg, err := algToEnum(c.Flags.GetRequiredString("alg"))
	if err != nil {
		cli.ExitWithError("Invalid algorithm", err)
	}

	mode, err := modeToEnum(c.Flags.GetRequiredString("mode"))
	if err != nil {
		cli.ExitWithError("Invalid mode", err)
	}

	providerConfigId := c.Flags.GetOptionalString("providerConfigId")
	publicKeyCtxArg := c.Flags.GetOptionalString("publicKeyCtx")
	privateKeyCtxArg := c.Flags.GetOptionalString("privateKeyCtx")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	var pubKeyCtxBytes []byte
	var privKeyCtxBytes []byte
	if mode == policy.KeyMode_KEY_MODE_LOCAL {
		wrappingKey := c.Flags.GetRequiredString("wrappingKey")
		privateKey, publicKey, err := generateKeys(wrappingKey, alg)
		if err != nil {
			cli.ExitWithError("Failed to generate keys", err)
		}

		// Unmarshal the public key context.
		pubKeyCtx := map[string]any{}
		if publicKeyCtxArg != "" {
			err := json.Unmarshal([]byte(publicKeyCtxArg), &pubKeyCtx)
			if err != nil {
				cli.ExitWithError("Failed to unmarshal public key context", err)
			}
		}
		// Unmarshal the private key context.
		privKeyCtx := map[string]any{}
		if privateKeyCtxArg != "" {
			err := json.Unmarshal([]byte(privateKeyCtxArg), &privKeyCtx)
			if err != nil {
				cli.ExitWithError("Failed to unmarshal private key context", err)
			}
		}

		pubKeyCtx["pubKey"] = string(ocrypto.Base64Encode(publicKey))
		pubKeyCtxBytes, err = json.Marshal(pubKeyCtx)
		if err != nil {
			cli.ExitWithError("Failed to marshal public key context", err)
		}

		privKeyCtx["wrappedKey"] = string(ocrypto.Base64Encode(privateKey))
		privKeyCtxBytes, err = json.Marshal(privKeyCtx)
		if err != nil {
			cli.ExitWithError("Failed to marshal private key context", err)
		}
	} else {
		if publicKeyCtxArg == "" {
			cli.ExitWithError("Public key context is required for remote mode", nil)
		}
		// REMOTE
		pubKeyCtxBytes = []byte(publicKeyCtxArg)
		privKeyCtxBytes = []byte(privateKeyCtxArg)
	}

	asymkey, err := h.CreateKasKey(
		c.Context(),
		kasId,
		keyId,
		alg,
		mode,
		pubKeyCtxBytes,
		privKeyCtxBytes,
		providerConfigId,
		getMetadataMutable(metadataLabels),
	)
	if err != nil {
		cli.ExitWithError("Failed to create kas key", err)
	}

	rows := getTableRows(asymkey)
	if mdRows := getMetadataRows(asymkey.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, asymkey.GetId(), t, asymkey)
}

func policy_getKASKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetOptionalID("id")
	keyId := c.Flags.GetOptionalString("keyId")

	asymkey, err := h.GetKasKey(c.Context(), id, keyId)
	if err != nil {
		cli.ExitWithError("Failed to get kas key", err)
	}

	rows := getTableRows(asymkey)
	if mdRows := getMetadataRows(asymkey.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, asymkey.GetId(), t, asymkey)
}

func policy_updateKASKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	keyStatusArg := c.Flags.GetOptionalString("status")
	var err error
	var keyStatus policy.KeyStatus
	if keyStatusArg != "" {
		keyStatus, err = statusToEnum(keyStatusArg)
		if err != nil {
			cli.ExitWithError("Invalid status", err)
		}
	}
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	if keyStatus == policy.KeyStatus_KEY_STATUS_UNSPECIFIED && len(metadataLabels) == 0 {
		cli.ExitWithError("Either status or metadata labels must be specified", nil)
	}

	resp, err := h.UpdateKasKey(
		c.Context(),
		id,
		keyStatus,
		getMetadataMutable(metadataLabels),
		getMetadataUpdateBehavior())
	if err != nil {
		cli.ExitWithError("Failed to update kas key", err)
	}

	// Get KAS Key.
	asymkey, err := h.GetKasKey(c.Context(), resp.GetId(), "")
	if err != nil {
		cli.ExitWithError("Failed to get kas key", err)
	}

	rows := getTableRows(asymkey)
	if mdRows := getMetadataRows(asymkey.GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, asymkey.GetId(), t, asymkey)
}

func policy_listKASKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")
	algArg := c.Flags.GetOptionalString("alg")
	var alg policy.Algorithm
	if algArg != "" {
		var err error
		alg, err = algToEnum(algArg)
		if err != nil {
			cli.ExitWithError("Invalid algorithm", err)
		}
	}
	kasId := c.Flags.GetOptionalString("kasId")
	kasName := c.Flags.GetOptionalString("kasName")
	kasUri := c.Flags.GetOptionalString("kasUri")

	// Get the list of keys.
	keys, page, err := h.ListKasKeys(c.Context(), limit, offset, alg, kasId, kasName, kasUri)
	if err != nil {
		cli.ExitWithError("Failed to list kas keys", err)
	}

	t := cli.NewTable(
		// columns should be id, name, config, labels, created_at, updated_at
		table.NewFlexColumn("id", "ID", cli.FlexColumnWidthOne),
		table.NewFlexColumn("keyId", "Key ID", cli.FlexColumnWidthOne),
		table.NewFlexColumn("keyAlgorithm", "Key Algorithm", cli.FlexColumnWidthOne),
		table.NewFlexColumn("keyStatus", "Key Status", cli.FlexColumnWidthOne),
		table.NewFlexColumn("keyMode", "Key Mode", cli.FlexColumnWidthOne),
		table.NewFlexColumn("pubKeyCtx", "Public Key Context", cli.FlexColumnWidthThree),
		table.NewFlexColumn("privateKeyCtx", "Private Key Context", cli.FlexColumnWidthThree),
		table.NewFlexColumn("providerConfig", "Provider Configuration", cli.FlexColumnWidthThree),
		table.NewFlexColumn("labels", "Labels", cli.FlexColumnWidthOne),
		table.NewFlexColumn("created_at", "Created At", cli.FlexColumnWidthOne),
		table.NewFlexColumn("updated_at", "Updated At", cli.FlexColumnWidthOne),
	)
	rows := []table.Row{}
	for _, key := range keys {
		metadata := cli.ConstructMetadata(key.GetMetadata())
		var providerConfig []byte
		if key.GetProviderConfig() != nil {
			providerConfig, err = proto.Marshal(key.GetProviderConfig())
			if err != nil {
				cli.ExitWithError("Failed to marshal provider config", err)
			}
		}
		statusStr, err := enumToStatus(key.GetKeyStatus())
		if err != nil {
			cli.ExitWithError("Failed to convert status", err)
		}
		modeStr, err := enumToMode(key.GetKeyMode())
		if err != nil {
			cli.ExitWithError("Failed to convert mode", err)
		}
		algStr, err := enumToAlg(key.GetKeyAlgorithm())
		if err != nil {
			cli.ExitWithError("Failed to convert algorithm", err)
		}
		rows = append(rows, table.NewRow(table.RowData{
			"id":             key.GetId(),
			"keyId":          key.GetKeyId(),
			"keyAlgorithm":   algStr,
			"keyStatus":      statusStr,
			"keyMode":        modeStr,
			"pubKeyCtx":      string(key.GetPublicKeyCtx()),
			"privateKeyCtx":  string(key.GetPrivateKeyCtx()),
			"providerConfig": string(providerConfig),
			"labels":         metadata["Labels"],
			"created_at":     metadata["Created At"],
			"updated_at":     metadata["Updated At"],
		}))
	}
	t = t.WithRows(rows)
	t = cli.WithListPaginationFooter(t, page)
	HandleSuccess(cmd, "", t, keys)
}

func init() {

	// Create Kas Key
	createDoc := man.Docs.GetCommand("policy/kas-registry/key/create",
		man.WithRun(policy_createKASKey),
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("kasId").Name,
		createDoc.GetDocFlag("kasId").Shorthand,
		createDoc.GetDocFlag("kasId").Default,
		createDoc.GetDocFlag("kasId").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("keyId").Name,
		createDoc.GetDocFlag("keyId").Shorthand,
		createDoc.GetDocFlag("keyId").Default,
		createDoc.GetDocFlag("keyId").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("alg").Name,
		createDoc.GetDocFlag("alg").Shorthand,
		createDoc.GetDocFlag("alg").Default,
		createDoc.GetDocFlag("alg").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("mode").Name,
		createDoc.GetDocFlag("mode").Shorthand,
		createDoc.GetDocFlag("mode").Default,
		createDoc.GetDocFlag("mode").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("wrappingKey").Name,
		createDoc.GetDocFlag("wrappingKey").Shorthand,
		createDoc.GetDocFlag("wrappingKey").Default,
		createDoc.GetDocFlag("wrappingKey").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("privateKeyCtx").Name,
		createDoc.GetDocFlag("privateKeyCtx").Shorthand,
		createDoc.GetDocFlag("privateKeyCtx").Default,
		createDoc.GetDocFlag("privateKeyCtx").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("publicKeyCtx").Name,
		createDoc.GetDocFlag("publicKeyCtx").Shorthand,
		createDoc.GetDocFlag("publicKeyCtx").Default,
		createDoc.GetDocFlag("publicKeyCtx").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("providerConfigId").Name,
		createDoc.GetDocFlag("providerConfigId").Shorthand,
		createDoc.GetDocFlag("providerConfigId").Default,
		createDoc.GetDocFlag("providerConfigId").Description,
	)
	injectLabelFlags(&createDoc.Command, false)

	// Get Kas Key
	getDoc := man.Docs.GetCommand("policy/kas-registry/key/get",
		man.WithRun(policy_getKASKey),
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("id").Name,
		getDoc.GetDocFlag("id").Shorthand,
		getDoc.GetDocFlag("id").Default,
		getDoc.GetDocFlag("id").Description,
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("keyId").Name,
		getDoc.GetDocFlag("keyId").Shorthand,
		getDoc.GetDocFlag("keyId").Default,
		getDoc.GetDocFlag("keyId").Description,
	)
	getDoc.MarkFlagsOneRequired("id", "keyId")
	getDoc.MarkFlagsMutuallyExclusive("id", "keyId")

	// Update Kas Key
	updateDoc := man.Docs.GetCommand("policy/kas-registry/key/update",
		man.WithRun(policy_updateKASKey),
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("status").Name,
		updateDoc.GetDocFlag("status").Shorthand,
		updateDoc.GetDocFlag("status").Default,
		updateDoc.GetDocFlag("status").Description,
	)
	injectLabelFlags(&updateDoc.Command, true)

	// List Kas Keys
	listDoc := man.Docs.GetCommand("policy/kas-registry/key/list",
		man.WithRun(policy_listKASKey),
	)
	listDoc.Flags().StringP(
		listDoc.GetDocFlag("alg").Name,
		listDoc.GetDocFlag("alg").Shorthand,
		listDoc.GetDocFlag("alg").Default,
		listDoc.GetDocFlag("alg").Description,
	)
	listDoc.Flags().StringP(
		listDoc.GetDocFlag("kasId").Name,
		listDoc.GetDocFlag("kasId").Shorthand,
		listDoc.GetDocFlag("kasId").Default,
		listDoc.GetDocFlag("kasId").Description,
	)
	listDoc.Flags().StringP(
		listDoc.GetDocFlag("kasName").Name,
		listDoc.GetDocFlag("kasName").Shorthand,
		listDoc.GetDocFlag("kasName").Default,
		listDoc.GetDocFlag("kasName").Description,
	)
	listDoc.Flags().StringP(
		listDoc.GetDocFlag("kasUri").Name,
		listDoc.GetDocFlag("kasUri").Shorthand,
		listDoc.GetDocFlag("kasUri").Default,
		listDoc.GetDocFlag("kasUri").Description,
	)
	injectListPaginationFlags(listDoc)
	listDoc.MarkFlagsMutuallyExclusive("kasId", "kasName", "kasUri")

	doc := man.Docs.GetCommand("policy/kas-registry/key",
		man.WithSubcommands(createDoc, getDoc, updateDoc, listDoc))
	policy_kasRegistryKeysCmd = &doc.Command
	policy_kasRegistryCmd.AddCommand(policy_kasRegistryKeysCmd)
}
