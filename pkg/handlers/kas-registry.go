package handlers

import (
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/kasregistry"
)

func (h Handler) GetKasRegistryEntry(id string) (*policy.KeyAccessServer, error) {
	resp, err := h.sdk.KeyAccessServerRegistry.GetKeyAccessServer(h.ctx, &kasregistry.GetKeyAccessServerRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetKeyAccessServer(), nil
}

func (h Handler) ListKasRegistryEntries() ([]*policy.KeyAccessServer, error) {
	resp, err := h.sdk.KeyAccessServerRegistry.ListKeyAccessServers(h.ctx, &kasregistry.ListKeyAccessServersRequest{})
	if err != nil {
		return nil, err
	}

	return resp.GetKeyAccessServers(), nil
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
func (h Handler) UpdateKasRegistryEntry(id, uri, name string, publickey *policy.PublicKey, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.KeyAccessServer, error) {
	_, err := h.sdk.KeyAccessServerRegistry.UpdateKeyAccessServer(h.ctx, &kasregistry.UpdateKeyAccessServerRequest{
		Id:                     id,
		Uri:                    uri,
		Name:                   name,
		PublicKey:              publickey,
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
