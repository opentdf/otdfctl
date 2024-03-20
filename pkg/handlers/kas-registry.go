package handlers

import (
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/kasregistry"
)

func (h Handler) GetKasRegistryEntry(id string) (*kasregistry.KeyAccessServer, error) {
	resp, err := h.sdk.KeyAccessServerRegistry.GetKeyAccessServer(h.ctx, &kasregistry.GetKeyAccessServerRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.KeyAccessServer, nil
}

func (h Handler) ListKasRegistryEntries() ([]*kasregistry.KeyAccessServer, error) {
	req := &kasregistry.ListKeyAccessServersRequest{}

	resp, err := h.sdk.KeyAccessServerRegistry.ListKeyAccessServers(h.ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.KeyAccessServers, nil
}

func (h Handler) CreateKasRegistryEntry(uri string, publicKey *kasregistry.PublicKey, metadata *common.MetadataMutable) (*kasregistry.KeyAccessServer, error) {
	req := &kasregistry.CreateKeyAccessServerRequest{
		Uri:       uri,
		PublicKey: publicKey,
		Metadata:  metadata,
	}

	resp, err := h.sdk.KeyAccessServerRegistry.CreateKeyAccessServer(h.ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.KeyAccessServer, nil
}

// TODO: verify updation behavior
func (h Handler) UpdateKasRegistryEntry(id string, kasUpdateReq *kasregistry.UpdateKeyAccessServerRequest) (*kasregistry.KeyAccessServer, error) {
	resp, err := h.sdk.KeyAccessServerRegistry.UpdateKeyAccessServer(h.ctx, kasUpdateReq)
	if err != nil {
		return nil, err
	}

	return resp.KeyAccessServer, nil
}

func (h Handler) DeleteKasRegistryEntry(id string) error {
	req := &kasregistry.DeleteKeyAccessServerRequest{
		Id: id,
	}

	_, err := h.sdk.KeyAccessServerRegistry.DeleteKeyAccessServer(h.ctx, req)
	if err != nil {
		return err
	}

	return nil
}
