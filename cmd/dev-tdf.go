package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/platform/sdk"
	"github.com/spf13/cobra"
)

var (
	dev_tdfCmd = &cobra.Command{
		Use:   "tdf",
		Short: "Demo of TDF",
		Long:  `Demonstration of Trusted Data Format (TDF) functionalities: [encrypt, decrypt] enabled by the OpenTDF platform.`,
	}

	dev_tdfEncryptCmd = &cobra.Command{
		Use:   "encrypt",
		Short: "Encrypt a file or string of text utilizing the TDF and OpenTDF platform.",
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}

	dev_tdfDecryptCmd = &cobra.Command{
		Use:   "decrypt",
		Short: "Decrypt a TDF file empowered by the OpenTDF platform.",
		Run: func(cmd *cobra.Command, args []string) {
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
		},
	}
)

func init() {
	devCmd.AddCommand(dev_tdfCmd)

	dev_tdfCmd.AddCommand(dev_tdfEncryptCmd)
	dev_tdfEncryptCmd.Flags().StringP("file", "f", "", "A file to encrypt that will be encrypted and saved as '<filename>.<extension>.tdf' in the same directory.")
	dev_tdfEncryptCmd.Flags().StringP("text", "t", "", "A string of text to encrypt that will be saved as 'sensitive.txt.tdf' in the $HOME directory.")
	dev_tdfEncryptCmd.Flags().StringSliceVarP(&attrValues, "attr-value", "v", []string{}, "Attribute value Fully Qualified Names (FQNs, i.e. 'https://example.com/attr/attr1/value/value1') to apply to the encrypted data.")
	// TODO: should we have auth values pulled in from a config like the platform? See config: https://github.com/opentdf/platform/blob/main/opentdf-example.yaml#L16
	// NOTE: starting with unauthenticated, insecure platform implementation

	dev_tdfCmd.AddCommand(dev_tdfDecryptCmd)
	dev_tdfDecryptCmd.Flags().StringP("tdf", "t", "sensitive.txt.tdf", "The TDF file with path from $HOME being decrypted (default 'sensitive.txt.tdf')")
	dev_tdfDecryptCmd.Flags().StringP("output", "o", "file", "The decrypted output destination (default 'file', options: 'file', 'stdout')")
}
