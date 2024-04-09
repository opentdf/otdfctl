package handlers

import (
	"os"
	"strings"

	"github.com/opentdf/platform/sdk"
)

func (h Handler) EncryptText(text string, values []string) (*sdk.TDFObject, error) {
	tdfFile, err := os.Create("sensitive.txt.tdf")
	if err != nil {
		return nil, err
	}
	defer tdfFile.Close()

	// TODO: validate values are FQNs and add to TDF

	// Encrypt the text
	return h.sdk.CreateTDF(tdfFile, strings.NewReader(text),
		// sdk.WithDataAttributes(values...),
		sdk.WithKasInformation(sdk.KASInfo{
			URL:       h.platformEndpoint,
			PublicKey: "",
		},
		),
	)
}

func (h Handler) EncryptFile(filePath string, values []string) (*sdk.TDFObject, error) {
	return nil, nil
}

func (h Handler) DecryptTDF(filePath string) (*sdk.Reader, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Decrypt the TDF
	return h.sdk.LoadTDF(file)
}
