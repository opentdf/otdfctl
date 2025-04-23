package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
	"github.com/opentdf/platform/protocol/go/policy/unsafe"
)

func (h Handler) GetNamespace(identifier string) (*policy.Namespace, error) {
	req := &namespaces.GetNamespaceRequest{
		Identifier: &namespaces.GetNamespaceRequest_NamespaceId{
			NamespaceId: identifier,
		},
	}
	if _, err := uuid.Parse(identifier); err != nil {
		req.Identifier = &namespaces.GetNamespaceRequest_Fqn{
			Fqn: identifier,
		}
	}

	resp, err := h.sdk.Namespaces.GetNamespace(h.ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetNamespace(), nil
}

func (h Handler) ListNamespaces(state common.ActiveStateEnum, limit, offset int32) ([]*policy.Namespace, *policy.PageResponse, error) {
	resp, err := h.sdk.Namespaces.ListNamespaces(h.ctx, &namespaces.ListNamespacesRequest{
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
func (h Handler) UnsafeReactivateNamespace(id string) (*policy.Namespace, error) {
	_, err := h.sdk.Unsafe.UnsafeReactivateNamespace(h.ctx, &unsafe.UnsafeReactivateNamespaceRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(id)
}

// Deletes and returns the deleted namespace
func (h Handler) UnsafeDeleteNamespace(id string, fqn string) error {
	_, err := h.sdk.Unsafe.UnsafeDeleteNamespace(h.ctx, &unsafe.UnsafeDeleteNamespaceRequest{
		Id:  id,
		Fqn: fqn,
	})
	return err
}

// Unsafely updates the namespace and returns the renamed namespace
func (h Handler) UnsafeUpdateNamespace(id, name string) (*policy.Namespace, error) {
	_, err := h.sdk.Unsafe.UnsafeUpdateNamespace(h.ctx, &unsafe.UnsafeUpdateNamespaceRequest{
		Id:   id,
		Name: name,
	})
	if err != nil {
		return nil, err
	}

	return h.GetNamespace(id)
}

// AssignKeyToAttributeNamespace assigns a KAS key to an attribute namespace
func (h *Handler) AssignKeyToAttributeNamespace(ctx context.Context, namespace, keyId string) (*namespaces.NamespaceKey, error) {
	namespaceKey := &namespaces.NamespaceKey{
		KeyId:       keyId,
		NamespaceId: namespace,
	}

	if _, err := uuid.Parse(namespace); err != nil {
		ns, err := h.GetNamespace(namespace)
		if err != nil {
			return nil, err
		}
		namespaceKey.NamespaceId = ns.GetId()
	}

	resp, err := h.sdk.Namespaces.AssignPublicKeyToNamespace(ctx, &namespaces.AssignPublicKeyToNamespaceRequest{
		NamespaceKey: namespaceKey,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetNamespaceKey(), nil
}

// RemoveKeyFromAttributeNamespace removes a KAS key from an attribute namespace
func (h *Handler) RemoveKeyFromAttributeNamespace(ctx context.Context, namespace, keyId string) error {
	namespaceKey := &namespaces.NamespaceKey{
		KeyId:       keyId,
		NamespaceId: namespace,
	}

	if _, err := uuid.Parse(namespace); err != nil {
		ns, err := h.GetNamespace(namespace)
		if err != nil {
			return err
		}
		namespaceKey.NamespaceId = ns.GetId()
	}

	_, err := h.sdk.Namespaces.RemovePublicKeyFromNamespace(ctx, &namespaces.RemovePublicKeyFromNamespaceRequest{
		NamespaceKey: namespaceKey,
	})
	return err
}
