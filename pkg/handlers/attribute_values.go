package handlers

import (
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
