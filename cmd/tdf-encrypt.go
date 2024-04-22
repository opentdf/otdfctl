package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

func dev_tdfEncryptCmd(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	var filePath string
	if len(args) > 0 {
		filePath = args[0]
	}
	out := flagHelper.GetOptionalString("out")
	values := flagHelper.GetStringSlice("attr", attrValues, cli.FlagHelperStringSliceOptions{Min: 0})

	// Read bytes from stdin without blocking by checking size first
	var piped []byte
	in := os.Stdin
	stdin, err := in.Stat()
	if err != nil {
		cli.ExitWithError("Failed to read from stdin", err)
	}
	size := stdin.Size()
	if size > 0 {
		piped, err = io.ReadAll(os.Stdin)
		if err != nil {
			cli.ExitWithError("Failed to read from stdin", err)
		}
	}

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

	var bytes []byte
	if filePath != "" {
		fileToEncrypt, err := os.Open(filePath)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to open file at path: %s", filePath), err)
		}
		defer fileToEncrypt.Close()

		bytes, err = io.ReadAll(fileToEncrypt)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to read bytes from file at path: %s", filePath), err)
		}
		// default <filename.extension>.tdf as output
		if out == "" {
			out = filePath
		}
	} else {
		bytes = piped
	}
	tdfFile, err := h.EncryptBytes(bytes, values, out)
	if err != nil {
		cli.ExitWithError("Failed to encrypt", err)
	}
	defer tdfFile.Close()

	f, err := os.Open(tdfFile.Name())
	if err != nil {
		cli.ExitWithError(fmt.Sprintf("Failed to write encrypted file %s to stdout", tdfFile.Name()), err)
	}

	_, e := io.Copy(os.Stdout, f)
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
