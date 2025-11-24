package tdf

import (
	"errors"
	"fmt"
	"os"

	"github.com/opentdf/otdfctl/cmd/common"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/lib/ocrypto"
	"github.com/spf13/cobra"
)

var assertionVerification string
var kasAllowList []string

func decryptRun(cmd *cobra.Command, args []string) {
	c := cli.New(cmd, args, cli.WithPrintJSON())
	h := common.NewHandler(c)
	defer h.Close()

	output := c.Flags.GetOptionalString("out")
	disableAssertionVerification := c.Flags.GetOptionalBool("no-verify-assertions")
	sessionKeyAlgStr := c.Flags.GetOptionalString("session-key-algorithm")
	var sessionKeyAlgorithm ocrypto.KeyType
	switch sessionKeyAlgStr {
	case string(ocrypto.RSA2048Key):
		sessionKeyAlgorithm = ocrypto.RSA2048Key
	case string(ocrypto.EC256Key):
		sessionKeyAlgorithm = ocrypto.EC256Key
	case string(ocrypto.EC384Key):
		sessionKeyAlgorithm = ocrypto.EC384Key
	case string(ocrypto.EC521Key):
		sessionKeyAlgorithm = ocrypto.EC521Key
	default:
		sessionKeyAlgorithm = ocrypto.RSA2048Key
	}

	// check for piped input
	piped := readPipedStdin()

	// Prefer file argument over piped input over default filename
	bytesToDecrypt := piped
	var tdfFile string
	var err error
	if len(args) > 0 {
		tdfFile = args[0]
		bytesToDecrypt, err = utils.ReadBytesFromFile(tdfFile, MaxFileSize)
		if err != nil {
			cli.ExitWithError("Failed to read file:", err)
		}
	}

	if len(bytesToDecrypt) == 0 {
		cli.ExitWithError("Must provide ONE of the following to decrypt: [file argument, stdin input]", errors.New("no input provided"))
	}

	ignoreAllowlist := len(kasAllowList) == 1 && kasAllowList[0] == "*"

	decrypted, err := h.DecryptBytes(
		c.Context(),
		bytesToDecrypt,
		assertionVerification,
		disableAssertionVerification,
		sessionKeyAlgorithm,
		kasAllowList,
		ignoreAllowlist,
		nil,
	)
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
	decrypt := man.Docs.GetCommand("decrypt",
		man.WithRun(decryptRun),
	)
	decrypt.Flags().StringP(
		decrypt.GetDocFlag("out").Name,
		decrypt.GetDocFlag("out").Shorthand,
		decrypt.GetDocFlag("out").Default,
		decrypt.GetDocFlag("out").Description,
	)
	// deprecated flag
	decrypt.Flags().StringP(
		decrypt.GetDocFlag("tdf-type").Name,
		decrypt.GetDocFlag("tdf-type").Shorthand,
		decrypt.GetDocFlag("tdf-type").Default,
		decrypt.GetDocFlag("tdf-type").Description,
	)
	decrypt.Flags().StringVarP(
		&assertionVerification,
		decrypt.GetDocFlag("with-assertion-verification-keys").Name,
		decrypt.GetDocFlag("with-assertion-verification-keys").Shorthand,
		"",
		decrypt.GetDocFlag("with-assertion-verification-keys").Description,
	)
	decrypt.Flags().String(
		decrypt.GetDocFlag("session-key-algorithm").Name,
		decrypt.GetDocFlag("session-key-algorithm").Default,
		decrypt.GetDocFlag("session-key-algorithm").Description,
	)
	decrypt.Flags().Bool(
		decrypt.GetDocFlag("no-verify-assertions").Name,
		decrypt.GetDocFlag("no-verify-assertions").DefaultAsBool(),
		decrypt.GetDocFlag("no-verify-assertions").Description,
	)
	decrypt.Flags().StringSliceVarP(
		&kasAllowList,
		decrypt.GetDocFlag("kas-allowlist").Name,
		decrypt.GetDocFlag("kas-allowlist").Shorthand,
		nil,
		decrypt.GetDocFlag("kas-allowlist").Description,
	)

	decrypt.GroupID = TDF

	DecryptCmd = &decrypt.Command
}
