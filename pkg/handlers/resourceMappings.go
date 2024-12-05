package handlers

import (
	"context"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/resourcemapping"
)

type ResourceMapping struct {
	Id          string
	AttributeId string
	Terms       []string
}

// Creates and returns the created resource mapping
func (h *Handler) CreateResourceMapping(attributeId string, terms []string, metadata *common.MetadataMutable) (*policy.ResourceMapping, error) {
	res, err := h.sdk.ResourceMapping.CreateResourceMapping(context.Background(), &resourcemapping.CreateResourceMappingRequest{
		AttributeValueId: attributeId,
		Terms:            terms,
		Metadata:         metadata,
	})
	if err != nil {
		return nil, err
	}

	return h.GetResourceMapping(res.GetResourceMapping().GetId())
}

func (h *Handler) GetResourceMapping(id string) (*policy.ResourceMapping, error) {
	res, err := h.sdk.ResourceMapping.GetResourceMapping(context.Background(), &resourcemapping.GetResourceMappingRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return res.GetResourceMapping(), nil
}

func (h *Handler) ListResourceMappings(ctx context.Context, limit, offset int32) ([]*policy.ResourceMapping, *policy.PageResponse, error) {
	res, err := h.sdk.ResourceMapping.ListResourceMappings(ctx, &resourcemapping.ListResourceMappingsRequest{
		Pagination: &policy.PageRequest{
			Limit:  limit,
			Offset: offset,
		},
	})
	if err != nil {
		return nil, nil, err
	}

	return res.GetResourceMappings(), res.GetPagination(), nil
}

// TODO: verify updation behavior
// Updates and returns the updated resource mapping
func (h *Handler) UpdateResourceMapping(id string, attrValueId string, terms []string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.ResourceMapping, error) {
	_, err := h.sdk.ResourceMapping.UpdateResourceMapping(context.Background(), &resourcemapping.UpdateResourceMappingRequest{
		Id:                     id,
		AttributeValueId:       attrValueId,
		Terms:                  terms,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	})
	if err != nil {
		return nil, err
	}

	return h.GetResourceMapping(id)
}

func (h *Handler) DeleteResourceMapping(id string) (*policy.ResourceMapping, error) {
	resp, err := h.sdk.ResourceMapping.DeleteResourceMapping(context.Background(), &resourcemapping.DeleteResourceMappingRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetResourceMapping(), nil
}
