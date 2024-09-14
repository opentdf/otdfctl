package handlers

import (
	"connectrpc.com/connect"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
	"github.com/opentdf/platform/protocol/go/policy/unsafe"
)

func (h Handler) GetNamespace(id string) (*policy.Namespace, error) {
	resp, err := h.sdk.Namespaces.GetNamespace(h.ctx, &connect.Request[namespaces.GetNamespaceRequest]{
		Msg: &namespaces.GetNamespaceRequest{
			Id: id,
		}})
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetNamespace(), nil
}

func (h Handler) ListNamespaces(state common.ActiveStateEnum) ([]*policy.Namespace, error) {
	resp, err := h.sdk.Namespaces.ListNamespaces(h.ctx, &connect.Request[namespaces.ListNamespacesRequest]{
		Msg: &namespaces.ListNamespacesRequest{State: state}})
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetNamespaces(), nil
}

// Creates and returns the created n
func (h Handler) CreateNamespace(name string, metadata *common.MetadataMutable) (*policy.Namespace, error) {
	resp, err := h.sdk.Namespaces.CreateNamespace(h.ctx, &connect.Request[namespaces.CreateNamespaceRequest]{
		Msg: &namespaces.CreateNamespaceRequest{
			Name:     name,
			Metadata: metadata,
		}})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(resp.Msg.GetNamespace().GetId())
}

// Updates and returns the updated namespace
func (h Handler) UpdateNamespace(id string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.Namespace, error) {
	_, err := h.sdk.Namespaces.UpdateNamespace(h.ctx, &connect.Request[namespaces.UpdateNamespaceRequest]{
		Msg: &namespaces.UpdateNamespaceRequest{
			Id:                     id,
			Metadata:               metadata,
			MetadataUpdateBehavior: behavior,
		}})
	if err != nil {
		return nil, err
	}
	return h.GetNamespace(id)
}

// Deactivates and returns the deactivated namespace
func (h Handler) DeactivateNamespace(id string) (*policy.Namespace, error) {
	_, err := h.sdk.Namespaces.DeactivateNamespace(h.ctx, &connect.Request[namespaces.DeactivateNamespaceRequest]{
		Msg: &namespaces.DeactivateNamespaceRequest{
			Id: id,
		}})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(id)
}

// Reactivates and returns the reactivated namespace
func (h Handler) UnsafeReactivateNamespace(id string) (*policy.Namespace, error) {
	_, err := h.sdk.Unsafe.UnsafeReactivateNamespace(h.ctx, &connect.Request[unsafe.UnsafeReactivateNamespaceRequest]{
		Msg: &unsafe.UnsafeReactivateNamespaceRequest{
			Id: id,
		}})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(id)
}

// Deletes and returns the deleted namespace
func (h Handler) UnsafeDeleteNamespace(id string, fqn string) error {
	_, err := h.sdk.Unsafe.UnsafeDeleteNamespace(h.ctx, &connect.Request[unsafe.UnsafeDeleteNamespaceRequest]{
		Msg: &unsafe.UnsafeDeleteNamespaceRequest{
			Id:  id,
			Fqn: fqn,
		}})
	return err
}

// Unsafely updates the namespace and returns the renamed namespace
func (h Handler) UnsafeUpdateNamespace(id, name string) (*policy.Namespace, error) {
	_, err := h.sdk.Unsafe.UnsafeUpdateNamespace(h.ctx, &connect.Request[unsafe.UnsafeUpdateNamespaceRequest]{
		Msg: &unsafe.UnsafeUpdateNamespaceRequest{
			Id:   id,
			Name: name,
		}})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(id)
}
