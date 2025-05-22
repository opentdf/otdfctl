package handlers

import (
	"context"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
	"github.com/opentdf/platform/protocol/go/policy/unsafe"
)

func (h Handler) GetNamespace(ctx context.Context, id string) (*policy.Namespace, error) {
	resp, err := h.sdk.Namespaces.GetNamespace(ctx, &namespaces.GetNamespaceRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetNamespace(), nil
}

func (h Handler) ListNamespaces(ctx context.Context, state common.ActiveStateEnum, limit, offset int32) ([]*policy.Namespace, *policy.PageResponse, error) {
	resp, err := h.sdk.Namespaces.ListNamespaces(ctx, &namespaces.ListNamespacesRequest{
		State: state,
		Pagination: &policy.PageRequest{
			Limit:  limit,
			Offset: offset,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return resp.GetNamespaces(), resp.GetPagination(), nil
}

// Creates and returns the created n
func (h Handler) CreateNamespace(ctx context.Context, name string, metadata *common.MetadataMutable) (*policy.Namespace, error) {
	resp, err := h.sdk.Namespaces.CreateNamespace(ctx, &namespaces.CreateNamespaceRequest{
		Name:     name,
		Metadata: metadata,
	})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(ctx, resp.GetNamespace().GetId())
}

// Updates and returns the updated namespace
func (h Handler) UpdateNamespace(ctx context.Context, id string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.Namespace, error) {
	_, err := h.sdk.Namespaces.UpdateNamespace(ctx, &namespaces.UpdateNamespaceRequest{
		Id:                     id,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	})
	if err != nil {
		return nil, err
	}
	return h.GetNamespace(ctx, id)
}

// Deactivates and returns the deactivated namespace
func (h Handler) DeactivateNamespace(ctx context.Context, id string) (*policy.Namespace, error) {
	_, err := h.sdk.Namespaces.DeactivateNamespace(ctx, &namespaces.DeactivateNamespaceRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(ctx, id)
}

// Reactivates and returns the reactivated namespace
func (h Handler) UnsafeReactivateNamespace(ctx context.Context, id string) (*policy.Namespace, error) {
	_, err := h.sdk.Unsafe.UnsafeReactivateNamespace(ctx, &unsafe.UnsafeReactivateNamespaceRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(ctx, id)
}

// Deletes and returns the deleted namespace
func (h Handler) UnsafeDeleteNamespace(ctx context.Context, id string, fqn string) error {
	_, err := h.sdk.Unsafe.UnsafeDeleteNamespace(ctx, &unsafe.UnsafeDeleteNamespaceRequest{
		Id:  id,
		Fqn: fqn,
	})
	return err
}

// Unsafely updates the namespace and returns the renamed namespace
func (h Handler) UnsafeUpdateNamespace(ctx context.Context, id, name string) (*policy.Namespace, error) {
	_, err := h.sdk.Unsafe.UnsafeUpdateNamespace(ctx, &unsafe.UnsafeUpdateNamespaceRequest{
		Id:   id,
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(ctx, id)
}
