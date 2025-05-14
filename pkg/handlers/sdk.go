package handlers

import (
	"context"
	"errors"
	"net"

	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/sdk"
	"golang.org/x/net/proxy"
	"google.golang.org/grpc"
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
	grpcProxy   string

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

func WithProxy(proxy string, parent handlerOptsFunc) handlerOptsFunc {
	if proxy == "" {
		return parent
	}
	return func(opts handlerOpts) handlerOpts {
		opts = parent(opts)
		opts.grpcProxy = proxy
		return opts
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

	// get auth
	authSDKOpt, err := auth.GetSDKAuthOptionFromProfile(o.profile)
	if err != nil {
		return Handler{}, err
	}
	o.sdkOpts = append(o.sdkOpts, authSDKOpt, sdk.WithConnectionValidation())

	if u.Scheme == "http" {
		o.sdkOpts = append(o.sdkOpts, sdk.WithInsecurePlaintextConn())
	}

	if o.grpcProxy != "" {
		o.sdkOpts = append(o.sdkOpts, sdk.WithExtraDialOptions(grpcSocks5Proxy(o.grpcProxy)))
	}

	s, err := sdk.New(u.String(), o.sdkOpts...)
	if err != nil {
		return Handler{}, err
	}

	return Handler{
		sdk:              s,
		platformEndpoint: o.endpoint,
		ctx:              context.Background(),
	}, nil
}

func grpcSocks5Proxy(proxyAddr string) grpc.DialOption {
	return grpc.WithContextDialer(func(ctx context.Context, addr string) (net.Conn, error) {
		// Create a SOCKS5 dialer
		dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
		if err != nil {
			return nil, err
		}

		// Use the dialer to connect to the target address
		return dialer.Dial("tcp", addr)
	})
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
