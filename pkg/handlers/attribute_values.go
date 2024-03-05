package handlers

import (
	"github.com/opentdf/platform/protocol/go/policy/attributes"
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
