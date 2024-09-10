package handlers

import (
	"bytes"
	"io"

	"github.com/opentdf/platform/sdk"
)

func (h Handler) EncryptNanoBytes(b []byte, values []string, kasUrlPath string, ecdsaBinding bool) (*bytes.Buffer, error) {
	var encrypted []byte
	enc := bytes.NewBuffer(encrypted)

	options := []sdk.NanoTDFOption{
		sdk.WithKasURL(h.platformEndpoint + kasUrlPath),
		sdk.WithNanoDataAttributes(values),
	}
	if ecdsaBinding {
		options = append(options, sdk.WithECDSAPolicyBinding())
	}
	// TODO: validate values are FQNs or return an error [https://github.com/opentdf/platform/issues/515]
	_, err := h.sdk.CreateNanoTDF(enc, bytes.NewReader(b), options...)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (h Handler) DecryptNanoTDF(toDecrypt []byte) (*bytes.Buffer, error) {
	outBuf := bytes.Buffer{}
	_, err := h.sdk.ReadNanoTDF(io.Writer(&outBuf), bytes.NewReader(toDecrypt))
	if err != nil {
		return nil, err
	}
	return &outBuf, nil
}
