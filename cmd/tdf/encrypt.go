package tdf

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/opentdf/otdfctl/cmd/common"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/lib/ocrypto"
	"github.com/spf13/cobra"
)

var attrValues []string
var assertions string

func encryptRun(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args, cli.WithPrintJSON())
	h := common.NewHandler(c)
	defer h.Close()

	var filePath string
	var fileExt string
	if len(args) > 0 {
		filePath = args[0]
		fileExt = strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))
	}

	out := c.Flags.GetOptionalString("out")
	fileMimeType := c.Flags.GetOptionalString("mime-type")
	attrValues = c.Flags.GetStringSlice("attr", attrValues, cli.FlagsStringSliceOptions{Min: 0})
	tdfType := c.Flags.GetOptionalString("tdf-type")
	kasURLPath := c.Flags.GetOptionalString("kas-url-path")
	wrappingKeyAlgStr := c.Flags.GetOptionalString("wrapping-key-algorithm")
	policyMode := c.Flags.GetOptionalString("policy-mode")
	targetMode := c.Flags.GetOptionalString("target-mode")
	var wrappingKeyAlgorithm ocrypto.KeyType
	switch wrappingKeyAlgStr {
	case string(ocrypto.RSA2048Key):
		wrappingKeyAlgorithm = ocrypto.RSA2048Key
	case string(ocrypto.EC256Key):
		wrappingKeyAlgorithm = ocrypto.EC256Key
	case string(ocrypto.EC384Key):
		wrappingKeyAlgorithm = ocrypto.EC384Key
	case string(ocrypto.EC521Key):
		wrappingKeyAlgorithm = ocrypto.EC521Key
	default:
		wrappingKeyAlgorithm = ocrypto.RSA2048Key
	}

	piped := readPipedStdin()

	inputCount := 0
	if filePath != "" {
		inputCount++
	}
	if len(piped) > 0 {
		inputCount++
	}

	cliExit := func(s string) {
		cli.ExitWithError("Must provide "+s+" of the following to encrypt: [file argument, stdin input]", nil)
	}
	if inputCount == 0 {
		cliExit("ONE")
	} else if inputCount > 1 {
		cliExit("ONLY ONE")
	}

	// prefer filepath argument over stdin input
	bytesSlice := piped
	var err error
	if filePath != "" {
		bytesSlice, err = utils.ReadBytesFromFile(filePath, MaxFileSize)
		if err != nil {
			cli.ExitWithError("Failed to read file:", err)
		}
	}

	// auto-detect mime type if not provided
	if fileMimeType == "" {
		slog.Debug("Detecting mime type of file")
		// get the mime type of the file
		mimetype.SetLimit(Size1MB) // limit to 1MB
		m := mimetype.Detect(bytesSlice)
		// default to application/octet-stream if no mime type is detected
		fileMimeType = m.String()

		if fileMimeType == "application/octet-stream" {
			if fileExt != "" {
				fileMimeType = mimetype.Lookup(fileExt).String()
			}
		}
	}
	slog.Debug("Encrypting file",
		slog.Int("file-len", len(bytesSlice)),
		slog.String("mime-type", fileMimeType),
	)

	// Do the encryption
	encrypted, err := h.EncryptBytes(
		tdfType,
		bytesSlice,
		attrValues,
		fileMimeType,
		kasURLPath,
		c.Flags.GetOptionalBool("ecdsa-binding"),
		assertions,
		wrappingKeyAlgorithm,
		policyMode,
		targetMode,
	)
	if err != nil {
		cli.ExitWithError("Failed to encrypt", err)
	}

	// Find the destination as the output flag filename or stdout
	var dest *os.File
	if out != "" {
		// make sure output ends in .tdf extension
		if !strings.HasSuffix(out, ".tdf") {
			out += ".tdf"
		}
		tdfFile, err := os.Create(out)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to write encrypted file %s", out), err)
		}
		defer tdfFile.Close()
		dest = tdfFile
	} else {
		dest = os.Stdout
	}

	_, e := io.Copy(dest, encrypted)
	if e != nil {
		cli.ExitWithError("Failed to write encrypted data to stdout", e)
	}
}

func init() {
	encrypt := man.Docs.GetCommand("encrypt",
		man.WithRun(encryptRun),
	)
	encrypt.Flags().StringP(
		encrypt.GetDocFlag("out").Name,
		encrypt.GetDocFlag("out").Shorthand,
		encrypt.GetDocFlag("out").Default,
		encrypt.GetDocFlag("out").Description,
	)
	encrypt.Flags().StringSliceVarP(
		&attrValues,
		encrypt.GetDocFlag("attr").Name,
		encrypt.GetDocFlag("attr").Shorthand,
		[]string{},
		encrypt.GetDocFlag("attr").Description,
	)
	encrypt.Flags().StringVarP(
		&assertions,
		encrypt.GetDocFlag("with-assertions").Name,
		encrypt.GetDocFlag("with-assertions").Shorthand,
		"",
		encrypt.GetDocFlag("with-assertions").Description,
	)
	encrypt.Flags().String(
		encrypt.GetDocFlag("mime-type").Name,
		encrypt.GetDocFlag("mime-type").Default,
		encrypt.GetDocFlag("mime-type").Description,
	)
	encrypt.Flags().String(
		encrypt.GetDocFlag("tdf-type").Name,
		encrypt.GetDocFlag("tdf-type").Default,
		encrypt.GetDocFlag("tdf-type").Description,
	)
	encrypt.Flags().StringP(
		encrypt.GetDocFlag("wrapping-key-algorithm").Name,
		encrypt.GetDocFlag("wrapping-key-algorithm").Shorthand,
		encrypt.GetDocFlag("wrapping-key-algorithm").Default,
		encrypt.GetDocFlag("wrapping-key-algorithm").Description,
	)
	encrypt.Flags().Bool(
		encrypt.GetDocFlag("ecdsa-binding").Name,
		false,
		encrypt.GetDocFlag("ecdsa-binding").Description,
	)
	encrypt.Flags().String(
		encrypt.GetDocFlag("kas-url-path").Name,
		encrypt.GetDocFlag("kas-url-path").Default,
		encrypt.GetDocFlag("kas-url-path").Description,
	)
	encrypt.Flags().String(
		encrypt.GetDocFlag("policy-mode").Name,
		encrypt.GetDocFlag("policy-mode").Default,
		encrypt.GetDocFlag("policy-mode").Description,
	)
	encrypt.Flags().String(
		encrypt.GetDocFlag("target-mode").Name,
		encrypt.GetDocFlag("target-mode").Default,
		encrypt.GetDocFlag("target-mode").Description,
	)
	encrypt.GroupID = TDF

	EncryptCmd = &encrypt.Command
}
