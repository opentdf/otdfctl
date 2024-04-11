package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

func dev_tdfDecryptCmd(cmd *cobra.Command, args []string) {
	h := cli.NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	tdfFile := flagHelper.GetOptionalString("file")
	output := flagHelper.GetOptionalString("output")

	if tdfFile == "" {
		cli.ExitWithError("Must provide a TDF file to decrypt", nil)
	}

	decrypted, err := h.DecryptTDF(tdfFile)
	if err != nil {
		cli.ExitWithError("Failed to decrypt file", err)
	}

	if output == "file" {
		// Write decrypted string to file with stripped .tdf extension
		f, err := os.Create(strings.Replace(tdfFile, ".tdf", "", 1))
		if err != nil {
			cli.ExitWithError("Failed to write decrypted data to file", err)
		}
		defer f.Close()
		_, err = f.Write(decrypted.Bytes())
		if err != nil {
			cli.ExitWithError("Failed to write decrypted data to file", err)
		}
		return
	}
	// Print decrypted content to stdout
	fmt.Print(decrypted.String())
}

func init() {
	decryptCmd := man.Docs.GetCommand("decrypt",
		man.WithRun(dev_tdfDecryptCmd),
	)
	decryptCmd.Flags().StringP(
		decryptCmd.GetDocFlag("file").Name,
		decryptCmd.GetDocFlag("file").Shorthand,
		decryptCmd.GetDocFlag("file").Default,
		decryptCmd.GetDocFlag("file").Description,
	)
	decryptCmd.Flags().StringP(
		decryptCmd.GetDocFlag("output").Name,
		decryptCmd.GetDocFlag("output").Shorthand,
		decryptCmd.GetDocFlag("output").Default,
		decryptCmd.GetDocFlag("output").Description,
	)
	decryptCmd.Command.GroupID = "tdf"

	rootCmd.AddCommand(&decryptCmd.Command)
}
