package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/sdk"
	"github.com/spf13/cobra"
)

var dev_tdfCmd *cobra.Command

func dev_tdfEncryptCmd(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	filePath := flagHelper.GetOptionalString("file")
	text := flagHelper.GetOptionalString("text")
	values := flagHelper.GetStringSlice("attr-value", attrValues, cli.FlagHelperStringSliceOptions{Min: 0})

	if filePath == "" && text == "" {
		cli.ExitWithError("Must provide a file or text to encrypt", nil)
	} else if filePath != "" && text != "" {
		cli.ExitWithError("Cannot provide both a file and text to encrypt", nil)
	}

	var (
		tdfFile *sdk.TDFObject
		err     error
	)
	if filePath != "" {
		cli.ExitWithError("File encryption not yet implemented", nil)
		// tdfFile, err = h.EncryptFile(filePath, values)
		// if err != nil {
		// 	cli.ExitWithError("Failed to encrypt file", err)
		// }
	} else {
		tdfFile, err = h.EncryptText(text, values)
		if err != nil {
			cli.ExitWithError("Failed to encrypt text", err)
		}
	}

	cli.SuccessMessage(fmt.Sprintf("Successfully encrypted data.\nTDF manifest: %+v", tdfFile.Manifest()))
}

func dev_tdfDecryptCmd(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	tdfFile := flagHelper.GetOptionalString("tdf")
	output := flagHelper.GetOptionalString("output")

	if tdfFile == "" {
		cli.ExitWithError("Must provide a TDF file to decrypt", nil)
	}

	decrypted, err := h.DecryptTDF(tdfFile)
	if err != nil {
		cli.ExitWithError("Failed to decrypt file", err)
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, decrypted)
	if err != nil && err != io.EOF {
		cli.ExitWithError("Failed to read decrypted data", err)
	}

	cli.SuccessMessage("Successfully decrypted TDF file: " + tdfFile)
	if output == "file" {
		// Write decrypted string to file with stripped .tdf extension
		f, err := os.Create(strings.Replace(tdfFile, ".tdf", "", 1))
		if err != nil {
			cli.ExitWithError("Failed to write decrypted data to file", err)
		}
		defer f.Close()
		_, err = f.WriteString(buf.String())
		if err != nil {
			cli.ExitWithError("Failed to write decrypted data to file", err)
		}
	} else {
		// Print decrypted string
		fmt.Println(buf.String())
	}
}

func init() {
	encryptCmd := man.Docs.GetCommand("dev/tdf/encrypt",
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
	encryptCmd.Flags().StringSliceVarP(
		&attrValues,
		encryptCmd.GetDocFlag("attr-value").Name,
		encryptCmd.GetDocFlag("attr-value").Shorthand,
		[]string{},
		encryptCmd.GetDocFlag("attr-value").Description,
	)
	// TODO: should we have auth values pulled in from a config like the platform? See config: https://github.com/opentdf/platform/blob/main/opentdf-example.yaml#L16
	// NOTE: starting with unauthenticated, insecure platform implementation

	decryptCmd := man.Docs.GetCommand("dev/tdf/decrypt",
		man.WithRun(dev_tdfDecryptCmd),
	)
	decryptCmd.Flags().StringP(
		decryptCmd.GetDocFlag("tdf").Name,
		decryptCmd.GetDocFlag("tdf").Shorthand,
		decryptCmd.GetDocFlag("tdf").Default,
		decryptCmd.GetDocFlag("tdf").Description,
	)
	decryptCmd.Flags().StringP(
		decryptCmd.GetDocFlag("output").Name,
		decryptCmd.GetDocFlag("output").Shorthand,
		decryptCmd.GetDocFlag("output").Default,
		decryptCmd.GetDocFlag("output").Description,
	)

	doc := man.Docs.GetCommand("dev/tdf",
		man.WithSubcommands(encryptCmd, decryptCmd),
	)
	dev_tdfCmd = &doc.Command
	devCmd.AddCommand(dev_tdfCmd)
}
