package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/opentdf/platform/sdk"
)

var (
	ErrTDFInspectFailNotValidTDF          = errors.New("file or input is not a valid TDF")
	ErrTDFInspectFailNotInspectable       = errors.New("file or input is not inspectable")
	ErrTDFUnableToReadAttributes          = errors.New("unable to read attributes from TDF")
	ErrTDFUnableToReadUnencryptedMetadata = errors.New("unable to read unencrypted metadata from TDF")
	ErrTDFUnableToReadAssertions          = errors.New("unable to read assertions")
	minBytesLength                        = 3
)

func (h Handler) EncryptBytes(b []byte, values []string, mimeType string, kasUrlPath string, assertions string) (*bytes.Buffer, error) {
	var encrypted []byte
	enc := bytes.NewBuffer(encrypted)

	var assertionConfigs []sdk.AssertionConfig
	if assertions != "" {
		err := json.Unmarshal([]byte(assertions), &assertionConfigs)
		if err != nil {
			return nil, errors.Join(ErrTDFUnableToReadAssertions, err)
		}
	}

	// TODO: validate values are FQNs or return an error [https://github.com/opentdf/platform/issues/515]
	_, err := h.sdk.CreateTDF(enc, bytes.NewReader(b),
		sdk.WithDataAttributes(values...),
		sdk.WithKasInformation(sdk.KASInfo{
			URL: h.platformEndpoint + kasUrlPath,
		}),
		sdk.WithAssertions(assertionConfigs...),
		sdk.WithMimeType(mimeType),
	)
	if err != nil {
		return nil, err
	}
	return enc, nil
}

func (h Handler) DecryptTDF(toDecrypt []byte, disableAssertionVerification bool) (*bytes.Buffer, error) {
	tdfreader, err := h.sdk.LoadTDF(bytes.NewReader(toDecrypt),
		sdk.WithDisableAssertionVerification(disableAssertionVerification),
	)
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, tdfreader)
	//nolint:errorlint // callers intended to test error equality directly
	if err != nil && err != io.EOF {
		return nil, err
	}
	return buf, nil
}

type TDFInspect struct {
	NanoHeader          *sdk.NanoTDFHeader
	ZTDFManifest        *sdk.Manifest
	Attributes          []string
	UnencryptedMetadata []byte
}

func (h Handler) InspectTDF(toInspect []byte) (TDFInspect, []error) {
	if len(toInspect) < minBytesLength {
		return TDFInspect{}, []error{fmt.Errorf("tdf too small [%d] bytes", len(toInspect))}
	}
	switch {
	case bytes.Equal([]byte("PK"), toInspect[0:2]):
		return h.InspectZTDF(toInspect)
	case bytes.Equal([]byte("L1L"), toInspect[0:3]):
		return h.InspectNanoTDF(toInspect)
	}
	return TDFInspect{}, []error{fmt.Errorf("tdf format unrecognized")}
}

func (h Handler) InspectZTDF(toInspect []byte) (TDFInspect, []error) {
	// grouping errors so we don't impact the piping of the data
	errs := []error{}

	tdfreader, err := h.sdk.LoadTDF(bytes.NewReader(toInspect))
	if err != nil {
		if strings.Contains(err.Error(), "zip: not a valid zip file") {
			return TDFInspect{}, []error{ErrTDFInspectFailNotInspectable}
		}
		return TDFInspect{}, []error{errors.Join(ErrTDFInspectFailNotValidTDF, err)}
	}

	attributes, err := tdfreader.DataAttributes()
	if err != nil {
		errs = append(errs, errors.Join(ErrTDFUnableToReadAttributes, err))
	}

	unencryptedMetadata, err := tdfreader.UnencryptedMetadata()
	if err != nil {
		errs = append(errs, errors.Join(ErrTDFUnableToReadUnencryptedMetadata, err))
	}

	m := tdfreader.Manifest()
	return TDFInspect{
		ZTDFManifest:        &m,
		Attributes:          attributes,
		UnencryptedMetadata: unencryptedMetadata,
	}, errs
}

//nolint:gosec,mnd // SDK should secure lengths of inputs/outputs
func (h Handler) InspectNanoTDF(toInspect []byte) (TDFInspect, []error) {
	header, size, err := sdk.NewNanoTDFHeaderFromReader(bytes.NewReader(toInspect))
	if err != nil {
		return TDFInspect{}, []error{errors.Join(ErrTDFInspectFailNotValidTDF, err)}
	}
	r := TDFInspect{
		NanoHeader: &header,
	}
	remainder := uint32(len(toInspect)) - size
	if remainder < 18 {
		return r, []error{ErrTDFInspectFailNotValidTDF}
	}
	return r, nil
}
