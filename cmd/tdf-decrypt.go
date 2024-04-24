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
	output := flagHelper.GetOptionalString("out")

	// check for a TDF flag argument
	var tdfFile string
	if len(args) > 0 {
		tdfFile = args[0]
	}
	// default to sensitive.txt.tdf if no file is provided
	if tdfFile == "" {
		tdfFile = "sensitive.txt.tdf"
	}

	decrypted, err := h.DecryptTDF(tdfFile)
	if err != nil {
		cli.ExitWithError("Failed to decrypt file", err)
	}

	if output == "stdout" {
		// Print decrypted content to stdout
		fmt.Print(decrypted.String())
		return
	} else if output != "" {
		// Here 'output' is the filename given with -o
		// Write decrypted string to file with stripped .tdf extension
		fmt.Println("Decrypting " + tdfFile)
		f, err := os.Create(output)
		if err != nil {
			cli.ExitWithError("Failed to write decrypted data to file", err)
		}
		defer f.Close()
		_, err = f.Write(decrypted.Bytes())
		if err != nil {
			cli.ExitWithError("Failed to write decrypted data to file", err)
		}
		return
	} else {
		// Here the 'output' filename is based on the input filename
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
}

func init() {
	decryptCmd := man.Docs.GetCommand("decrypt",
		man.WithRun(dev_tdfDecryptCmd),
	)
	decryptCmd.Flags().StringP(
		decryptCmd.GetDocFlag("out").Name,
		decryptCmd.GetDocFlag("out").Shorthand,
		decryptCmd.GetDocFlag("out").Default,
		decryptCmd.GetDocFlag("out").Description,
	)
	decryptCmd.Command.GroupID = "tdf"

	RootCmd.AddCommand(&decryptCmd.Command)
}
