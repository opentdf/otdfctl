package cmd

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

const (
	TDF3     = "tdf3"
	NANO     = "nano"
	Size_1MB = 1024 * 1024
)

var attrValues []string
var assertions string

func dev_tdfEncryptCmd(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
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
	if tdfType == "" {
		tdfType = TDF3
	}
	kasURLPath := c.Flags.GetOptionalString("kas-url-path")

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
	if filePath != "" {
		bytesSlice = readBytesFromFile(filePath)
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
	var encrypted *bytes.Buffer
	var err error
	switch tdfType {
	case TDF3:
		encrypted, err = h.EncryptBytes(bytesSlice, attrValues, fileMimeType, kasURLPath, assertions)
	case NANO:
		ecdsaBinding := c.Flags.GetOptionalBool("ecdsa-binding")
		encrypted, err = h.EncryptNanoBytes(bytesSlice, attrValues, kasURLPath, ecdsaBinding)
	default:
		cli.ExitWithError("Failed to encrypt", fmt.Errorf("unrecognized tdf-type: %s", tdfType))
	}
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
	encryptCmd.Flags().StringP(
		encryptCmd.GetDocFlag("tdf-type").Name,
		encryptCmd.GetDocFlag("tdf-type").Shorthand,
		encryptCmd.GetDocFlag("tdf-type").Default,
		encryptCmd.GetDocFlag("tdf-type").Description,
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
	encryptCmd.Command.GroupID = TDF

	RootCmd.AddCommand(&encryptCmd.Command)
}
