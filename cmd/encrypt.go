package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var dev_tdfCmd *cobra.Command

func dev_tdfEncryptCmd(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	filePath := flagHelper.GetOptionalString("file")
	text := flagHelper.GetOptionalString("text")
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

	if filePath == "" && text == "" && len(piped) == 0 {
		cli.ExitWithError("Must provide ONE of the following to encrypt: [text, file, stdin input]", nil)
	} else if filePath != "" && text != "" && len(piped) > 0 {
		cli.ExitWithError("Must provide ONLY ONE of the following to encrypt: [text, file, stdin input]", nil)
	}

	var bytes []byte
	if filePath != "" {
		fileToEncrypt, err := os.Open(filePath)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to open file at path: %s", filePath), err)
		}
		defer fileToEncrypt.Close()

		bytes, err = ioutil.ReadAll(fileToEncrypt)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to read bytes from file at path: %s", filePath), err)
		}
		// default <filename.extension>.tdf as output
		if out == "" {
			out = filePath
		}
	} else if text != "" {
		bytes = []byte(text)
	} else {
		bytes = piped
	}
	tdfFile, err := h.EncryptBytes(bytes, values, out)
	if err != nil {
		cli.ExitWithError("Failed to encrypt text", err)
	}

	cli.SuccessMessage(fmt.Sprintf("Successfully encrypted data.\nTDF manifest: %+v", tdfFile.Manifest()))
}

func init() {
	encryptCmd := man.Docs.GetCommand("encrypt",
		man.WithRun(dev_tdfEncryptCmd),
	)
	encryptCmd.Flags().StringP(
		encryptCmd.GetDocFlag("file").Name,
		encryptCmd.GetDocFlag("file").Shorthand,
		encryptCmd.GetDocFlag("file").Default,
		encryptCmd.GetDocFlag("file").Description,
	)
	encryptCmd.Flags().StringP(
		encryptCmd.GetDocFlag("text").Name,
		encryptCmd.GetDocFlag("text").Shorthand,
		encryptCmd.GetDocFlag("text").Default,
		encryptCmd.GetDocFlag("text").Description,
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

	rootCmd.AddCommand(&encryptCmd.Command)
}
