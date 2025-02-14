package handlers

import (
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/kasregistry"
)

func (h Handler) GetKasRegistryEntry(identifier string) (*policy.KeyAccessServer, error) {
	req := &kasregistry.GetKeyAccessServerRequest{}

	switch {
	case utils.IsUUID(identifier):
		req.Identifier = &kasregistry.GetKeyAccessServerRequest_KasId{KasId: identifier}
	case utils.IsURI(identifier):
		req.Identifier = &kasregistry.GetKeyAccessServerRequest_Uri{Uri: identifier}
	default:
		// At this point we assume its a kas name
		req.Identifier = &kasregistry.GetKeyAccessServerRequest_Name{Name: identifier}
	}

	resp, err := h.sdk.KeyAccessServerRegistry.GetKeyAccessServer(h.ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetKeyAccessServer(), nil
}

func (h Handler) ListKasRegistryEntries(limit, offset int32) ([]*policy.KeyAccessServer, *policy.PageResponse, error) {
	resp, err := h.sdk.KeyAccessServerRegistry.ListKeyAccessServers(h.ctx, &kasregistry.ListKeyAccessServersRequest{
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
func (h Handler) CreateKasRegistryEntry(uri string, publicKey *policy.PublicKey, name string, metadata *common.MetadataMutable) (*policy.KeyAccessServer, error) {
	req := &kasregistry.CreateKeyAccessServerRequest{
		Uri:       uri,
		PublicKey: publicKey,
		Name:      name,
		Metadata:  metadata,
	}

	resp, err := h.sdk.KeyAccessServerRegistry.CreateKeyAccessServer(h.ctx, req)
	if err != nil {
		return nil, err
	}

	return h.GetKasRegistryEntry(resp.GetKeyAccessServer().GetId())
}

// Updates the KAS registry and then returns the KAS
func (h Handler) UpdateKasRegistryEntry(id, uri, name string, pubKey *policy.PublicKey, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.KeyAccessServer, error) {
	_, err := h.sdk.KeyAccessServerRegistry.UpdateKeyAccessServer(h.ctx, &kasregistry.UpdateKeyAccessServerRequest{
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

	return h.GetKasRegistryEntry(id)
}

// Deletes the KAS registry and returns the deleted KAS
func (h Handler) DeleteKasRegistryEntry(id string) (*policy.KeyAccessServer, error) {
	req := &kasregistry.DeleteKeyAccessServerRequest{
		Id: id,
	}

	resp, err := h.sdk.KeyAccessServerRegistry.DeleteKeyAccessServer(h.ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetKeyAccessServer(), nil
}
