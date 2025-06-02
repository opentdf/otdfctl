package cmd

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/lib/ocrypto"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/spf13/cobra"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	defaultAlg           = 0
	defaultMode          = 0
	defaultStatus        = 0
	rsa2048Len           = 2048
	rsa4096Len           = 4096
	ecSecp256Len         = 256
	ecSecp384Len         = 384
	ecSecp521Len         = 521
	keyStatusActive      = "active"
	keyStatusRotated     = "rotated"
	keyModeLocal         = "local"
	keyModeProvider      = "provider"
	keyModeRemote        = "remote"
	keyModePublicKeyOnly = "public_key"
)

var policyKasRegistryKeysCmd = man.Docs.GetCommand("policy/kas-registry/key")

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
		key, err = generateRSAKey(rsa2048Len)
	case policy.Algorithm_ALGORITHM_RSA_4096:
		key, err = generateRSAKey(rsa4096Len)
	case policy.Algorithm_ALGORITHM_EC_P256:
		key, err = generateECCKey(ecSecp256Len)
	case policy.Algorithm_ALGORITHM_EC_P384:
		key, err = generateECCKey(ecSecp384Len)
	case policy.Algorithm_ALGORITHM_EC_P521:
		key, err = generateECCKey(ecSecp521Len)
	case policy.Algorithm_ALGORITHM_UNSPECIFIED:
		fallthrough
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
	case policy.KeyStatus_KEY_STATUS_ROTATED:
		return keyStatusRotated, nil
	default:
		return "", errors.New("invalid enum status")
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
func policyCreateKasKey(cmd *cobra.Command, args []string) {
	var (
		wrappingKeyID    string
		providerConfigID string
	)

	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	keyIdentifier := c.Flags.GetRequiredString("key-id")
	kasIdentifier := c.Flags.GetRequiredString("kas")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	alg, err := algToEnum(c.Flags.GetRequiredString("algorithm"))
	if err != nil {
		cli.ExitWithError("Invalid algorithm", err)
	}

	mode, err := modeToEnum(c.Flags.GetRequiredString("mode"))
	if err != nil {
		cli.ExitWithError("Invalid mode", err)
	}

	wrappingKeyID = c.Flags.GetOptionalString("wrapping-key-id")
	if mode != policy.KeyMode_KEY_MODE_PUBLIC_KEY_ONLY && wrappingKeyID == "" {
		formattedMode, _ := enumToMode(mode)
		cli.ExitWithError(fmt.Sprintf("wrapping-key-id is required for mode %s", formattedMode), nil)
	}

	providerConfigID = c.Flags.GetOptionalString("provider-config-id")
	if (mode == policy.KeyMode_KEY_MODE_PROVIDER_ROOT_KEY || mode == policy.KeyMode_KEY_MODE_REMOTE) && providerConfigID == "" {
		formattedMode, _ := enumToMode(mode)
		cli.ExitWithError(fmt.Sprintf("provider-config-id is required for mode %s", formattedMode), nil)
	}

	var publicKeyCtx *policy.PublicKeyCtx
	var privateKeyCtx *policy.PrivateKeyCtx
	switch mode {
	case policy.KeyMode_KEY_MODE_CONFIG_ROOT_KEY:
		wrappingKey := c.Flags.GetRequiredString("wrapping-key")
		privateKeyPem, publicKeyPem, err := generateKeys(alg)
		if err != nil {
			cli.ExitWithError("Failed to generate keys", err)
		}

		privateKey, err := wrapKey(privateKeyPem, wrappingKey)
		if err != nil {
			cli.ExitWithError("Failed to wrap key", err)
		}

		pubPemBase64 := base64.StdEncoding.EncodeToString([]byte(publicKeyPem))
		privPemBase64 := base64.StdEncoding.EncodeToString(privateKey)
		publicKeyCtx = &policy.PublicKeyCtx{
			Pem: pubPemBase64,
		}
		privateKeyCtx = &policy.PrivateKeyCtx{
			KeyId:      wrappingKeyID,
			WrappedKey: privPemBase64,
		}
	case policy.KeyMode_KEY_MODE_PROVIDER_ROOT_KEY:
		providerConfigID = c.Flags.GetRequiredString("provider-config-id")
		publicPem := c.Flags.GetRequiredString("public-key-pem")
		privatePem := c.Flags.GetRequiredString("private-key-pem")
		_, err = base64.StdEncoding.DecodeString(publicPem)
		if err != nil {
			cli.ExitWithError("pem must be base64 encoded", err)
		}
		_, err = base64.StdEncoding.DecodeString(privatePem)
		if err != nil {
			cli.ExitWithError("pem must be base64 encoded", err)
		}
		publicKeyCtx = &policy.PublicKeyCtx{
			Pem: publicPem,
		}
		privateKeyCtx = &policy.PrivateKeyCtx{
			KeyId:      wrappingKeyID,
			WrappedKey: privatePem,
		}
	case policy.KeyMode_KEY_MODE_REMOTE:
		pem := c.Flags.GetRequiredString("public-key-pem")
		providerConfigID = c.Flags.GetRequiredString("provider-config-id")

		_, err = base64.StdEncoding.DecodeString(pem)
		if err != nil {
			cli.ExitWithError("pem must be base64 encoded", err)
		}

		publicKeyCtx = &policy.PublicKeyCtx{
			Pem: pem,
		}
		privateKeyCtx = &policy.PrivateKeyCtx{
			KeyId: wrappingKeyID,
		}
	case policy.KeyMode_KEY_MODE_PUBLIC_KEY_ONLY:
		pem := c.Flags.GetRequiredString("public-key-pem")
		_, err = base64.StdEncoding.DecodeString(pem)
		if err != nil {
			cli.ExitWithError("pem must be base64 encoded", err)
		}
		publicKeyCtx = &policy.PublicKeyCtx{
			Pem: pem,
		}
	case policy.KeyMode_KEY_MODE_UNSPECIFIED:
		fallthrough
	default:
		cli.ExitWithError("Invalid mode", nil)
	}

	kasIdentifier, err = resolveKasIdentifier(c.Context(), kasIdentifier, h)
	if err != nil {
		cli.ExitWithError("Invalid kas identifier", err)
	}

	kasKey, err := h.CreateKasKey(
		c.Context(),
		kasIdentifier,
		keyIdentifier,
		alg,
		mode,
		publicKeyCtx,
		privateKeyCtx,
		providerConfigID,
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

func policyGetKasKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetOptionalID("id")

	identifier, err := getKasKeyIdentifier(c)
	if err != nil {
		cli.ExitWithError("Invalid key identifier", err)
	}
	kasKey, err := h.GetKasKey(c.Context(), id, identifier)
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

func policyUpdateKasKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetRequiredID("id")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	resp, err := h.UpdateKasKey(
		c.Context(),
		id,
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

func policyListKasKeys(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	limit := c.Flags.GetRequiredInt32("limit")
	offset := c.Flags.GetRequiredInt32("offset")
	algArg := c.Flags.GetOptionalString("algorithm")
	var alg policy.Algorithm
	if algArg != "" {
		var err error
		alg, err = algToEnum(algArg)
		if err != nil {
			cli.ExitWithError("Invalid algorithm", err)
		}
	}
	kasIdentifier := c.Flags.GetRequiredString("kas")

	kasIdentifier, err := resolveKasIdentifier(c.Context(), kasIdentifier, h)
	if err != nil {
		cli.ExitWithError("Invalid kas identifier", err)
	}

	// Get the list of keys.
	keys, page, err := h.ListKasKeys(c.Context(), limit, offset, alg, handlers.KasIdentifier{
		ID: kasIdentifier,
	})
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

func resolveKasIdentifier(ctx context.Context, ident string, h handlers.Handler) (string, error) {
	// Use the ClassifyString helper to determine how to look up the KAS
	kasLookup := handlers.KasIdentifier{}
	kasInputType := utils.ClassifyString(ident)

	switch kasInputType { //nolint:exhaustive // default catches unknown
	case utils.StringTypeUUID:
		return ident, nil
	case utils.StringTypeURI:
		kasLookup.URI = ident
	case utils.StringTypeGeneric:
		kasLookup.Name = ident
	default:
		return "", errors.New("invalid kas identifier")
	}

	if kasInputType != utils.StringTypeUUID {
		resolvedKas, err := h.GetKasRegistryEntry(ctx, kasLookup)
		if err != nil {
			return "", errors.Join(errors.New("failed to get kas registry entry"), err)
		}
		return resolvedKas.GetId(), nil
	}
	return "", nil
}

func init() {
	// Create Kas Key
	createDoc := man.Docs.GetCommand("policy/kas-registry/key/create",
		man.WithRun(policyCreateKasKey),
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("key-id").Name,
		createDoc.GetDocFlag("key-id").Shorthand,
		createDoc.GetDocFlag("key-id").Default,
		createDoc.GetDocFlag("key-id").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("algorithm").Name,
		createDoc.GetDocFlag("algorithm").Shorthand,
		createDoc.GetDocFlag("algorithm").Default,
		createDoc.GetDocFlag("algorithm").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("mode").Name,
		createDoc.GetDocFlag("mode").Shorthand,
		createDoc.GetDocFlag("mode").Default,
		createDoc.GetDocFlag("mode").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("kas").Name,
		createDoc.GetDocFlag("kas").Shorthand,
		createDoc.GetDocFlag("kas").Default,
		createDoc.GetDocFlag("kas").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("wrapping-key-id").Name,
		createDoc.GetDocFlag("wrapping-key-id").Shorthand,
		createDoc.GetDocFlag("wrapping-key-id").Default,
		createDoc.GetDocFlag("wrapping-key-id").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("wrapping-key").Name,
		createDoc.GetDocFlag("wrapping-key").Shorthand,
		createDoc.GetDocFlag("wrapping-key").Default,
		createDoc.GetDocFlag("wrapping-key").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("provider-config-id").Name,
		createDoc.GetDocFlag("provider-config-id").Shorthand,
		createDoc.GetDocFlag("provider-config-id").Default,
		createDoc.GetDocFlag("provider-config-id").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("public-key-pem").Name,
		createDoc.GetDocFlag("public-key-pem").Shorthand,
		createDoc.GetDocFlag("public-key-pem").Default,
		createDoc.GetDocFlag("public-key-pem").Description,
	)
	createDoc.Flags().StringP(
		createDoc.GetDocFlag("private-key-pem").Name,
		createDoc.GetDocFlag("private-key-pem").Shorthand,
		createDoc.GetDocFlag("private-key-pem").Default,
		createDoc.GetDocFlag("private-key-pem").Description,
	)
	injectLabelFlags(&createDoc.Command, false)

	// Get Kas Key
	getDoc := man.Docs.GetCommand("policy/kas-registry/key/get",
		man.WithRun(policyGetKasKey),
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("key").Name,
		getDoc.GetDocFlag("key").Shorthand,
		getDoc.GetDocFlag("key").Default,
		getDoc.GetDocFlag("key").Description,
	)
	getDoc.Flags().StringP(
		getDoc.GetDocFlag("kas").Name,
		getDoc.GetDocFlag("kas").Shorthand,
		getDoc.GetDocFlag("kas").Default,
		getDoc.GetDocFlag("kas").Description,
	)
	// Update Kas Key
	updateDoc := man.Docs.GetCommand("policy/kas-registry/key/update",
		man.WithRun(policyUpdateKasKey),
	)
	updateDoc.Flags().StringP(
		updateDoc.GetDocFlag("id").Name,
		updateDoc.GetDocFlag("id").Shorthand,
		updateDoc.GetDocFlag("id").Default,
		updateDoc.GetDocFlag("id").Description,
	)
	injectLabelFlags(&updateDoc.Command, true)

	// List Kas Keys
	listDoc := man.Docs.GetCommand("policy/kas-registry/key/list",
		man.WithRun(policyListKasKeys),
	)
	listDoc.Flags().StringP(
		listDoc.GetDocFlag("algorithm").Name,
		listDoc.GetDocFlag("algorithm").Shorthand,
		listDoc.GetDocFlag("algorithm").Default,
		listDoc.GetDocFlag("algorithm").Description,
	)
	listDoc.Flags().StringP(
		listDoc.GetDocFlag("kas").Name,
		listDoc.GetDocFlag("kas").Shorthand,
		listDoc.GetDocFlag("kas").Default,
		listDoc.GetDocFlag("kas").Description,
	)
	injectListPaginationFlags(listDoc)

	policyKasRegistryKeysCmd.AddSubcommands(createDoc, getDoc, updateDoc, listDoc)
	policyKasRegCmd.AddCommand(&policyKasRegistryKeysCmd.Command)
}
