package cmd

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/lib/ocrypto"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/kasregistry"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	defaultAlg           = 0
	defaultMode          = 0
	defaultStatus        = 0
	keyStatusActive      = "active"
	keyStatusInactive    = "inactive"
	keyStatusCompromised = "compromised"
	keyModeLocal         = "local"
	keyModeProvider      = "provider"
	keyModeRemote        = "remote"
	keyModePublicKeyOnly = "public_key"
)

var policy_kasRegistryKeysCmd *cobra.Command

func wrapKey(key string, wrappingKey string) ([]byte, error) {
	wrappingKeyBytes, err := ocrypto.Base64Decode([]byte(wrappingKey))
	if err != nil {
		return nil, errors.Join(errors.New("failed to decode wrapping key"), err)
	}

	aesKey, err := ocrypto.NewAESGcm(wrappingKeyBytes)
	if err != nil {
		return nil, errors.Join(errors.New("failed to create AES key"), err)
	}

	wrappedKek, err := aesKey.Encrypt([]byte(key))
	if err != nil {
		return nil, errors.Join(errors.New("failed to wrap key"), err)
	}

	return wrappedKek, nil
}

func generateKeys(alg policy.Algorithm) (string, string, error) {
	kek, err := generateKeyPair(alg)
	if err != nil {
		return "", "", errors.Join(errors.New("failed to generate key pair"), err)
	}

	kekPrivPem, err := kek.PrivateKeyInPemFormat()
	if err != nil {
		return "", "", errors.Join(errors.New("failed to get private key in pem format"), err)
	}

	kekPubPem, err := kek.PublicKeyInPemFormat()
	if err != nil {
		return "", "", errors.Join(errors.New("failed to get public key in pem format"), err)
	}

	return kekPrivPem, kekPubPem, nil
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
		return keyStatusActive, nil
	case policy.KeyStatus_KEY_STATUS_INACTIVE:
		return keyStatusInactive, nil
	case policy.KeyStatus_KEY_STATUS_COMPROMISED:
		return keyStatusCompromised, nil
	default:
		return "", errors.New("invalid enum status")
	}
}

func statusToEnum(status string) (policy.KeyStatus, error) {
	switch strings.ToLower(status) {
	case keyStatusActive:
		return policy.KeyStatus_KEY_STATUS_ACTIVE, nil
	case keyStatusInactive:
		return policy.KeyStatus_KEY_STATUS_INACTIVE, nil
	case keyStatusCompromised:
		return policy.KeyStatus_KEY_STATUS_COMPROMISED, nil
	default:
		return policy.KeyStatus_KEY_STATUS_UNSPECIFIED, errors.New("invalid status")
	}
}

func enumToMode(enum policy.KeyMode) (string, error) {
	switch enum { //nolint:exhaustive // UNSPECIFIED is not needed here
	case policy.KeyMode_KEY_MODE_CONFIG_ROOT_KEY:
		return keyModeLocal, nil
	case policy.KeyMode_KEY_MODE_PROVIDER_ROOT_KEY:
		return keyModeProvider, nil
	case policy.KeyMode_KEY_MODE_REMOTE:
		return keyModeRemote, nil
	case policy.KeyMode_KEY_MODE_PUBLIC_KEY_ONLY:
		return keyModePublicKeyOnly, nil
	default:
		return "", errors.New("invalid enum mode")
	}
}

