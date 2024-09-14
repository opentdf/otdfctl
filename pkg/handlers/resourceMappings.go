package handlers

import (
	"context"

	"connectrpc.com/connect"
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
	res, err := h.sdk.ResourceMapping.CreateResourceMapping(context.Background(), &connect.Request[resourcemapping.CreateResourceMappingRequest]{
		Msg: &resourcemapping.CreateResourceMappingRequest{
			AttributeValueId: attributeId,
			Terms:            terms,
			Metadata:         metadata,
		}})
	if err != nil {
		return nil, err
	}

	return h.GetResourceMapping(res.Msg.GetResourceMapping().GetId())
}

func (h *Handler) GetResourceMapping(id string) (*policy.ResourceMapping, error) {
	res, err := h.sdk.ResourceMapping.GetResourceMapping(context.Background(), &connect.Request[resourcemapping.GetResourceMappingRequest]{
		Msg: &resourcemapping.GetResourceMappingRequest{
			Id: id,
		}})
	if err != nil {
		return nil, err
	}

	return res.Msg.GetResourceMapping(), nil
}

func (h *Handler) ListResourceMappings() ([]*policy.ResourceMapping, error) {
	res, err := h.sdk.ResourceMapping.ListResourceMappings(context.Background(), &connect.Request[resourcemapping.ListResourceMappingsRequest]{
		Msg: &resourcemapping.ListResourceMappingsRequest{}})
	if err != nil {
		return nil, err
	}

	return res.Msg.GetResourceMappings(), nil
}

// TODO: verify updation behavior
// Updates and returns the updated resource mapping
func (h *Handler) UpdateResourceMapping(id string, attrValueId string, terms []string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.ResourceMapping, error) {
	_, err := h.sdk.ResourceMapping.UpdateResourceMapping(context.Background(), &connect.Request[resourcemapping.UpdateResourceMappingRequest]{
		Msg: &resourcemapping.UpdateResourceMappingRequest{
			Id:                     id,
			AttributeValueId:       attrValueId,
			Terms:                  terms,
			Metadata:               metadata,
			MetadataUpdateBehavior: behavior,
		}})
	if err != nil {
		return nil, err
	}

	return h.GetResourceMapping(id)
}

func (h *Handler) DeleteResourceMapping(id string) (*policy.ResourceMapping, error) {
	resp, err := h.sdk.ResourceMapping.DeleteResourceMapping(context.Background(), &connect.Request[resourcemapping.DeleteResourceMappingRequest]{
		Msg: &resourcemapping.DeleteResourceMappingRequest{
			Id: id,
		}})
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetResourceMapping(), nil
}
