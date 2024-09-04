package handlers

import (
	"bytes"
	"io"
)

func (h Handler) EncryptNanoBytes(b []byte, values []string, kasUrlPath string) (*bytes.Buffer, error) {
	var encrypted []byte
	enc := bytes.NewBuffer(encrypted)

	nanoTDFConfig, err := h.sdk.NewNanoTDFConfig()
	if err != nil {
		return nil, err
	}

	err = nanoTDFConfig.SetKasURL(h.platformEndpoint + kasUrlPath)
	if err != nil {
		return nil, err
	}
	err = nanoTDFConfig.SetAttributes(values)
	if err != nil {
		return nil, err
	}

	// TODO: validate values are FQNs or return an error [https://github.com/opentdf/platform/issues/515]
	_, err = h.sdk.CreateNanoTDF(enc, bytes.NewReader(b), *nanoTDFConfig)
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
