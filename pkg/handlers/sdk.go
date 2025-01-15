package handlers

import (
	"context"
	"errors"

	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/sdk"
)

var (
	SDK *sdk.SDK

	ErrUnauthenticated = errors.New("unauthenticated")
)

type Handler struct {
	sdk *sdk.SDK
	//nolint:containedctx // need to handle in a separate dedicated issue [https://github.com/opentdf/otdfctl/issues/364]
	ctx              context.Context
	platformEndpoint string
	profile          *profiles.ProfileCLI
}

// Creates a new handler wrapping the SDK, which is authenticated through the cached client-credentials flow tokens
func New(ctx context.Context, p *profiles.ProfileCLI) (Handler, error) {
	sdkOpts, err := profiles.GetSDKOptionsFromProfile(p)
	if err != nil {
		return Handler{}, err
	}
	u, err := utils.NormalizeEndpoint(p.GetEndpoint())
	if err != nil {
		return Handler{}, err
	}

	// TODO let's make sure we still support plaintext connections

	if u.Scheme == "http" {
		sdkOpts = append(sdkOpts, sdk.WithInsecurePlaintextConn())
	}

	s, err := sdk.New(u.Host, sdkOpts...)
	if err != nil {
		return Handler{}, err
	}

	return Handler{
		sdk:              s,
		platformEndpoint: u.String(),
		profile:          p,
		ctx:              ctx,
	}, nil
}

func (h Handler) Close() error {
	return h.sdk.Close()
}

func (h Handler) Direct() *sdk.SDK {
	return h.sdk
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
		labels := metadata.GetLabels()
		labels[key] = value
		nextMetadata := &common.MetadataMutable{
			Labels: labels,
		}
		return nextMetadata
	}
}
