package handlers

import (
	"context"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/kasregistry"
)

func (h Handler) GetKasRegistryEntry(ctx context.Context, id string) (*policy.KeyAccessServer, error) {
	resp, err := h.sdk.KeyAccessServerRegistry.GetKeyAccessServer(ctx, &kasregistry.GetKeyAccessServerRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetKeyAccessServer(), nil
}

func (h Handler) ListKasRegistryEntries(ctx context.Context, limit, offset int32) ([]*policy.KeyAccessServer, *policy.PageResponse, error) {
	resp, err := h.sdk.KeyAccessServerRegistry.ListKeyAccessServers(ctx, &kasregistry.ListKeyAccessServersRequest{
		Pagination: &policy.PageRequest{
			Limit:  limit,
			Offset: offset,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return resp.GetKeyAccessServers(), resp.GetPagination(), nil
}

// Creates the KAS registry and then returns the KAS
func (h Handler) CreateKasRegistryEntry(ctx context.Context, uri string, publicKey *policy.PublicKey, name string, metadata *common.MetadataMutable) (*policy.KeyAccessServer, error) {
	req := &kasregistry.CreateKeyAccessServerRequest{
		Uri:       uri,
		PublicKey: publicKey,
		Name:      name,
		Metadata:  metadata,
	}

	resp, err := h.sdk.KeyAccessServerRegistry.CreateKeyAccessServer(ctx, req)
	if err != nil {
		return nil, err
	}

	return h.GetKasRegistryEntry(ctx, resp.GetKeyAccessServer().GetId())
}

// Updates the KAS registry and then returns the KAS
func (h Handler) UpdateKasRegistryEntry(ctx context.Context, id, uri, name string, pubKey *policy.PublicKey, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.KeyAccessServer, error) {
	_, err := h.sdk.KeyAccessServerRegistry.UpdateKeyAccessServer(ctx, &kasregistry.UpdateKeyAccessServerRequest{
		Id:                     id,
		Uri:                    uri,
		Name:                   name,
		PublicKey:              pubKey,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	})
	if err != nil {
		return nil, err
	}

	return h.GetKasRegistryEntry(ctx, id)
}

// Deletes the KAS registry and returns the deleted KAS
func (h Handler) DeleteKasRegistryEntry(ctx context.Context, id string) (*policy.KeyAccessServer, error) {
	req := &kasregistry.DeleteKeyAccessServerRequest{
		Id: id,
	}

	resp, err := h.sdk.KeyAccessServerRegistry.DeleteKeyAccessServer(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetKeyAccessServer(), nil
}
