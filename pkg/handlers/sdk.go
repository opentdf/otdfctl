package handlers

import (
	"context"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/sdk"
)

var SDK *sdk.SDK

type Handler struct {
	sdk        *sdk.SDK
	ctx        context.Context
	OIDC_TOKEN string
}

func New(platformEndpoint string) (Handler, error) {
	// define the scopes in an array
	// scopes := []string{"email"}
	// normally, we should try to retrieve an active OICD token here, however, the SDK has no option for passing a token
	// so instead, we'll check if we have a clientId and clientSecret stored, and if so, we'll use those to init the SDK, otherwise, we'll use the insecure connection (which will stop working once we enforce auth on the backend)
	// clientSecret, clientId, err := GetClientIdAndSecretFromCache()

	// if err != nil {
	// 	return Handler{}, err
	// }

	// scopes := []string{"email"}
	// create the sdk with the client credentials
	//NOTE FROM AVERY: The below line is commented out because although it should work, the SDK
	// is having trouble with the "WithClientCredentials" endpoint
	// so although the commented out line should work, and will work in the future, today it doesn't, so
	// to facilitate development, we're leaving it commented, until the SDK is fixed, and using the insecure connection instead
	// note that for now we're hard coding the TOKEN_URL until we have an endpoint to get the config from
	// sdk, err := sdk.New(platformEndpoint, sdk.WithClientCredentials(clientId, clientSecret, scopes), sdk.WithTokenEndpoint(TOKEN_URL))
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
			Labels: metadata.Labels,
		}
		return nextMetadata
	}
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
