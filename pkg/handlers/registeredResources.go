package handlers

import (
	"context"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/registeredresources"
)

type RegisteredResource struct {
	Id       string
	Name     string
	Metadata *common.MetadataMutable
}

func (h *Handler) CreateRegisteredResource(name string, metadata *common.MetadataMutable) (*policy.RegisteredResource, error) {
	res, err := h.sdk.RegisteredResources.CreateRegisteredResource(context.Background(), &registeredresources.CreateRegisteredResourceRequest{
		Name:     name,
		Metadata: metadata,
	})
	if err != nil {
		return nil, err
	}

	return h.GetRegisteredResource(res.GetRegisteredResource().GetId())
}

func (h *Handler) GetRegisteredResource(id string) (*policy.RegisteredResource, error) {
	res, err := h.sdk.RegisteredResources.GetRegisteredResource(context.Background(), &registeredresources.GetRegisteredResourceRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return res.GetRegisteredResource(), nil
}

func (h *Handler) ListRegisteredResources(ctx context.Context, limit, offset int32) ([]*policy.RegisteredResource, *policy.PageResponse, error) {
	res, err := h.sdk.RegisteredResources.ListRegisteredResources(ctx, &registeredresources.ListRegisteredResourcesRequest{
		Pagination: &policy.PageRequest{
			Limit:  limit,
			Offset: offset,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return res.GetRegisteredResources(), res.GetPagination(), nil
}

func (h *Handler) UpdateRegisteredResource(id, name string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.RegisteredResource, error) {
	_, err := h.sdk.RegisteredResources.UpdateRegisteredResource(context.Background(), &registeredresources.UpdateRegisteredResourceRequest{
		Id:                     id,
		Name:                   name,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	})
	if err != nil {
		return nil, err
	}

	return h.GetRegisteredResource(id)
}

func (h *Handler) DeleteRegisteredResource(id string) (*policy.RegisteredResource, error) {
	resp, err := h.sdk.RegisteredResources.DeleteRegisteredResource(context.Background(), &registeredresources.DeleteRegisteredResourceRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetRegisteredResource(), nil
}
