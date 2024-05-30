package handlers

import (
	"context"
	"errors"
	"net/url"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/sdk"
)

var (
	SDK *sdk.SDK

	ErrUnauthenticated = errors.New("unauthenticated")
)

type Handler struct {
	sdk              *sdk.SDK
	ctx              context.Context
	OIDC_TOKEN       string
	platformEndpoint string
}

// Creates a new handler wrapping the SDK, which is authenticated through the cached client-credentials flow tokens
func New(platformEndpoint, clientID, clientSecret string, tlsNoVerify bool) (Handler, error) {
	scopes := []string{"email"}

	opts := []sdk.Option{
		sdk.WithClientCredentials(clientID, clientSecret, scopes),
		sdk.WithTokenEndpoint(TOKEN_URL),
	}

	// Try an parse scheme out of platformEndpoint
	// If it fails, use the default scheme of https
	// There has to be a better way to do this
	platformURL, err := url.Parse(platformEndpoint)
	if err != nil {
		return Handler{}, err
	}

	switch platformURL.Scheme {
	case "http":
		opts = append(opts, sdk.WithInsecurePlaintextConn())
		if platformURL.Port() == "" {
			platformURL.Host += ":80"
		}
	case "https":
		if platformURL.Port() == "" {
			platformURL.Host += ":443"
		}
		if tlsNoVerify {
			opts = append(opts, sdk.WithInsecureSkipVerifyConn())
		}
	default:
		return Handler{}, errors.New("invalid scheme")
	}

	sdk, err := sdk.New(platformURL.Host, opts...)
	if err != nil {
		return Handler{}, err
	}

	return Handler{
		sdk:              sdk,
		platformEndpoint: platformEndpoint,
		ctx:              context.Background(),
	}, nil
}

func (h Handler) Close() error {
	return h.sdk.Close()
}

// Replace all labels in the metadata
func (h Handler) WithReplaceLabelsMetadata(metadata *common.MetadataMutable, labels map[string]string) func(*common.MetadataMutable) *common.MetadataMutable {
	return func(*common.MetadataMutable) *common.MetadataMutable {
		nextMetadata := &common.MetadataMutable{
			Labels: labels,
		}
		return nextMetadata
	}
}

// Append a label to the metadata
func (h Handler) WithLabelMetadata(metadata *common.MetadataMutable, key, value string) func(*common.MetadataMutable) *common.MetadataMutable {
	return func(*common.MetadataMutable) *common.MetadataMutable {
		labels := metadata.Labels
		labels[key] = value
		nextMetadata := &common.MetadataMutable{
			Labels: labels,
		}
		return nextMetadata
	}
}

func buildMetadata(metadata *common.MetadataMutable, fns ...func(*common.MetadataMutable) *common.MetadataMutable) *common.MetadataMutable {
	if metadata == nil {
		metadata = &common.MetadataMutable{}
	}
	if len(fns) == 0 {
		return metadata
	}
	for _, fn := range fns {
		metadata = fn(metadata)
	}
	return metadata
}
