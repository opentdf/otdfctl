package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

func dev_tdfEncryptCmd(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	var filePath string
	if len(args) > 0 {
		filePath = args[0]
	}
	out := flagHelper.GetOptionalString("out")
	values := flagHelper.GetStringSlice("attr", attrValues, cli.FlagHelperStringSliceOptions{Min: 0})

	piped := readPipedStdin()

	inputCount := 0
	if filePath != "" {
		inputCount++
	}
	if len(piped) > 0 {
		inputCount++
	}

	if inputCount == 0 {
		cli.ExitWithError("Must provide ONE of the following to encrypt: [file argument, stdin input]", nil)
	} else if inputCount > 1 {
		cli.ExitWithError("Must provide ONLY ONE of the following to encrypt: [file argument, stdin input]", nil)
	}

	// prefer filepath argument over stdin input
	var bytes []byte
	if filePath != "" {
		bytes = readBytesFromFile(filePath)
	} else {
		bytes = piped
	}

	// Do the encryption
	encrypted, err := h.EncryptBytes(bytes, values)
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
	encryptCmd.Command.GroupID = "tdf"

	RootCmd.AddCommand(&encryptCmd.Command)
}