func modeToEnum(mode string) (policy.KeyMode, error) {
	switch strings.ToLower(mode) {
	case keyModeLocal:
		return policy.KeyMode_KEY_MODE_CONFIG_ROOT_KEY, nil
	case keyModeProvider:
		return policy.KeyMode_KEY_MODE_PROVIDER_ROOT_KEY, nil
	case keyModeRemote:
		return policy.KeyMode_KEY_MODE_REMOTE, nil
	case keyModePublicKeyOnly:
		return policy.KeyMode_KEY_MODE_PUBLIC_KEY_ONLY, nil
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

func getTableRows(kasKey *policy.KasKey) [][]string {
	var providerConfig []byte
	var err error
	asymkey := kasKey.GetKey()
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

	pubCtxBytes, err := protojson.Marshal(asymkey.GetPublicKeyCtx())
	if err != nil {
		cli.ExitWithError("Failed to marshal public key context", err)
	}
	privateKeyBytes, err := protojson.Marshal(asymkey.GetPrivateKeyCtx())
	if err != nil {
		cli.ExitWithError("Failed to marshal private key context", err)
	}

	rows := [][]string{
		{"ID", asymkey.GetId()},
		{"KasUri", kasKey.GetKasUri()},
		{"KeyId", asymkey.GetKeyId()},
		{"Algorithm", algStr},
		{"Status", statusStr},
		{"Mode", modeStr},
		{"PubKeyCtx", string(pubCtxBytes)},
		{"PrivateKeyCtx", string(privateKeyBytes)},
		{"ProviderConfig", string(providerConfig)},
	}
	return rows
}

// TODO: Handle wrapping the generated key with provider config.
func policy_createKASKey(cmd *cobra.Command, args []string) {
	var (
		wrappingKeyId    string
		providerConfigId string
	)

	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	keyId := c.Flags.GetRequiredString("keyId")
	kasId := c.Flags.GetOptionalString("kasId")
	kasUri := c.Flags.GetOptionalString("kasUri")
	kasName := c.Flags.GetOptionalString("kasName")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	alg, err := algToEnum(c.Flags.GetRequiredString("alg"))
	if err != nil {
		cli.ExitWithError("Invalid algorithm", err)
	}

	mode, err := modeToEnum(c.Flags.GetRequiredString("mode"))
	if err != nil {
		cli.ExitWithError("Invalid mode", err)
	}

	wrappingKeyId = c.Flags.GetOptionalString("wrappingKeyId")
	if mode != policy.KeyMode_KEY_MODE_PUBLIC_KEY_ONLY && wrappingKeyId == "" {
		formattedMode, _ := enumToMode(mode)
		cli.ExitWithError(fmt.Sprintf("wrappingKeyId is required for mode %s", formattedMode), nil)
	}

	providerConfigId = c.Flags.GetOptionalString("providerConfigId")
	if (mode == policy.KeyMode_KEY_MODE_PROVIDER_ROOT_KEY || mode == policy.KeyMode_KEY_MODE_REMOTE) && providerConfigId == "" {
		formattedMode, _ := enumToMode(mode)
		cli.ExitWithError(fmt.Sprintf("providerConfigId is required for mode %s", formattedMode), nil)
	}

	var publicKeyCtx *policy.KasPublicKeyCtx
	var privateKeyCtx *policy.KasPrivateKeyCtx
	if mode == policy.KeyMode_KEY_MODE_CONFIG_ROOT_KEY {
		wrappingKey := c.Flags.GetRequiredString("wrappingKey")
		privateKeyPem, publicKeyPem, err := generateKeys(alg)
		if err != nil {
			cli.ExitWithError("Failed to generate keys", err)
		}

		privateKey, err := wrapKey(privateKeyPem, wrappingKey)
		if err != nil {
			cli.ExitWithError("Failed to wrap key", err)
		}

		pubPemBase64 := base64.StdEncoding.EncodeToString([]byte(publicKeyPem))
		if err != nil {
			cli.ExitWithError("Failed to base64 encode public key", err)
		}

		privPemBase64 := base64.StdEncoding.EncodeToString(privateKey)
		if err != nil {
			cli.ExitWithError("Failed to base64 encode private key", err)
		}

		publicKeyCtx = &policy.KasPublicKeyCtx{
			Pem: pubPemBase64,
		}
		privateKeyCtx = &policy.KasPrivateKeyCtx{
			KeyId:      wrappingKeyId,
			WrappedKey: privPemBase64,
		}
	} else if mode == policy.KeyMode_KEY_MODE_PROVIDER_ROOT_KEY {
		providerConfigId = c.Flags.GetRequiredString("providerConfigId")
		publicPem := c.Flags.GetRequiredString("pubPem")
		privatePem := c.Flags.GetRequiredString("privatePem")
		_, err = base64.StdEncoding.DecodeString(publicPem)
		if err != nil {
			cli.ExitWithError("pem must be base64 encoded", err)
		}
		_, err = base64.StdEncoding.DecodeString(privatePem)
		if err != nil {
			cli.ExitWithError("pem must be base64 encoded", err)
		}
		publicKeyCtx = &policy.KasPublicKeyCtx{
			Pem: publicPem,
		}
		privateKeyCtx = &policy.KasPrivateKeyCtx{
			KeyId:      wrappingKeyId,
			WrappedKey: privatePem,
		}

	} else if mode == policy.KeyMode_KEY_MODE_REMOTE {
		pem := c.Flags.GetRequiredString("pubPem")
		providerConfigId = c.Flags.GetRequiredString("providerConfigId")

		_, err = base64.StdEncoding.DecodeString(pem)
		if err != nil {
			cli.ExitWithError("pem must be base64 encoded", err)
		}

		publicKeyCtx = &policy.KasPublicKeyCtx{
			Pem: pem,
		}
		privateKeyCtx = &policy.KasPrivateKeyCtx{
			KeyId: wrappingKeyId,
		}
	} else if mode == policy.KeyMode_KEY_MODE_PUBLIC_KEY_ONLY {
		pem := c.Flags.GetRequiredString("pubPem")
		_, err = base64.StdEncoding.DecodeString(pem)
		if err != nil {
			cli.ExitWithError("pem must be base64 encoded", err)
		}
		publicKeyCtx = &policy.KasPublicKeyCtx{
			Pem: pem,
		}
	}

	if kasId == "" {
		kas, err := h.GetKasRegistryEntry(c.Context(), handlers.KasIdentifier{
			Name: kasName,
			URI:  kasUri,
		})
		if err != nil {
			cli.ExitWithError("Failed to get kas registry entry", err)
		}
		kasId = kas.GetId()
	}

	kasKey, err := h.CreateKasKey(
		c.Context(),
		kasId,
		keyId,
		alg,
		mode,
		publicKeyCtx,
		privateKeyCtx,
		providerConfigId,
		getMetadataMutable(metadataLabels),
	)
	if err != nil {
		cli.ExitWithError("Failed to create kas key", err)
	}

	rows := getTableRows(kasKey)
	if mdRows := getMetadataRows(kasKey.GetKey().GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, kasKey.GetKey().GetId(), t, kasKey)
}

func getKasKeyIdentifier(c *cli.Cli) *kasregistry.KasKeyIdentifier {
	keyId := c.Flags.GetOptionalString("keyId")
	kasId := c.Flags.GetOptionalString("kasId")
	kasName := c.Flags.GetOptionalString("kasName")
	kasUri := c.Flags.GetOptionalString("kasUri")

	if keyId != "" {
		identifier := &kasregistry.KasKeyIdentifier{
			Kid: keyId,
		}

		if kasId != "" {
			identifier.Identifier = &kasregistry.KasKeyIdentifier_KasId{
				KasId: kasId,
			}
		} else if kasName != "" {
			identifier.Identifier = &kasregistry.KasKeyIdentifier_Name{
				Name: kasName,
			}
		} else if kasUri != "" {
			identifier.Identifier = &kasregistry.KasKeyIdentifier_Uri{
				Uri: kasUri,
			}
		} else {
			return nil
		}
		return identifier
	}

	return nil

}

func policy_getKASKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetOptionalID("id")

	kasKey, err := h.GetKasKey(c.Context(), id, getKasKeyIdentifier(c))
	if err != nil {
		cli.ExitWithError("Failed to get kas key", err)
	}

	rows := getTableRows(kasKey)
	if mdRows := getMetadataRows(kasKey.GetKey().GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, kasKey.GetKey().GetId(), t, kasKey)
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
	kasKey, err := h.GetKasKey(c.Context(), resp.GetKey().GetId(), nil)
	if err != nil {
		cli.ExitWithError("Failed to get kas key", err)
	}

	rows := getTableRows(kasKey)
	if mdRows := getMetadataRows(kasKey.GetKey().GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, kasKey.GetKey().GetId(), t, kasKey)
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
	for _, kasKey := range keys {
		key := kasKey.GetKey()
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

		pubCtxBytes, err := protojson.Marshal(key.GetPublicKeyCtx())
		if err != nil {
			cli.ExitWithError("Failed to marshal public key context", err)
		}
		privateKeyBytes, err := protojson.Marshal(key.GetPrivateKeyCtx())
		if err != nil {
			cli.ExitWithError("Failed to marshal private key context", err)
		}

		rows = append(rows, table.NewRow(table.RowData{
			"id":             key.GetId(),
			"keyId":          key.GetKeyId(),
			"keyAlgorithm":   algStr,
			"keyStatus":      statusStr,
			"keyMode":        modeStr,
			"pubKeyCtx":      string(pubCtxBytes),
			"privateKeyCtx":  string(privateKeyBytes),
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
		createDoc.GetDocFlag("kasId").Name,
		createDoc.GetDocFlag("kasId").Shorthand,
		createDoc.GetDocFlag("kasId").Default,
		createDoc.GetDocFlag("kasId").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("kasUri").Name,
		createDoc.GetDocFlag("kasUri").Shorthand,
		createDoc.GetDocFlag("kasUri").Default,
		createDoc.GetDocFlag("kasUri").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("kasName").Name,
		createDoc.GetDocFlag("kasName").Shorthand,
		createDoc.GetDocFlag("kasName").Default,
		createDoc.GetDocFlag("kasName").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("wrappingKeyId").Name,
		createDoc.GetDocFlag("wrappingKeyId").Shorthand,
		createDoc.GetDocFlag("wrappingKeyId").Default,
		createDoc.GetDocFlag("wrappingKeyId").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("wrappingKey").Name,
		createDoc.GetDocFlag("wrappingKey").Shorthand,
		createDoc.GetDocFlag("wrappingKey").Default,
		createDoc.GetDocFlag("wrappingKey").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("providerConfigId").Name,
		createDoc.GetDocFlag("providerConfigId").Shorthand,
		createDoc.GetDocFlag("providerConfigId").Default,
		createDoc.GetDocFlag("providerConfigId").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("pubPem").Name,
		createDoc.GetDocFlag("pubPem").Shorthand,
		createDoc.GetDocFlag("pubPem").Default,
		createDoc.GetDocFlag("pubPem").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("privatePem").Name,
		createDoc.GetDocFlag("privatePem").Shorthand,
		createDoc.GetDocFlag("privatePem").Default,
		createDoc.GetDocFlag("privatePem").Description,
	)
	createDoc.MarkFlagsOneRequired("kasId", "kasUri", "kasName")
	createDoc.MarkFlagsMutuallyExclusive("kasId", "kasUri", "kasName")
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
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("kasUri").Name,
		getDoc.GetDocFlag("kasUri").Shorthand,
		getDoc.GetDocFlag("kasUri").Default,
		getDoc.GetDocFlag("kasUri").Description,
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("kasId").Name,
		getDoc.GetDocFlag("kasId").Shorthand,
		getDoc.GetDocFlag("kasId").Default,
		getDoc.GetDocFlag("kasId").Description,
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("kasName").Name,
		getDoc.GetDocFlag("kasName").Shorthand,
		getDoc.GetDocFlag("kasName").Default,
		getDoc.GetDocFlag("kasName").Description,
	)
	getDoc.MarkFlagsMutuallyExclusive("id", "keyId")
	getDoc.MarkFlagsMutuallyExclusive("kasUri", "kasId", "kasName")

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
