package handlers

import (
	"bytes"
	"io"
)

func (h Handler) EncryptNanoBytes(b []byte, values []string) (*bytes.Buffer, error) {
	var encrypted []byte
	enc := bytes.NewBuffer(encrypted)

	nanoTDFCOnfig, err := h.sdk.NewNanoTDFConfig()
	if err != nil {
		return nil, err
	}

	nanoTDFCOnfig.SetKasUrl(h.platformEndpoint)
	nanoTDFCOnfig.SetAttributes(values)

	// TODO: validate values are FQNs or return an error [https://github.com/opentdf/platform/issues/515]
	_, err = h.sdk.CreateNanoTDF(enc, bytes.NewReader(b), *nanoTDFCOnfig)
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
