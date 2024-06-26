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
	TDF3 = "tdf3"
	NANO = "nano"
)

func dev_tdfEncryptCmd(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	var filePath string
	var fileExt string
	if len(args) > 0 {
		filePath = args[0]
		fileExt = strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))
	}

	out := flagHelper.GetOptionalString("out")
	fileMimeType := flagHelper.GetOptionalString("mime-type")
	values := flagHelper.GetStringSlice("attr", attrValues, cli.FlagHelperStringSliceOptions{Min: 0})
	tdfType := flagHelper.GetOptionalString("tdf-type")
	if tdfType == "" {
		tdfType = TDF3
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
	if filePath != "" {
		bytesSlice = readBytesFromFile(filePath)
	}

	// auto-detect mime type if not provided
	if fileMimeType == "" {
		slog.Debug("Detecting mime type of file")
		// get the mime type of the file
		mimetype.SetLimit(1024 * 1024) // limit to 1MB
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
	if tdfType == TDF3 {
		encrypted, err = h.EncryptBytes(bytesSlice, values, fileMimeType)
	} else if tdfType == NANO {
		encrypted, err = h.EncryptNanoBytes(bytesSlice, values)
	} else {
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
	encryptCmd.Command.GroupID = "tdf"

	RootCmd.AddCommand(&encryptCmd.Command)
}
