package handlers

import (
	"context"
	"errors"

	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
	"github.com/opentdf/platform/protocol/go/policy/unsafe"
	"google.golang.org/grpc/status"
)

func (h Handler) GetNamespace(identifier string) (*policy.Namespace, error) {
	nsReq := new(namespaces.GetNamespaceRequest)

	if utils.IsUUID(identifier) {
		nsReq.Identifier = &namespaces.GetNamespaceRequest_NamespaceId{
			NamespaceId: identifier,
		}
	} else {
		nsReq.Identifier = &namespaces.GetNamespaceRequest_Fqn{
			Fqn: identifier,
		}
	}

	resp, err := h.sdk.Namespaces.GetNamespace(h.ctx, nsReq)
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

func (h Handler) AddPublicKeyToNamespace(ctx context.Context, nameSpace, publicKeyID string) (*namespaces.NamespaceKey, error) {
	nk := &namespaces.NamespaceKey{
		KeyId: publicKeyID,
	}

	if utils.IsUUID(nameSpace) {
		nk.NamespaceId = nameSpace
	} else {
		nss, err := h.GetNamespace(nameSpace)
		if err != nil {
			return nil, err
		}
		nk.NamespaceId = nss.GetId()
	}

	resp, err := h.sdk.Namespaces.AssignKeyToNamespace(ctx, &namespaces.AssignKeyToNamespaceRequest{
		NamespaceKey: nk,
	})
	if err != nil {
		s := status.Convert(err)
		return nil, errors.New(s.Message())
	}

	return resp.GetNamespaceKey(), nil
}

func (h Handler) RemovePublicKeyFromNamespace(ctx context.Context, nameSpace, publicKeyID string) (*namespaces.NamespaceKey, error) {
	nk := &namespaces.NamespaceKey{
		KeyId: publicKeyID,
	}

	if utils.IsUUID(nameSpace) {
		nk.NamespaceId = nameSpace
	} else {
		nss, err := h.GetNamespace(nameSpace)
		if err != nil {
			return nil, err
		}
		nk.NamespaceId = nss.GetId()
	}
	_, err := h.sdk.Namespaces.RemoveKeyFromNamespace(ctx, &namespaces.RemoveKeyFromNamespaceRequest{
		NamespaceKey: nk,
	})
	if err != nil {
		s := status.Convert(err)
		return nil, errors.New(s.Message())
	}

	return nk, nil
}
