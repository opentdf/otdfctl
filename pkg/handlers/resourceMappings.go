package handlers

import (
	"context"

	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/resourcemapping"
)

type ResourceMapping struct {
	Id          string
	AttributeId string
	Terms       []string
}

func (h *Handler) CreateResourceMapping(attributeId string, terms []string) (*policy.ResourceMapping, error) {
	res, err := h.sdk.ResourceMapping.CreateResourceMapping(context.Background(), &resourcemapping.CreateResourceMappingRequest{
		AttributeValueId: attributeId,
		Terms:            terms,
	})
	if err != nil {
		return nil, err
	}

	return res.ResourceMapping, nil
}

func (h *Handler) GetResourceMapping(id string) (*policy.ResourceMapping, error) {
	res, err := h.sdk.ResourceMapping.GetResourceMapping(context.Background(), &resourcemapping.GetResourceMappingRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return res.ResourceMapping, nil
}

func (h *Handler) ListResourceMappings() ([]*policy.ResourceMapping, error) {
	res, err := h.sdk.ResourceMapping.ListResourceMappings(context.Background(), &resourcemapping.ListResourceMappingsRequest{})
	if err != nil {
		return nil, err
	}

	return res.ResourceMappings, nil
}

// TODO: verify updation behavior
func (h *Handler) UpdateResourceMapping(id string, attrValueId string, terms []string) (*policy.ResourceMapping, error) {
	res, err := h.sdk.ResourceMapping.UpdateResourceMapping(context.Background(), &resourcemapping.UpdateResourceMappingRequest{
		Id:               id,
		AttributeValueId: attrValueId,
		Terms:            terms,
	})
	if err != nil {
		return nil, err
	}

	return res.ResourceMapping, nil
}

func (h *Handler) DeleteResourceMapping(id string) (*policy.ResourceMapping, error) {
	resp, err := h.sdk.ResourceMapping.DeleteResourceMapping(context.Background(), &resourcemapping.DeleteResourceMappingRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.ResourceMapping, nil
}
