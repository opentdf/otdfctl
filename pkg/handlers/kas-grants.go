package handlers

import (
	"github.com/opentdf/platform/protocol/go/policy/attributes"
)

func (h Handler) UpdateKasGrantForAttribute(attr_id string, kas_id string) (*attributes.AttributeKeyAccessServer, error) {
	kas := &attributes.AttributeKeyAccessServer{
		AttributeId:       attr_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Attributes.AssignKeyAccessServerToAttribute(h.ctx, &attributes.AssignKeyAccessServerToAttributeRequest{
		AttributeKeyAccessServer: kas,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetAttributeKeyAccessServer(), nil
}

func (h Handler) DeleteKasGrantFromAttribute(attr_id string, kas_id string) (*attributes.AttributeKeyAccessServer, error) {
	kas := &attributes.AttributeKeyAccessServer{
		AttributeId:       attr_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Attributes.RemoveKeyAccessServerFromAttribute(h.ctx, &attributes.RemoveKeyAccessServerFromAttributeRequest{
		AttributeKeyAccessServer: kas,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetAttributeKeyAccessServer(), nil
}

func (h Handler) UpdateKasGrantForValue(val_id string, kas_id string) (*attributes.ValueKeyAccessServer, error) {
	kas := &attributes.ValueKeyAccessServer{
		ValueId:           val_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Attributes.AssignKeyAccessServerToValue(h.ctx, &attributes.AssignKeyAccessServerToValueRequest{
		ValueKeyAccessServer: kas,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetValueKeyAccessServer(), nil
}

func (h Handler) DeleteKasGrantFromValue(val_id string, kas_id string) (*attributes.ValueKeyAccessServer, error) {
	kas := &attributes.ValueKeyAccessServer{
		ValueId:           val_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Attributes.RemoveKeyAccessServerFromValue(h.ctx, &attributes.RemoveKeyAccessServerFromValueRequest{
		ValueKeyAccessServer: kas,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetValueKeyAccessServer(), nil
}
