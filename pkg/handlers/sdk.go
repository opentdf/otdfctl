package handlers

import (
	"context"
	"errors"

	"github.com/opentdf/otdfctl/pkg/auth"
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
	OIDC_TOKEN       string
	platformEndpoint string
}

type handlerOpts struct {
	endpoint    string
	tlsNoVerify bool

	profile *profiles.ProfileStore

	sdkOpts []sdk.Option
}

type handlerOptsFunc func(handlerOpts) handlerOpts

func WithEndpoint(endpoint string, tlsNoVerify bool) handlerOptsFunc {
	return func(c handlerOpts) handlerOpts {
		c.endpoint = endpoint
		c.tlsNoVerify = tlsNoVerify
		return c
	}
}

func WithProfile(profile *profiles.ProfileStore) handlerOptsFunc {
	return func(c handlerOpts) handlerOpts {
		c.profile = profile
		c.endpoint = profile.GetEndpoint()
		c.tlsNoVerify = profile.GetTLSNoVerify()

		// get sdk opts
		opts, err := auth.GetSDKAuthOptionFromProfile(profile)
		if err != nil {
			return c
		}
		c.sdkOpts = append(c.sdkOpts, opts)

		return c
	}
}

func WithSDKOpts(opts ...sdk.Option) handlerOptsFunc {
	return func(c handlerOpts) handlerOpts {
		c.sdkOpts = opts
		return c
	}
}

// Creates a new handler wrapping the SDK, which is authenticated through the cached client-credentials flow tokens
func New(opts ...handlerOptsFunc) (Handler, error) {
	var o handlerOpts
	for _, f := range opts {
		o = f(o)
	}

	u, err := utils.NormalizeEndpoint(o.endpoint)
	if err != nil {
		return Handler{}, err
	}

	if o.tlsNoVerify {
		o.sdkOpts = append(o.sdkOpts, sdk.WithInsecureSkipVerifyConn())
	}

	// TODO let's make sure we still support plaintext connections

	// get auth
	ao, err := auth.GetSDKAuthOptionFromProfile(o.profile)
	if err != nil {
		return Handler{}, err
	}
	o.sdkOpts = append(o.sdkOpts, ao)

	if u.Scheme == "http" {
		o.sdkOpts = append(o.sdkOpts, sdk.WithInsecurePlaintextConn())
	}

	s, err := sdk.New(u.Host, o.sdkOpts...)
	if err != nil {
		return Handler{}, err
	}

	return Handler{
		sdk:              s,
		platformEndpoint: o.endpoint,
		ctx:              context.Background(),
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

// func buildMetadata(metadata *common.MetadataMutable, fns ...func(*common.MetadataMutable) *common.MetadataMutable) *common.MetadataMutable {
// 	if metadata == nil {
// 		metadata = &common.MetadataMutable{}
// 	}
// 	if len(fns) == 0 {
// 		return metadata
// 	}
// 	for _, fn := range fns {
// 		metadata = fn(metadata)
// 	}
// 	return metadata
// }
