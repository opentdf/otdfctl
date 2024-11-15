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

func (h Handler) EncryptBytes(tdfType string, unencrypted []byte, attrValues []string, mimeType string, kasUrlPath string, ecdsaBinding bool, assertions string) (*bytes.Buffer, error) {
	var encrypted []byte
	enc := bytes.NewBuffer(encrypted)

	switch tdfType {
	// Encrypt the data as a ZTDF
	case "", TDF_TYPE_TDF3, TDF_TYPE_ZTDF:
		if ecdsaBinding {
			return nil, errors.New("ECDSA policy binding is not supported for ZTDF")
		}

		opts := []sdk.TDFOption{
			sdk.WithDataAttributes(attrValues...),
			sdk.WithKasInformation(sdk.KASInfo{
				URL: h.platformEndpoint + kasUrlPath,
			}),
			sdk.WithMimeType(mimeType),
		}

		var assertionConfigs []sdk.AssertionConfig
		if assertions != "" {
			err := json.Unmarshal([]byte(assertions), &assertionConfigs)
			if err != nil {
				return nil, errors.Join(ErrTDFUnableToReadAssertions, err)
			}
			opts = append(opts, sdk.WithAssertions(assertionConfigs...))
		}

		_, err := h.sdk.CreateTDF(enc, bytes.NewReader(unencrypted), opts...)
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
		if err = nanoTDFConfig.SetAttributes(attrValues); err != nil {
			return nil, err
		}
		// enable ECDSA policy binding
		if ecdsaBinding {
			nanoTDFConfig.EnableECDSAPolicyBinding()
		}
		// create the nano TDF
		if _, err = h.sdk.CreateNanoTDF(enc, bytes.NewReader(unencrypted), *nanoTDFConfig); err != nil {
			return nil, err
		}
		return enc, nil
	default:
		return nil, errors.New("unknown TDF type")
	}
}

func (h Handler) DecryptBytes(toDecrypt []byte, disableAssertionCheck bool) (*bytes.Buffer, error) {
	out := &bytes.Buffer{}
	pt := io.Writer(out)
	ec := bytes.NewReader(toDecrypt)
	switch sdk.GetTdfType(ec) {
	case sdk.Nano:
		if _, err := h.sdk.ReadNanoTDF(pt, ec); err != nil {
			return nil, err
		}
	case sdk.Standard:
		r, err := h.sdk.LoadTDF(ec, sdk.WithDisableAssertionVerification(disableAssertionCheck))
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

// TODO: Rename. Not sure what this value is at present
const inspectTDFEighteen = 18

func (h Handler) InspectTDF(toInspect []byte) (TDFInspect, []error) {
	b := bytes.NewReader(toInspect)
	switch sdk.GetTdfType(b) {
	case sdk.Standard:
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
	case sdk.Nano:
		header, size, err := sdk.NewNanoTDFHeaderFromReader(b)
		if err != nil {
			return TDFInspect{}, []error{errors.Join(ErrTDFInspectFailNotValidTDF, err)}
		}
		r := TDFInspect{
			NanoHeader: &header,
		}
		remainder := len(toInspect) - int(size)
		if remainder < inspectTDFEighteen {
			return r, []error{ErrTDFInspectFailNotValidTDF}
		}
		return r, nil
	case sdk.Invalid:
		return TDFInspect{}, []error{ErrTDFInspectFailNotValidTDF}
	default:
		return TDFInspect{}, []error{fmt.Errorf("tdf format unrecognized")}
	}
}
