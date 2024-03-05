package handlers

import (
	common "github.com/opentdf/platform/protocol/go/common"
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

// ListKasRegistryEntries lists the KeyAccessServer entries  in the project.
func (h Handler) ListKasRegistryEntries() ([]*kasregistry.KeyAccessServer, error) {
	// Create a request to list the KeyAccessServer entries.
	req := &kasregistry.ListKeyAccessServersRequest{}

	// List the KeyAccessServer entries  using the SDK.
	resp, err := h.sdk.KeyAccessServerRegistry.ListKeyAccessServers(h.ctx, req)
	if err != nil {
		return nil, err
	}

	// Return the list of KeyAccessServer entries.
	return resp.KeyAccessServers, nil
}

// CreateKasRegistryEntry creates a KeyAccessServer entry  in the project.
// map[string]interface{} used to handle arbitarily structured metadata
func (h Handler) CreateKasRegistryEntry(uri string, publicKey *kasregistry.PublicKey, metadata *common.MetadataMutable) (*kasregistry.KeyAccessServer, error) {
	// Create a request to create a KeyAccessServer entry.
	req := &kasregistry.CreateKeyAccessServerRequest{
		KeyAccessServer: &kasregistry.KeyAccessServerCreateUpdate{
			Uri:       uri,
			PublicKey: publicKey,
			Metadata:  metadata,
		},
	}

	// Create the KeyAccessServer entry using the SDK.
	resp, err := h.sdk.KeyAccessServerRegistry.CreateKeyAccessServer(h.ctx, req)
	if err != nil {
		return nil, err
	}

	// Return the created KeyAccessServer entry.
	return resp.KeyAccessServer, nil
}

// UpdateKasRegistryEntry updates a KeyAccessServer  entry  in the project.
// note: we are specifically building the request on the otherside, due to so manu of the options being optional
func (h Handler) UpdateKasRegistryEntry(id string, kasUpdateReq *kasregistry.UpdateKeyAccessServerRequest) (*kasregistry.KeyAccessServer, error) {

	// Update the KeyAccessServer entry using the SDK.
	resp, err := h.sdk.KeyAccessServerRegistry.UpdateKeyAccessServer(h.ctx, kasUpdateReq)
	if err != nil {
		return nil, err
	}

	// Return the updated KeyAccess Server entry.
	return resp.KeyAccessServer, nil
}

// DeleteKasRegistryEntry deletes a KeyAccessServer entry  from the project.
func (h Handler) DeleteKasRegistryEntry(id string) error {
	// Create a request to delete a KeyAccessServer entry.
	req := &kasregistry.DeleteKeyAccessServerRequest{
		Id: id,
	}

	// Delete  the KeyAccessServer entry using the SDK.
	_, err := h.sdk.KeyAccessServerRegistry.DeleteKeyAccessServer(h.ctx, req)
	if err != nil {
		return err
	}

	return nil
}
