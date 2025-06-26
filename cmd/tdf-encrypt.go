package cmd

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	encryptgenerated "github.com/opentdf/otdfctl/cmd/generated/encrypt"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/lib/ocrypto"
	"github.com/spf13/cobra"
)

const (
	TDFTYPE_ZTDF = "ztdf"
	TDF3         = "tdf3"
	NANO         = "nano"
	Size_1MB     = 1024 * 1024
)

var attrValues []string
var assertions string

const INPUT_MAX_FILE_SIZE = int64(10 * 1024 * 1024 * 1024) // 10 GB

// handleTdfEncrypt implements the business logic for the encrypt command
func handleTdfEncrypt(cmd *cobra.Command, req *encryptgenerated.EncryptRequest) error {
	// Handle file arguments - encrypt can take a file argument or read from stdin
	args := []string{}
	if cmd.Context() != nil {
		if ctxArgs := cmd.Context().Value("args"); ctxArgs != nil {
			args = ctxArgs.([]string)
		}
	}
	
	c := cli.New(cmd, args, cli.WithPrintJson())
	h := NewHandler(c)
	defer h.Close()

	var filePath string
	var fileExt string
	if len(args) > 0 {
		filePath = args[0]
		fileExt = strings.ToLower(strings.TrimPrefix(filepath.Ext(filePath), "."))
	}

	// Extract flags from the generated request structure
	out := req.Flags.Out
	fileMimeType := req.Flags.MimeType
	attrValue := req.Flags.Attr
	tdfType := req.Flags.TdfType
	kasURLPath := req.Flags.KasUrlPath
	wrappingKeyAlgStr := req.Flags.WrappingKeyAlgorithm
	targetMode := req.Flags.TargetMode
	ecdsaBinding := req.Flags.EcdsaBinding
	withAssertions := req.Flags.WithAssertions
	
	// Convert single attr string to slice for compatibility with existing logic
	if attrValue != "" {
		attrValues = []string{attrValue}
	} else {
		attrValues = []string{}
	}
	
	// Set assertions for compatibility
	assertions = withAssertions
	var wrappingKeyAlgorithm ocrypto.KeyType
	switch wrappingKeyAlgStr {
	case string(ocrypto.RSA2048Key):
		wrappingKeyAlgorithm = ocrypto.RSA2048Key
	case string(ocrypto.EC256Key):
		wrappingKeyAlgorithm = ocrypto.EC256Key
	case string(ocrypto.EC384Key):
		wrappingKeyAlgorithm = ocrypto.EC384Key
	case string(ocrypto.EC521Key):
		wrappingKeyAlgorithm = ocrypto.EC521Key
	default:
		wrappingKeyAlgorithm = ocrypto.RSA2048Key
	}

	piped := readPipedStdin()

	inputCount := 0
	if filePath != "" {
		inputCount++
	}
	if len(piped) > 0 {
		inputCount++
	}

	cliExit := func(s string) {
		cli.ExitWithError("Must provide "+s+" of the following to encrypt: [file argument, stdin input]", nil)
	}
	if inputCount == 0 {
		cliExit("ONE")
	} else if inputCount > 1 {
		cliExit("ONLY ONE")
	}

	// prefer filepath argument over stdin input
	bytesSlice := piped
	var err error
	if filePath != "" {
		bytesSlice, err = utils.ReadBytesFromFile(filePath, INPUT_MAX_FILE_SIZE)
		if err != nil {
			cli.ExitWithError("Failed to read file:", err)
		}
	}

	// auto-detect mime type if not provided
	if fileMimeType == "" {
		slog.Debug("Detecting mime type of file")
		// get the mime type of the file
		mimetype.SetLimit(Size_1MB) // limit to 1MB
		m := mimetype.Detect(bytesSlice)
		// default to application/octet-stream if no mime type is detected
		fileMimeType = m.String()

		if fileMimeType == "application/octet-stream" {
			if fileExt != "" {
				fileMimeType = mimetype.Lookup(fileExt).String()
			}
		}
	}
	slog.Debug("Encrypting file",
		slog.Int("file-len", len(bytesSlice)),
		slog.String("mime-type", fileMimeType),
	)

	// Do the encryption  
	encrypted, err := h.EncryptBytes(
		tdfType,
		bytesSlice,
		attrValues,
		fileMimeType,
		kasURLPath,
		ecdsaBinding != "", // Convert string to bool (non-empty string means enabled)
		assertions,
		wrappingKeyAlgorithm,
		targetMode,
	)
	if err != nil {
		cli.ExitWithError("Failed to encrypt", err)
	}

	// Find the destination as the output flag filename or stdout
	var dest *os.File
	if out != "" {
		// make sure output ends in .tdf extension
		if !strings.HasSuffix(out, ".tdf") {
			out += ".tdf"
		}
		tdfFile, err := os.Create(out)
		if err != nil {
			cli.ExitWithError(fmt.Sprintf("Failed to write encrypted file %s", out), err)
		}
		defer tdfFile.Close()
		dest = tdfFile
	} else {
		dest = os.Stdout
	}

	_, e := io.Copy(dest, encrypted)
	if e != nil {
		cli.ExitWithError("Failed to write encrypted data to stdout", e)
	}
	
	return nil
}

func init() {
	// Create command using generated constructor with handler function
	encryptCmd := encryptgenerated.NewEncryptCommand(handleTdfEncrypt)
	encryptCmd.GroupID = TDF

	// Override the RunE to capture args properly for file handling
	originalRunE := encryptCmd.RunE
	encryptCmd.RunE = func(cmd *cobra.Command, args []string) error {
		// Store args in context for the handler to access
		ctx := context.WithValue(cmd.Context(), "args", args)
		cmd.SetContext(ctx)
		return originalRunE(cmd, args)
	}

	// Add to root command
	RootCmd.AddCommand(encryptCmd)
}
