package handlers

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/opentdf/platform/sdk"
)

func (h Handler) EncryptBytes(b []byte, values []string, out string) (*sdk.TDFObject, error) {
	if out == "" {
		out = "sensitive.txt"
	}
	tdfFile, err := os.Create(fmt.Sprintf("%s.tdf", out))
	if err != nil {
		return nil, err
	}
	defer tdfFile.Close()

	// TODO: validate values are FQNs and add to TDF

	return h.sdk.CreateTDF(tdfFile, bytes.NewReader(b),
		// sdk.WithDataAttributes(values...),
		sdk.WithKasInformation(sdk.KASInfo{
			URL:       fmt.Sprintf("http://%s", h.platformEndpoint),
			PublicKey: "",
		},
		),
	)
}

func (h Handler) DecryptTDF(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decrypt the TDF
	tdfreader, err := h.sdk.LoadTDF(file)
	if err != nil {
		return "", err
	}

	buf := new(strings.Builder)
	_, err = io.Copy(buf, tdfreader)
	if err != nil && err != io.EOF {
		return "", err
	}
	return buf.String(), nil
}
