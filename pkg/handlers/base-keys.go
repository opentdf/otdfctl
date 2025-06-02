package handlers

import (
	"context"

	"github.com/opentdf/platform/protocol/go/policy/kasregistry"
)

// GetBaseKey retrieves a base key from the KAS registry.
// This is a stub function and needs to be implemented.
func (h Handler) GetBaseKey(ctx context.Context) (*kasregistry.SimpleKasKey, error) {
	resp, err := h.sdk.KeyAccessServerRegistry.GetBaseKey(ctx, &kasregistry.GetBaseKeyRequest{})
	if err != nil {
		return nil, err
	}

	return resp.GetBaseKey(), nil
}

func (h Handler) SetBaseKey(ctx context.Context, key *kasregistry.KasKeyIdentifier) (*kasregistry.SetBaseKeyResponse, error) {
	req := kasregistry.SetBaseKeyRequest{}

	req.ActiveKey = &kasregistry.SetBaseKeyRequest_Key{
		Key: key,
	}

	return h.sdk.KeyAccessServerRegistry.SetBaseKey(ctx, &req)
}
