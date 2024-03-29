package handlers

import (
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/attributes"
)

func (h *Handler) CreateAttributeValue(attributeId string, value string) (*policy.Value, error) {
	resp, err := h.sdk.Attributes.CreateAttributeValue(h.ctx, &attributes.CreateAttributeValueRequest{
		AttributeId: attributeId,
		Value:       value,
	})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

func (h *Handler) GetAttributeValue(id string) (*policy.Value, error) {
	resp, err := h.sdk.Attributes.GetAttributeValue(h.ctx, &attributes.GetAttributeValueRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

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

	return resp.Value, nil
}

func (h *Handler) DeleteAttributeValue(id string) error {
	_, err := h.sdk.Attributes.DeactivateAttributeValue(h.ctx, &attributes.DeactivateAttributeValueRequest{
		Id: id,
	})
	return err
}
