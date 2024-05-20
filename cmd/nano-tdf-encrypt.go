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

func dev_nanoTdfEncryptCmd(cmd *cobra.Command, args []string) {
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
	encrypted, err := h.EncryptNanoBytes(bytes, values)
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
	encryptNanoCmd := man.Docs.GetCommand("encrypt-nano",
		man.WithRun(dev_nanoTdfEncryptCmd),
	)
	encryptNanoCmd.Flags().StringP(
		encryptNanoCmd.GetDocFlag("out").Name,
		encryptNanoCmd.GetDocFlag("out").Shorthand,
		encryptNanoCmd.GetDocFlag("out").Default,
		encryptNanoCmd.GetDocFlag("out").Description,
	)
	encryptNanoCmd.Flags().StringSliceVarP(
		&attrValues,
		encryptNanoCmd.GetDocFlag("attr").Name,
		encryptNanoCmd.GetDocFlag("attr").Shorthand,
		[]string{},
		encryptNanoCmd.GetDocFlag("attr").Description,
	)
	encryptNanoCmd.Command.GroupID = "tdf"

	RootCmd.AddCommand(&encryptNanoCmd.Command)
}
