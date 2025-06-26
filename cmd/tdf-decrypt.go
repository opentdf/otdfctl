package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"

	decryptgenerated "github.com/opentdf/otdfctl/cmd/generated/decrypt"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/lib/ocrypto"
	"github.com/spf13/cobra"
)

var TDF = "tdf"

var assertionVerification string
var kasAllowList []string

const TDF_MAX_FILE_SIZE = int64(10 * 1024 * 1024 * 1024) // 10 GB

// handleTdfDecrypt implements the business logic for the decrypt command
func handleTdfDecrypt(cmd *cobra.Command, req *decryptgenerated.DecryptRequest) error {
	// Handle file arguments - decrypt can take a file argument or read from stdin
	args := []string{}
	if cmd.Context() != nil {
		if ctxArgs := cmd.Context().Value("args"); ctxArgs != nil {
			args = ctxArgs.([]string)
		}
	}
	
	c := cli.New(cmd, args, cli.WithPrintJson())
	h := NewHandler(c)
	defer h.Close()

	// Extract flags from the generated request structure
	output := req.Flags.Out
	disableAssertionVerification := req.Flags.NoVerifyAssertions != "" // Convert string to bool
	sessionKeyAlgStr := req.Flags.SessionKeyAlgorithm
	withAssertionVerificationKeys := req.Flags.WithAssertionVerificationKeys
	kasAllowlistStr := req.Flags.KasAllowlist
	
	// Set global variables for compatibility with existing logic
	assertionVerification = withAssertionVerificationKeys
	if kasAllowlistStr != "" {
		kasAllowList = []string{kasAllowlistStr} // Convert single string to slice
	}
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
		bytesToDecrypt, err = utils.ReadBytesFromFile(tdfFile, TDF_MAX_FILE_SIZE)
		if err != nil {
			cli.ExitWithError("Failed to read file:", err)
		}
	}

	if len(bytesToDecrypt) == 0 {
		cli.ExitWithError("Must provide ONE of the following to decrypt: [file argument, stdin input]", errors.New("no input provided"))
	}

	ignoreAllowlist := len(kasAllowList) == 1 && kasAllowList[0] == "*"

	decrypted, err := h.DecryptBytes(
		bytesToDecrypt,
		assertionVerification,
		disableAssertionVerification,
		sessionKeyAlgorithm,
		kasAllowList,
		ignoreAllowlist,
	)
	if err != nil {
		cli.ExitWithError("Failed to decrypt file", err)
	}

	if output == "" {
		//nolint:forbidigo // printing decrypted content to stdout
		fmt.Print(decrypted.String())
		return nil
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
	
	return nil
}

func init() {
	// Create command using generated constructor with handler function
	decryptCmd := decryptgenerated.NewDecryptCommand(handleTdfDecrypt)
	decryptCmd.GroupID = TDF

	// Override the RunE to capture args properly for file handling
	originalRunE := decryptCmd.RunE
	decryptCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// Store args in context for the handler to access
		ctx := context.WithValue(cmd.Context(), "args", args)
		cmd.SetContext(ctx)
		return originalRunE(cmd, args)
	}

	// Add to root command
	RootCmd.AddCommand(decryptCmd)
}
