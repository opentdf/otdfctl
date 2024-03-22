package handlers

import (
	"github.com/opentdf/opentdf-v2-poc/sdk/attributes"
)

func (h *Handler) CreateAttributeValue(attributeId string, value string) (*attributes.Value, error) {
	resp, err := h.sdk.Attributes.CreateAttributeValue(h.ctx, &attributes.CreateAttributeValueRequest{
		AttributeId: attributeId,
		Value: &attributes.ValueCreateUpdate{
			Value: value,
		},
	})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

func (h *Handler) GetAttributeValue(id string) (*attributes.Value, error) {
	resp, err := h.sdk.Attributes.GetAttributeValue(h.ctx, &attributes.GetAttributeValueRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

func (h *Handler) UpdateAttributeValue(id string, value string) (*attributes.Value, error) {
	resp, err := h.sdk.Attributes.UpdateAttributeValue(h.ctx, &attributes.UpdateAttributeValueRequest{
		Id: id,
		Value: &attributes.ValueCreateUpdate{
			Value: value,
		},
	})
	if err != nil {
		return nil, err
	}

	return resp.Value, nil
}

func (h *Handler) DeleteAttributeValue(id string) error {
	_, err := h.sdk.Attributes.DeleteAttributeValue(h.ctx, &attributes.DeleteAttributeValueRequest{
		Id: id,
	})
	return err
}
