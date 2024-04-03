package handlers

import (
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/attributes"
)

func (h *Handler) ListAttributeValues(state common.ActiveStateEnum) ([]*policy.Value, error) {
	resp, err := h.sdk.Attributes.ListAttributeValues(h.ctx, &attributes.ListAttributeValuesRequest{State: state})
	if err != nil {
		return nil, err
	}
	return resp.Values, err
}

// Creates and returns the created value
func (h *Handler) CreateAttributeValue(attributeId string, value string, metadata *common.MetadataMutable) (*policy.Value, error) {
	resp, err := h.sdk.Attributes.CreateAttributeValue(h.ctx, &attributes.CreateAttributeValueRequest{
		AttributeId: attributeId,
		Value:       value,
		Metadata:    metadata,
	})
	if err != nil {
		return nil, err
	}

	return h.GetAttributeValue(resp.GetValue().GetId())
}

func (h *Handler) GetAttributeValue(id string) (*policy.Value, error) {
	resp, err := h.sdk.Attributes.GetAttributeValue(h.ctx, &attributes.GetAttributeValueRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetValue(), nil
}

// Updates and returns updated value
func (h *Handler) UpdateAttributeValue(id string, memberIds []string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.Value, error) {
	resp, err := h.sdk.Attributes.UpdateAttributeValue(h.ctx, &attributes.UpdateAttributeValueRequest{
		Id:                     id,
		Members:                memberIds,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	})
	if err != nil {
		return nil, err
	}

	return h.GetAttributeValue(resp.GetValue().GetId())
}

// Deactivates and returns deactivated value
func (h *Handler) DeactivateAttributeValue(id string) (*policy.Value, error) {
	_, err := h.sdk.Attributes.DeactivateAttributeValue(h.ctx, &attributes.DeactivateAttributeValueRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return h.GetAttributeValue(id)
}
