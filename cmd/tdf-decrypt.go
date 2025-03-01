package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/lib/ocrypto"
	"github.com/spf13/cobra"
)

const (
	RSA2048     = "rsa:2048"
	ECSECP256R1 = "ec:secp256r1"
	ECSECP384R1 = "ec:secp384r1"
	ECSECP521R1 = "ec:secp521r1"
)

var TDF = "tdf"

var assertionVerification string

const TDF_MAX_FILE_SIZE = int64(10 * 1024 * 1024 * 1024) // 10 GB

func dev_tdfDecryptCmd(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args)
	h := NewHandler(c)
	defer h.Close()

	output := c.Flags.GetOptionalString("out")
	disableAssertionVerification := c.Flags.GetOptionalBool("no-verify-assertions")
	sessionKeyAlgStr := c.Flags.GetOptionalString("session-key-algorithm")
	var sessionKeyAlgorithm ocrypto.KeyType
	switch sessionKeyAlgStr {
	case RSA2048:
		sessionKeyAlgorithm = RSA2048
	case ECSECP256R1:
		sessionKeyAlgorithm = ECSECP256R1
	case ECSECP384R1:
		sessionKeyAlgorithm = ECSECP384R1
	case ECSECP521R1:
		sessionKeyAlgorithm = ECSECP521R1
	default:
		sessionKeyAlgorithm = RSA2048
	}

	// check for piped input
	piped := readPipedStdin()

	// Prefer file argument over piped input over default filename
	bytesToDecrypt := piped
	var tdfFile string
	var err error
	if len(args) > 0 {
		tdfFile = args[0]
		bytesToDecrypt, err = utils.ReadBytesFromFile(tdfFile, TDF_MAX_FILE_SIZE)
		if err != nil {
			cli.ExitWithError("Failed to read file:", err)
		}
	}

	if len(bytesToDecrypt) == 0 {
		cli.ExitWithError("Must provide ONE of the following to decrypt: [file argument, stdin input]", errors.New("no input provided"))
	}

	decrypted, err := h.DecryptBytes(bytesToDecrypt, assertionVerification, disableAssertionVerification, sessionKeyAlgorithm)
	if err != nil {
		cli.ExitWithError("Failed to decrypt file", err)
	}

	if output == "" {
		//nolint:forbidigo // printing decrypted content to stdout
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
	// deprecated flag
	decryptCmd.Flags().StringP(
		decryptCmd.GetDocFlag("tdf-type").Name,
		decryptCmd.GetDocFlag("tdf-type").Shorthand,
		decryptCmd.GetDocFlag("tdf-type").Default,
		decryptCmd.GetDocFlag("tdf-type").Description,
	)
	decryptCmd.Flags().StringVarP(
		&assertionVerification,
		decryptCmd.GetDocFlag("with-assertion-verification-keys").Name,
		decryptCmd.GetDocFlag("with-assertion-verification-keys").Shorthand,
		"",
		decryptCmd.GetDocFlag("with-assertion-verification-keys").Description,
	)
	decryptCmd.Flags().StringP(
		decryptCmd.GetDocFlag("session-key-algorithm").Name,
		decryptCmd.GetDocFlag("session-key-algorithm").Shorthand,
		decryptCmd.GetDocFlag("session-key-algorithm").Default,
		decryptCmd.GetDocFlag("session-key-algorithm").Description,
	)
	decryptCmd.Flags().Bool(
		decryptCmd.GetDocFlag("no-verify-assertions").Name,
		decryptCmd.GetDocFlag("no-verify-assertions").DefaultAsBool(),
		decryptCmd.GetDocFlag("no-verify-assertions").Description,
	)

	decryptCmd.Command.GroupID = TDF

	RootCmd.AddCommand(&decryptCmd.Command)
}
