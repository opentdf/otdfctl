package cmd

import (
	"errors"

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
			c.ExitWithError("not a valid ZTDF", err)
		} else if errors.Is(err, handlers.ErrTDFInspectFailNotInspectable) {
			c.ExitWithError("failed to inspect TDF", err)
		}
	}

	m := tdfInspectResult{
		Manifest: tdfInspectManifest{
			Algorithm:             result.Manifest.Algorithm,
			KeyAccessType:         result.Manifest.KeyAccessType,
			MimeType:              result.Manifest.MimeType,
			Policy:                result.Manifest.Policy,
			Protocol:              result.Manifest.Protocol,
			SegmentHashAlgorithm:  result.Manifest.SegmentHashAlgorithm,
			Signature:             result.Manifest.Signature,
			Type:                  result.Manifest.Type,
			Method:                result.Manifest.Method,
			IntegrityInformation:  result.Manifest.IntegrityInformation,
			EncryptionInformation: result.Manifest.EncryptionInformation,
		},
		Attributes: result.Attributes,
	}

	c.PrintJson(m)
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
