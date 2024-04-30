package handlers

import (
	"context"
	"errors"

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
func New(platformEndpoint string) (Handler, error) {
	scopes := []string{"email"}
	clientSecret, clientId, err := GetClientIdAndSecretFromCache()
	if err != nil {
		return Handler{}, err
	}

	sdk, err := sdk.New(platformEndpoint, sdk.WithClientCredentials(clientId, clientSecret, scopes), sdk.WithTokenEndpoint(TOKEN_URL), sdk.WithInsecureConn())
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
