package cmd

import (
	"encoding/base64"
	"encoding/hex"
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
	"github.com/opentdf/platform/protocol/go/policy/kasregistry"
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

func wrapKey(key string, wrappingKey []byte) ([]byte, error) {
	aesKey, err := ocrypto.NewAESGcm(wrappingKey)
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
	var wrappingKeyID string

	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	keyIdentifier := c.Flags.GetRequiredString("key-id")
	kasIdentifier := c.Flags.GetRequiredString("kas")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	// Use the helper function to get and validate key parameters
	alg, mode, wrappingKeyID, err := prepareKeyParams(c)
	if err != nil {
		cli.ExitWithError("Invalid key parameters", err)
	}

	// Use the helper function to prepare key contexts
	publicKeyCtx, privateKeyCtx, providerConfigID, err := prepareKeyContexts(c, mode, alg, wrappingKeyID)
	if err != nil {
		cli.ExitWithError("Failed to prepare key contexts", err)
	}

	kasLookup, err := resolveKasIdentifier(kasIdentifier)
	if err != nil {
		cli.ExitWithError("Invalid kas identifier", err)
	}

	var resolvedKasID string
	if kasLookup.ID != "" {
		resolvedKasID = kasLookup.ID
	} else {
		// If not a UUID, resolve it to get the UUID
		kasEntry, err := h.GetKasRegistryEntry(c.Context(), kasLookup)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to resolve KAS identifier '%s'", kasIdentifier), err)
		}
		resolvedKasID = kasEntry.GetId()
	}

	kasKey, err := h.CreateKasKey(
		c.Context(),
		resolvedKasID,
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

	id := c.Flags.GetOptionalString("key")

	var identifier *kasregistry.KasKeyIdentifier
	var err error

	if utils.ClassifyString(id) != utils.StringTypeUUID {
		identifier, err = getKasKeyIdentifier(c)
		if err != nil {
			cli.ExitWithError("Invalid key identifier", err)
		}
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
	kasIdentifier := c.Flags.GetOptionalString("kas")

	kasLookup, err := resolveKasIdentifier(kasIdentifier)
	if err != nil {
		cli.ExitWithError("Invalid kas identifier", err)
	}

	// Get the list of keys.
	keys, page, err := h.ListKasKeys(c.Context(), limit, offset, alg, kasLookup)
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

func policyRotateKasKey(cmd *cobra.Command, args []string) {
	var wrappingKeyID string

	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	// Get parameters for the old key
	oldKey := c.Flags.GetRequiredString("key")

	// Get parameters for creating the new key
	newKeyID := c.Flags.GetRequiredString("key-id")
	metadataLabels = c.Flags.GetStringSlice("label", metadataLabels, cli.FlagsStringSliceOptions{Min: 0})

	// Use the helper function to get and validate key parameters
	alg, mode, wrappingKeyID, err := prepareKeyParams(c)
	if err != nil {
		cli.ExitWithError("Invalid key parameters", err)
	}

	// Use the helper function to prepare key contexts
	publicKeyCtx, privateKeyCtx, providerConfigID, err := prepareKeyContexts(c, mode, alg, wrappingKeyID)
	if err != nil {
		cli.ExitWithError("Failed to prepare key contexts", err)
	}

	// Create the new key request with the contexts created by the helper
	newKey := &kasregistry.RotateKeyRequest_NewKey{
		KeyId:            newKeyID,
		Algorithm:        alg,
		KeyMode:          mode,
		PublicKeyCtx:     publicKeyCtx,
		PrivateKeyCtx:    privateKeyCtx,
		ProviderConfigId: providerConfigID,
		Metadata:         getMetadataMutable(metadataLabels),
	}

	var identifier *kasregistry.KasKeyIdentifier
	if utils.ClassifyString(oldKey) != utils.StringTypeUUID {
		identifier, err = getKasKeyIdentifier(c)
		if err != nil {
			cli.ExitWithError("Invalid key identifier", err)
		}
	}

	// Call the rotate key function
	rotateKeyResult, err := h.RotateKasKey(
		c.Context(),
		oldKey,
		identifier,
		newKey,
	)
	if err != nil {
		cli.ExitWithError("Failed to rotate key", err)
	}

	rows := getTableRows(rotateKeyResult.KasKey)
	if mdRows := getMetadataRows(rotateKeyResult.KasKey.GetKey().GetMetadata()); mdRows != nil {
		rows = append(rows, mdRows...)
	}
	t := cli.NewTabular(rows...)
	HandleSuccess(cmd, rotateKeyResult.KasKey.GetKey().GetId(), t, rotateKeyResult)
}

func resolveKasIdentifier(ident string) (handlers.KasIdentifier, error) {
	// If the identifier is empty, it means no KAS filter is applied.
	// Return an empty KasIdentifier and no error.
	if ident == "" {
		return handlers.KasIdentifier{}, nil
	}

	// Use the ClassifyString helper to determine how to look up the KAS
	kasLookup := handlers.KasIdentifier{}
	kasInputType := utils.ClassifyString(ident)

	switch kasInputType { //nolint:exhaustive // default catches unknown
	case utils.StringTypeUUID:
		kasLookup.ID = ident
	case utils.StringTypeURI:
		kasLookup.URI = ident
	case utils.StringTypeGeneric:
		kasLookup.Name = ident
	default:
		return kasLookup, errors.New("invalid kas identifier")
	}

	return kasLookup, nil
}

// prepareKeyParams parses and validates the common key parameters used by both create and rotate operations.
// It returns the algorithm, mode, wrapping key ID, and any error that occurred.
func prepareKeyParams(c *cli.Cli) (policy.Algorithm, policy.KeyMode, string, error) {
	// Parse algorithm
	alg, err := algToEnum(c.Flags.GetRequiredString("algorithm"))
	if err != nil {
		return alg, 0, "", err
	}

	// Parse mode
	mode, err := modeToEnum(c.Flags.GetRequiredString("mode"))
	if err != nil {
		return alg, mode, "", err
	}

	// Get wrapping key ID and validate based on mode
	wrappingKeyID := c.Flags.GetOptionalString("wrapping-key-id")
	if mode != policy.KeyMode_KEY_MODE_PUBLIC_KEY_ONLY && wrappingKeyID == "" {
		formattedMode, _ := enumToMode(mode)
		return alg, mode, "", fmt.Errorf("wrapping-key-id is required for mode %s", formattedMode)
	}

	return alg, mode, wrappingKeyID, nil
}

// prepareKeyContexts prepares the key contexts based on the specified mode and parameters.
// This function encapsulates the common logic between key creation and key rotation.
func prepareKeyContexts(
	c *cli.Cli,
	mode policy.KeyMode,
	alg policy.Algorithm,
	wrappingKeyID string,
) (*policy.PublicKeyCtx, *policy.PrivateKeyCtx, string, error) {
	var publicKeyCtx *policy.PublicKeyCtx
	var privateKeyCtx *policy.PrivateKeyCtx
	var providerConfigID string

	switch mode {
	case policy.KeyMode_KEY_MODE_CONFIG_ROOT_KEY:
		// Local mode: generate keys locally and wrap with provided wrapping key
		wrappingKey := c.Flags.GetRequiredString("wrapping-key")
		wrappedKeyBytes, err := hex.DecodeString(wrappingKey)
		if err != nil {
			return nil, nil, "", errors.Join(errors.New("wrapping-key must be hex encoded"), err)
		}

		privateKeyPem, publicKeyPem, err := generateKeys(alg)
		if err != nil {
			return nil, nil, "", errors.Join(errors.New("failed to generate keys"), err)
		}

		privateKey, err := wrapKey(privateKeyPem, wrappedKeyBytes)
		if err != nil {
			return nil, nil, "", errors.Join(errors.New("failed to wrap key"), err)
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
		_, err := base64.StdEncoding.DecodeString(publicPem)
		if err != nil {
			return nil, nil, "", errors.Join(errors.New("public key pem must be base64 encoded"), err)
		}
		_, err = base64.StdEncoding.DecodeString(privatePem)
		if err != nil {
			return nil, nil, "", errors.Join(errors.New("private key pem must be base64 encoded"), err)
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

		_, err := base64.StdEncoding.DecodeString(pem)
		if err != nil {
			return nil, nil, "", errors.Join(errors.New("pem must be base64 encoded"), err)
		}

		publicKeyCtx = &policy.PublicKeyCtx{
			Pem: pem,
		}
		privateKeyCtx = &policy.PrivateKeyCtx{
			KeyId: wrappingKeyID,
		}
	case policy.KeyMode_KEY_MODE_PUBLIC_KEY_ONLY:
		pem := c.Flags.GetRequiredString("public-key-pem")
		_, err := base64.StdEncoding.DecodeString(pem)
		if err != nil {
			return nil, nil, "", errors.Join(errors.New("pem must be base64 encoded"), err)
		}
		publicKeyCtx = &policy.PublicKeyCtx{
			Pem: pem,
		}
	case policy.KeyMode_KEY_MODE_UNSPECIFIED:
		fallthrough
	default:
		return nil, nil, "", errors.New("invalid mode")
	}

	return publicKeyCtx, privateKeyCtx, providerConfigID, nil
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

	// Rotate Kas Key
	rotateDoc := man.Docs.GetCommand("policy/kas-registry/key/rotate",
		man.WithRun(policyRotateKasKey),
	)
	rotateDoc.Flags().StringP(
		rotateDoc.GetDocFlag("key").Name,
		rotateDoc.GetDocFlag("key").Shorthand,
		rotateDoc.GetDocFlag("key").Default,
		rotateDoc.GetDocFlag("key").Description,
	)
	rotateDoc.Flags().StringP(
		rotateDoc.GetDocFlag("kas").Name,
		rotateDoc.GetDocFlag("kas").Shorthand,
		rotateDoc.GetDocFlag("kas").Default,
		rotateDoc.GetDocFlag("kas").Description,
	)
	rotateDoc.Flags().StringP(
		rotateDoc.GetDocFlag("key-id").Name,
		rotateDoc.GetDocFlag("key-id").Shorthand,
		rotateDoc.GetDocFlag("key-id").Default,
		rotateDoc.GetDocFlag("key-id").Description,
	)
	rotateDoc.Flags().StringP(
		rotateDoc.GetDocFlag("algorithm").Name,
		rotateDoc.GetDocFlag("algorithm").Shorthand,
		rotateDoc.GetDocFlag("algorithm").Default,
		rotateDoc.GetDocFlag("algorithm").Description,
	)
	rotateDoc.Flags().StringP(
		rotateDoc.GetDocFlag("mode").Name,
		rotateDoc.GetDocFlag("mode").Shorthand,
		rotateDoc.GetDocFlag("mode").Default,
		rotateDoc.GetDocFlag("mode").Description,
	)
	rotateDoc.Flags().StringP(
		rotateDoc.GetDocFlag("wrapping-key-id").Name,
		rotateDoc.GetDocFlag("wrapping-key-id").Shorthand,
		rotateDoc.GetDocFlag("wrapping-key-id").Default,
		rotateDoc.GetDocFlag("wrapping-key-id").Description,
	)
	rotateDoc.Flags().StringP(
		rotateDoc.GetDocFlag("wrapping-key").Name,
		rotateDoc.GetDocFlag("wrapping-key").Shorthand,
		rotateDoc.GetDocFlag("wrapping-key").Default,
		rotateDoc.GetDocFlag("wrapping-key").Description,
	)
	rotateDoc.Flags().StringP(
		rotateDoc.GetDocFlag("provider-config-id").Name,
		rotateDoc.GetDocFlag("provider-config-id").Shorthand,
		rotateDoc.GetDocFlag("provider-config-id").Default,
		rotateDoc.GetDocFlag("provider-config-id").Description,
	)
	rotateDoc.Flags().StringP(
		rotateDoc.GetDocFlag("public-key-pem").Name,
		rotateDoc.GetDocFlag("public-key-pem").Shorthand,
		rotateDoc.GetDocFlag("public-key-pem").Default,
		rotateDoc.GetDocFlag("public-key-pem").Description,
	)
	rotateDoc.Flags().StringP(
		rotateDoc.GetDocFlag("private-key-pem").Name,
		rotateDoc.GetDocFlag("private-key-pem").Shorthand,
		rotateDoc.GetDocFlag("private-key-pem").Default,
		rotateDoc.GetDocFlag("private-key-pem").Description,
	)
	injectLabelFlags(&rotateDoc.Command, true)

	policyKasRegistryKeysCmd.AddSubcommands(createDoc, getDoc, updateDoc, listDoc, rotateDoc)
	policyKasRegCmd.AddCommand(&policyKasRegistryKeysCmd.Command)
}
