package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"os"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

func dev_tdfDecryptCmd(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	output := c.Flags.GetOptionalString("out")
	tdfType := c.Flags.GetOptionalString("tdf-type")
	if tdfType == "" {
		tdfType = TDF3
	}

	// check for piped input
	piped := readPipedStdin()

	// Prefer file argument over piped input over default filename
	bytesToDecrypt := piped
	var tdfFile string
	if len(args) > 0 {
		tdfFile = args[0]
		bytesToDecrypt = readBytesFromFile(tdfFile)
	}

	if len(bytesToDecrypt) == 0 {
		cli.ExitWithError("Must provide ONE of the following to decrypt: [file argument, stdin input]", errors.New("no input provided"))
	}

	var decrypted *bytes.Buffer
	var err error
	if tdfType == TDF3 {
		decrypted, err = h.DecryptTDF(bytesToDecrypt)
	} else if tdfType == NANO {
		decrypted, err = h.DecryptNanoTDF(bytesToDecrypt)
	} else {
		cli.ExitWithError("Failed to decrypt", fmt.Errorf("unrecognized tdf-type: %s", tdfType))
	}
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
	decryptCmd := man.Docs.GetCommand("decrypt",
		man.WithRun(dev_tdfDecryptCmd),
	)
	decryptCmd.Flags().StringP(
		decryptCmd.GetDocFlag("out").Name,
		decryptCmd.GetDocFlag("out").Shorthand,
		decryptCmd.GetDocFlag("out").Default,
		decryptCmd.GetDocFlag("out").Description,
	)
	decryptCmd.Flags().StringP(
		decryptCmd.GetDocFlag("tdf-type").Name,
		decryptCmd.GetDocFlag("tdf-type").Shorthand,
		decryptCmd.GetDocFlag("tdf-type").Default,
		decryptCmd.GetDocFlag("tdf-type").Description,
	)
	decryptCmd.Command.GroupID = "tdf"

	RootCmd.AddCommand(&decryptCmd.Command)
}
