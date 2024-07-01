package handlers

import (
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
	"github.com/opentdf/platform/protocol/go/policy/unsafe"
)

func (h Handler) GetNamespace(id string) (*policy.Namespace, error) {
	resp, err := h.sdk.Namespaces.GetNamespace(h.ctx, &namespaces.GetNamespaceRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetNamespace(), nil
}

func (h Handler) ListNamespaces(state common.ActiveStateEnum) ([]*policy.Namespace, error) {
	resp, err := h.sdk.Namespaces.ListNamespaces(h.ctx, &namespaces.ListNamespacesRequest{State: state})
	if err != nil {
		return nil, err
	}

	return resp.GetNamespaces(), nil
}

// Creates and returns the created n
func (h Handler) CreateNamespace(name string, metadata *common.MetadataMutable) (*policy.Namespace, error) {
	resp, err := h.sdk.Namespaces.CreateNamespace(h.ctx, &namespaces.CreateNamespaceRequest{
		Name:     name,
		Metadata: metadata,
	})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(resp.GetNamespace().GetId())
}

// Updates and returns the updated namespace
func (h Handler) UpdateNamespace(id string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.Namespace, error) {
	_, err := h.sdk.Namespaces.UpdateNamespace(h.ctx, &namespaces.UpdateNamespaceRequest{
		Id:                     id,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	})
	if err != nil {
		return nil, err
	}
	return h.GetNamespace(id)
}

// Deactivates and returns the deactivated namespace
func (h Handler) DeactivateNamespace(id string) (*policy.Namespace, error) {
	_, err := h.sdk.Namespaces.DeactivateNamespace(h.ctx, &namespaces.DeactivateNamespaceRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(id)
}

// Reactivates and returns the reactivated namespace
func (h Handler) ReactivateNamespace(id string) (*policy.Namespace, error) {
	_, err := h.sdk.Unsafe.ReactivateNamespace(h.ctx, &unsafe.ReactivateNamespaceRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(id)
}

// Deletes and returns the deleted namespace
func (h Handler) DeleteNamespace(id string, fqn string) (error) {
	_, err := h.sdk.Unsafe.DeleteNamespace(h.ctx, &unsafe.DeleteNamespaceRequest{
		Id:  id,
		Fqn: fqn,
	})
	if err != nil {
		return err
	}

	return nil
}

// Unsafely updates the namespace and returns the renamed namespace
func (h Handler) UnsafeUpdateNamespace(id, name string) (*policy.Namespace, error) {
	_, err := h.sdk.Unsafe.UpdateNamespace(h.ctx, &unsafe.UpdateNamespaceRequest{
		Id:   id,
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(id)
}
