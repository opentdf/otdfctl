package cmd

import (
	"errors"
	"fmt"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/protocol/go/policy/kasregistry"
	"github.com/spf13/cobra"
)

const (
	kasURIKey       = "kas_uri"
	kasURIColumn    = "Kas URI"
	algKey          = "algorithm"
	algColumn       = "Algorithm"
	pubPemKey       = "public_key_pem"
	pubPemColumn    = "Public Key PEM"
	kasKidKey       = "kas_key_id"
	kasKidColumn    = "Key ID"
	isBaseKey       = "is_base_key"
	isBaseKeyColumn = "Is Base Key"
)

// KAS Registry Base Keys Command
var policyKasRegistryBaseKeysCmd *cobra.Command

func getKasKeyIdentifier(c *cli.Cli) (*kasregistry.KasKeyIdentifier, error) {
	// This function is called when the user provides a human-readable key ID
	// via the --key flag and identifies the KAS it belongs to via the --kas flag.
	humanReadableKeyID := c.Flags.GetRequiredString("key") // Flag for the key's human-readable ID
	kasIdentifierInput := c.Flags.GetRequiredString("kas") // Flag for KAS ID, URI, or Name

	// Basic validation, though GetRequiredString should handle empty inputs.
	if humanReadableKeyID == "" {
		return nil, errors.New("--key (human-readable key ID) cannot be empty")
	}
	if kasIdentifierInput == "" {
		return nil, errors.New("--kas (KAS identifier) cannot be empty")
	}

	identifier := &kasregistry.KasKeyIdentifier{
		Kid: humanReadableKeyID,
	}

	kasInputType := utils.ClassifyString(kasIdentifierInput)
	switch kasInputType {
	case utils.StringTypeUUID:
		identifier.Identifier = &kasregistry.KasKeyIdentifier_KasId{KasId: kasIdentifierInput}
	case utils.StringTypeURI:
		identifier.Identifier = &kasregistry.KasKeyIdentifier_Uri{Uri: kasIdentifierInput}
	case utils.StringTypeGeneric:
		identifier.Identifier = &kasregistry.KasKeyIdentifier_Name{Name: kasIdentifierInput}
	default: // Catches StringTypeUnknown and any other unexpected types
		return nil, fmt.Errorf("invalid KAS identifier: '%s'. Must be a KAS UUID, URI, or Name", kasIdentifierInput)
	}
	return identifier, nil
}

func getBaseKeyTableRows(simpleKey *kasregistry.SimpleKasKey, additionalInfo map[string]string) table.Row {
	rowData := table.RowData{
		kasKidKey: simpleKey.GetPublicKey().GetKid(),
		pubPemKey: simpleKey.GetPublicKey().GetPem(),
		algKey:    simpleKey.GetPublicKey().GetAlgorithm(),
		kasURIKey: simpleKey.GetKasUri(),
	}

	if len(additionalInfo) > 0 {
		for key, value := range additionalInfo {
			rowData[key] = value
		}
	}

	return table.NewRow(rowData)
}

func getBaseKeyTable(additionalColumns []table.Column) table.Model {
	columns := []table.Column{
		table.NewFlexColumn(kasURIKey, kasURIColumn, cli.FlexColumnWidthOne),
		table.NewFlexColumn(kasKidKey, kasKidColumn, cli.FlexColumnWidthOne),
		table.NewFlexColumn(pubPemKey, pubPemColumn, cli.FlexColumnWidthOne),
		table.NewFlexColumn(algKey, algColumn, cli.FlexColumnWidthOne),
	}
	columns = append(columns, additionalColumns...)

	return cli.NewTable(
		columns...,
	)
}

func getBaseKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	baseKey, err := h.GetBaseKey(c.Context())
	if err != nil {
		cli.ExitWithError("Failed to get base key", err)
	}

	if baseKey == nil {
		cli.ExitWithError("No base key found", nil)
	}

	t := getBaseKeyTable(nil)
	t = t.WithRows([]table.Row{getBaseKeyTableRows(baseKey, nil)})
	HandleSuccess(cmd, "", t, baseKey)
}

func setBaseKey(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	id := c.Flags.GetOptionalID("id")

	identifier, err := getKasKeyIdentifier(c)
	if err != nil {
		c.ExitWithError("Invalid key identifier", err)
	}
	baseKey, err := h.SetBaseKey(c.Context(), id, identifier)
	if err != nil {
		cli.ExitWithError("Failed to set base key", err)
	}

	t := getBaseKeyTable([]table.Column{
		table.NewFlexColumn(isBaseKey, isBaseKeyColumn, cli.FlexColumnWidthOne),
	})

	rows := []table.Row{
		getBaseKeyTableRows(baseKey.GetNewBaseKey(), map[string]string{
			isBaseKey: "true",
		}),
	}
	if baseKey.GetPreviousBaseKey() != nil {
		rows = append(rows, getBaseKeyTableRows(baseKey.GetPreviousBaseKey(), map[string]string{
			isBaseKey: "false",
		}))
	}

	t = t.WithRows(rows)
	HandleSuccess(cmd, "", t, baseKey)
}

func init() {
	getDoc := man.Docs.GetCommand("policy/kas-registry/key/base/get",
		man.WithRun(getBaseKey),
	)

	setDoc := man.Docs.GetCommand("policy/kas-registry/key/base/set",
		man.WithRun(setBaseKey),
	)
	setDoc.Flags().StringP(
		setDoc.GetDocFlag("id").Name,
		setDoc.GetDocFlag("id").Shorthand,
		setDoc.GetDocFlag("id").Default,
		setDoc.GetDocFlag("id").Description,
	)
	setDoc.Flags().StringP(
		setDoc.GetDocFlag("keyId").Name,
		setDoc.GetDocFlag("keyId").Shorthand,
		setDoc.GetDocFlag("keyId").Default,
		setDoc.GetDocFlag("keyId").Description,
	)
	setDoc.Flags().StringP(
		setDoc.GetDocFlag("kasUri").Name,
		setDoc.GetDocFlag("kasUri").Shorthand,
		setDoc.GetDocFlag("kasUri").Default,
		setDoc.GetDocFlag("kasUri").Description,
	)
	setDoc.Flags().StringP(
		setDoc.GetDocFlag("kasId").Name,
		setDoc.GetDocFlag("kasId").Shorthand,
		setDoc.GetDocFlag("kasId").Default,
		setDoc.GetDocFlag("kasId").Description,
	)
	setDoc.Flags().StringP(
		setDoc.GetDocFlag("kasName").Name,
		setDoc.GetDocFlag("kasName").Shorthand,
		setDoc.GetDocFlag("kasName").Default,
		setDoc.GetDocFlag("kasName").Description,
	)
	setDoc.MarkFlagsMutuallyExclusive("id", "keyId")
	setDoc.MarkFlagsOneRequired("id", "keyId")
	setDoc.MarkFlagsMutuallyExclusive("kasUri", "kasId", "kasName")

	doc := man.Docs.GetCommand("policy/kas-registry/key/base",
		man.WithSubcommands(getDoc, setDoc))
	policyKasRegistryBaseKeysCmd = &doc.Command
	policyKasRegistryKeysCmd.AddCommand(
		policyKasRegistryBaseKeysCmd,
	)
}
