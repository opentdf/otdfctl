package handlers

import (
	"bytes"
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
	minBytesLength                        = 3
)

const (
	TDF_TYPE_ZTDF = "ztdf"
	TDF_TYPE_TDF3 = "tdf3" // alias for TDF
	TDF_TYPE_NANO = "nano"
)

type TDFInspect struct {
	NanoHeader          *sdk.NanoTDFHeader
	ZTDFManifest        *sdk.Manifest
	Attributes          []string
	UnencryptedMetadata []byte
}

func (h Handler) EncryptBytes(tdfType string, b []byte, values []string, mimeType string, kasUrlPath string, ecdsaBinding bool) (*bytes.Buffer, error) {
	var encrypted []byte
	enc := bytes.NewBuffer(encrypted)

	switch tdfType {
	// Encrypt the data as a ZTDF
	case "", TDF_TYPE_TDF3, TDF_TYPE_ZTDF:
		if ecdsaBinding {
			return nil, errors.New("ECDSA policy binding is not supported for ZTDF")
		}

		_, err := h.sdk.CreateTDF(enc, bytes.NewReader(b),
			sdk.WithDataAttributes(values...),
			sdk.WithKasInformation(sdk.KASInfo{
				URL: h.platformEndpoint + kasUrlPath,
			}),
			sdk.WithMimeType(mimeType),
		)
		return enc, err

	// Encrypt the data as a Nano TDF
	case TDF_TYPE_NANO:
		nanoTDFConfig, err := h.sdk.NewNanoTDFConfig()
		if err != nil {
			return nil, err
		}
		// set the KAS URL
		if err = nanoTDFConfig.SetKasURL(h.platformEndpoint + kasUrlPath); err != nil {
			return nil, err
		}
		// set the attributes
		if err = nanoTDFConfig.SetAttributes(values); err != nil {
			return nil, err
		}
		// enable ECDSA policy binding
		if ecdsaBinding {
			nanoTDFConfig.EnableECDSAPolicyBinding()
		}
		// create the nano TDF
		if _, err = h.sdk.CreateNanoTDF(enc, bytes.NewReader(b), *nanoTDFConfig); err != nil {
			return nil, err
		}
		return enc, nil
	default:
		return nil, errors.New("unknown TDF type")
	}
}

func (h Handler) DecryptBytes(toDecrypt []byte) (*bytes.Buffer, error) {
	out := &bytes.Buffer{}
	pt := io.Writer(out)
	ec := bytes.NewReader(toDecrypt)
	tdfType := sdk.GetTdfType(ec)
	// reset the reader to the beginning
	ec.Reset(toDecrypt)
	switch tdfType {
	case sdk.Nano:
		if _, err := h.sdk.ReadNanoTDF(pt, ec); err != nil {
			return nil, err
		}
	case sdk.Standard:
		r, err := h.sdk.LoadTDF(ec)
		if err != nil {
			return nil, err
		}
		//nolint:errorlint // callers intended to test error equality directly
		if _, err = io.Copy(pt, r); err != nil && err != io.EOF {
			return nil, err
		}
	case sdk.Invalid:
		return nil, errors.New("invalid TDF")
	default:
		return nil, errors.New("unknown TDF type")
	}
	return out, nil
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
