package handlers

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/opentdf/platform/sdk"
)

func (h Handler) EncryptBytes(b []byte, values []string, out string) (*os.File, error) {
	if out == "" {
		out = "sensitive.txt"
	}

	tdfFile, err := os.Create(fmt.Sprintf("%s.tdf", out))
	if err != nil {
		return nil, err
	}

	// TODO: validate values are FQNs or return an error [https://github.com/opentdf/platform/issues/515]
	_, err = h.sdk.CreateTDF(tdfFile, bytes.NewReader(b),
		sdk.WithDataAttributes(values...),
		sdk.WithKasInformation(sdk.KASInfo{
			URL:       fmt.Sprintf("http://%s", h.platformEndpoint),
			PublicKey: "",
		},
		),
	)
	if err != nil {
		return nil, err
	}
	return tdfFile, nil
}

func (h Handler) DecryptTDF(toDecrypt []byte) (*bytes.Buffer, error) {
	tdfreader, err := h.sdk.LoadTDF(bytes.NewReader(toDecrypt))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, tdfreader)
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buf, nil
}
