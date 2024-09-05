package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/sdk"
	"github.com/spf13/cobra"
)

type tdfInspectManifest struct {
	Algorithm             string                    `json:"algorithm"`
	KeyAccessType         string                    `json:"keyAccessType"`
	MimeType              string                    `json:"mimeType"`
	Policy                string                    `json:"policy"`
	Protocol              string                    `json:"protocol"`
	SegmentHashAlgorithm  string                    `json:"segmentHashAlgorithm"`
	Signature             string                    `json:"signature"`
	Type                  string                    `json:"type"`
	Method                sdk.Method                `json:"method"`
	IntegrityInformation  sdk.IntegrityInformation  `json:"integrityInformation"`
	EncryptionInformation sdk.EncryptionInformation `json:"encryptionInformation"`
}

type nanoInspectResult struct {
	Cipher       string `json:"cipher"`
	ECDSAEnabled bool   `json:"ecdsaEnabled"`
	Kas          string `json:"kas"`
	KID          string `json:"kid"`
}

type tdfInspectResult struct {
	Manifest   tdfInspectManifest `json:"manifest"`
	Attributes []string           `json:"attributes"`
}

func tdf_InspectCmd(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args, cli.WithPrintJson())
	h := NewHandler(c)
	defer h.Close()

	data := cli.ReadFromArgsOrPipe(args, nil)
	if len(data) == 0 {
		c.ExitWithError("must provide ONE of the following: [file argument, stdin input]", errors.New("no input provided"))
	}

	result, errs := h.InspectTDF(data)
	for _, err := range errs {
		if errors.Is(err, handlers.ErrTDFInspectFailNotValidTDF) {
			c.ExitWithError("not a valid TDF", err)
		} else if errors.Is(err, handlers.ErrTDFInspectFailNotInspectable) {
			c.ExitWithError("failed to inspect TDF", err)
		}
	}

	if result.ZTDFManifest != nil {
		m := tdfInspectResult{
			Manifest: tdfInspectManifest{
				Algorithm:             result.ZTDFManifest.Algorithm,
				KeyAccessType:         result.ZTDFManifest.KeyAccessType,
				MimeType:              result.ZTDFManifest.MimeType,
				Policy:                result.ZTDFManifest.Policy,
				Protocol:              result.ZTDFManifest.Protocol,
				SegmentHashAlgorithm:  result.ZTDFManifest.SegmentHashAlgorithm,
				Signature:             result.ZTDFManifest.Signature,
				Type:                  result.ZTDFManifest.Type,
				Method:                result.ZTDFManifest.Method,
				IntegrityInformation:  result.ZTDFManifest.IntegrityInformation,
				EncryptionInformation: result.ZTDFManifest.EncryptionInformation,
			},
			Attributes: result.Attributes,
		}

		c.PrintJson(m)
	} else if result.NanoHeader != nil {
		kas, err := result.NanoHeader.GetKasURL().GetURL()
		if err != nil {
			c.ExitWithError("not a valid NanoTDF", err)
		}
		kid, err := result.NanoHeader.GetKasURL().GetIdentifier()
		if err != nil {
			c.ExitWithError("not a valid NanoTDF", err)
		}
		cipher := result.NanoHeader.GetCipher()
		cipherBytes, err := sdk.SizeOfAuthTagForCipher(cipher)
		if err != nil {
			c.ExitWithError("not a valid NanoTDF", err)
		}
		cipherName := fmt.Sprintf("AES-%d", 8*cipherBytes)

		n := nanoInspectResult{
			Kas:          kas,
			KID:          strings.TrimRight(kid, "\u0000"),
			ECDSAEnabled: result.NanoHeader.IsEcdsaBindingEnabled(),
			Cipher:       cipherName,
		}

		c.PrintJson(n)
	} else {
		c.ExitWithError("failed to inspect TDF", nil)
	}
}

func init() {
	tdf_InspectCmd := man.Docs.GetCommand("inspect",
		man.WithRun(tdf_InspectCmd),
	)
	tdf_InspectCmd.Command.GroupID = TDF

	tdf_InspectCmd.Command.PreRun = func(cmd *cobra.Command, args []string) {
		// Set the json flag to true since we only support json output
		cmd.SetArgs(append(args, "--json"))
	}

	RootCmd.AddCommand(&tdf_InspectCmd.Command)
}
