package handlers

import (
	"context"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/sdk"
)

var SDK *sdk.SDK

type Handler struct {
	sdk *sdk.SDK
	ctx context.Context
}

func New(platformEndpoint string) (Handler, error) {
	// define the scopes in an array
	// scopes := []string{"email"}

	// sdk, err := sdk.New(platformEndpoint, sdk.WithClientCredentials("client-id", "clientSecret", scopes), sdk.WithTokenEndpoint("http://dummy/token-endpoint"))
	sdk, err := sdk.New(platformEndpoint, sdk.WithInsecureConn())
	if err != nil {
		return Handler{}, err
	}

	return Handler{
		sdk: sdk,
		ctx: context.Background(),
	}, nil
}

func (h Handler) Close() error {
	return h.sdk.Close()
}

// Replace the description in the metadata
func (h Handler) WithDescriptionMetadata(metadata *common.Metadata, description string) func() *common.Metadata {
	return func() *common.Metadata {
		nextMetadata := &common.Metadata{
			Labels:      metadata.Labels,
			Description: description,
		}
		return nextMetadata
	}
}

// Replace all labels in the metadata
func (h Handler) WithReplaceLabelsMetadata(metadata *common.MetadataMutable, labels map[string]string) func(*common.MetadataMutable) *common.MetadataMutable {
	return func(*common.MetadataMutable) *common.MetadataMutable {
		nextMetadata := &common.MetadataMutable{
			Labels:      labels,
			Description: metadata.Description,
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
			Labels:      labels,
			Description: metadata.Description,
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
