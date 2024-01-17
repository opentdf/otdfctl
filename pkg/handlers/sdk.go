package handlers

import (
	"context"

	"github.com/opentdf/opentdf-v2-poc/sdk"
)

var SDK *sdk.SDK

type Handler struct {
	sdk *sdk.SDK
	ctx context.Context
}

func New(platformEndpoint string) (Handler, error) {
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
