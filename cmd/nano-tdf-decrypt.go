package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

func dev_nanoTdfDecryptCmd(cmd *cobra.Command, args []string) {
	h := NewHandler(cmd)
	defer h.Close()

	flagHelper := cli.NewFlagHelper(cmd)
	output := flagHelper.GetOptionalString("out")

	// check for piped input
	piped := readPipedStdin()

	// Prefer file argument over piped input over default filename
	var bytesToDecrypt []byte
	var tdfFile string
	if len(args) > 0 {
		tdfFile = args[0]
		bytesToDecrypt = readBytesFromFile(tdfFile)
	} else if len(piped) > 0 {
		bytesToDecrypt = piped
	} else {
		cli.ExitWithError("Must provide ONE of the following to decrypt: [file argument, stdin input]", errors.New("no input provided"))
	}

	decrypted, err := h.DecryptNanoTDF(bytesToDecrypt)
	if err != nil {
		cli.ExitWithError("Failed to decrypt file", err)
	}

	if output == "" {
		// Print decrypted content to stdout
		fmt.Print(decrypted.String())
		return
	}
	// Here 'output' is the filename given with -o
	f, err := os.Create(output)
	if err != nil {
		cli.ExitWithError("Failed to write decrypted data to file", err)
	}
	defer f.Close()
	_, err = f.Write(decrypted.Bytes())
	if err != nil {
		cli.ExitWithError("Failed to write decrypted data to file", err)
	}
}

func init() {
	decryptCmd := man.Docs.GetCommand("decrypt-nano",
		man.WithRun(dev_nanoTdfDecryptCmd),
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
