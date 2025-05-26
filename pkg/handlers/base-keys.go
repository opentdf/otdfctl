package handlers

import (
	"context"
	"errors"

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

func (h Handler) SetBaseKey(ctx context.Context, id string, key *kasregistry.KasKeyIdentifier) (*kasregistry.SetBaseKeyResponse, error) {
	req := kasregistry.SetBaseKeyRequest{}
	switch {
	case id != "":
		req.ActiveKey = &kasregistry.SetBaseKeyRequest_Id{
			Id: id,
		}
	case key != nil:
		req.ActiveKey = &kasregistry.SetBaseKeyRequest_Key{
			Key: key,
		}
	default:
		return nil, errors.New("id or key must be provided")
	}

	return h.sdk.KeyAccessServerRegistry.SetBaseKey(ctx, &req)
}
