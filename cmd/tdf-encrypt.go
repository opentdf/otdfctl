package cmd

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/lib/ocrypto"
	"github.com/spf13/cobra"
)

const (
	TDFTYPE_ZTDF = "ztdf"
	TDF3         = "tdf3"
	NANO         = "nano"
	Size_1MB     = 1024 * 1024
)

var attrValues []string
var assertions string

const INPUT_MAX_FILE_SIZE = int64(10 * 1024 * 1024 * 1024) // 10 GB

func dev_tdfEncryptCmd(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args, cli.WithPrintJSON())
	h := NewHandler(c)
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
		bytesSlice, err = utils.ReadBytesFromFile(filePath, INPUT_MAX_FILE_SIZE)
		if err != nil {
			cli.ExitWithError("Failed to read file:", err)
		}
	}

	// auto-detect mime type if not provided
	if fileMimeType == "" {
		slog.Debug("Detecting mime type of file")
		// get the mime type of the file
		mimetype.SetLimit(Size_1MB) // limit to 1MB
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
	encryptCmd := man.Docs.GetCommand("encrypt",
		man.WithRun(dev_tdfEncryptCmd),
	)
	encryptCmd.Flags().StringP(
		encryptCmd.GetDocFlag("out").Name,
		encryptCmd.GetDocFlag("out").Shorthand,
		encryptCmd.GetDocFlag("out").Default,
		encryptCmd.GetDocFlag("out").Description,
	)
	encryptCmd.Flags().StringSliceVarP(
		&attrValues,
		encryptCmd.GetDocFlag("attr").Name,
		encryptCmd.GetDocFlag("attr").Shorthand,
		[]string{},
		encryptCmd.GetDocFlag("attr").Description,
	)
	encryptCmd.Flags().StringVarP(
		&assertions,
		encryptCmd.GetDocFlag("with-assertions").Name,
		encryptCmd.GetDocFlag("with-assertions").Shorthand,
		"",
		encryptCmd.GetDocFlag("with-assertions").Description,
	)
	encryptCmd.Flags().String(
		encryptCmd.GetDocFlag("mime-type").Name,
		encryptCmd.GetDocFlag("mime-type").Default,
		encryptCmd.GetDocFlag("mime-type").Description,
	)
	encryptCmd.Flags().String(
		encryptCmd.GetDocFlag("tdf-type").Name,
		encryptCmd.GetDocFlag("tdf-type").Default,
		encryptCmd.GetDocFlag("tdf-type").Description,
	)
	encryptCmd.Flags().StringP(
		encryptCmd.GetDocFlag("wrapping-key-algorithm").Name,
		encryptCmd.GetDocFlag("wrapping-key-algorithm").Shorthand,
		encryptCmd.GetDocFlag("wrapping-key-algorithm").Default,
		encryptCmd.GetDocFlag("wrapping-key-algorithm").Description,
	)
	encryptCmd.Flags().Bool(
		encryptCmd.GetDocFlag("ecdsa-binding").Name,
		false,
		encryptCmd.GetDocFlag("ecdsa-binding").Description,
	)
	encryptCmd.Flags().String(
		encryptCmd.GetDocFlag("kas-url-path").Name,
		encryptCmd.GetDocFlag("kas-url-path").Default,
		encryptCmd.GetDocFlag("kas-url-path").Description,
	)
	encryptCmd.Flags().String(
		encryptCmd.GetDocFlag("policy-mode").Name,
		encryptCmd.GetDocFlag("policy-mode").Default,
		encryptCmd.GetDocFlag("policy-mode").Description,
	)
	encryptCmd.Flags().String(
		encryptCmd.GetDocFlag("target-mode").Name,
		encryptCmd.GetDocFlag("target-mode").Default,
		encryptCmd.GetDocFlag("target-mode").Description,
	)
	encryptCmd.GroupID = TDF

	RootCmd.AddCommand(&encryptCmd.Command)
}
