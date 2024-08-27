package handlers

import (
	"context"

	"github.com/opentdf/platform/protocol/go/policy/attributes"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
)

func (h Handler) AssignKasGrantToAttribute(ctx context.Context, attr_id string, kas_id string) (*attributes.AttributeKeyAccessServer, error) {
	kas := &attributes.AttributeKeyAccessServer{
		AttributeId:       attr_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Attributes.AssignKeyAccessServerToAttribute(ctx, &attributes.AssignKeyAccessServerToAttributeRequest{
		AttributeKeyAccessServer: kas,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetAttributeKeyAccessServer(), nil
}

func (h Handler) DeleteKasGrantFromAttribute(ctx context.Context, attr_id string, kas_id string) (*attributes.AttributeKeyAccessServer, error) {
	kas := &attributes.AttributeKeyAccessServer{
		AttributeId:       attr_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Attributes.RemoveKeyAccessServerFromAttribute(ctx, &attributes.RemoveKeyAccessServerFromAttributeRequest{
		AttributeKeyAccessServer: kas,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetAttributeKeyAccessServer(), nil
}

func (h Handler) AssignKasGrantToValue(ctx context.Context, val_id string, kas_id string) (*attributes.ValueKeyAccessServer, error) {
	kas := &attributes.ValueKeyAccessServer{
		ValueId:           val_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Attributes.AssignKeyAccessServerToValue(ctx, &attributes.AssignKeyAccessServerToValueRequest{
		ValueKeyAccessServer: kas,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetValueKeyAccessServer(), nil
}

func (h Handler) DeleteKasGrantFromValue(ctx context.Context, val_id string, kas_id string) (*attributes.ValueKeyAccessServer, error) {
	kas := &attributes.ValueKeyAccessServer{
		ValueId:           val_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Attributes.RemoveKeyAccessServerFromValue(ctx, &attributes.RemoveKeyAccessServerFromValueRequest{
		ValueKeyAccessServer: kas,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetValueKeyAccessServer(), nil
}

func (h Handler) AssignKasGrantToNamespace(ctx context.Context, ns_id string, kas_id string) (*namespaces.NamespaceKeyAccessServer, error) {
	kas := &namespaces.NamespaceKeyAccessServer{
		NamespaceId:       ns_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Namespaces.AssignKeyAccessServerToNamespace(ctx, &namespaces.AssignKeyAccessServerToNamespaceRequest{
		NamespaceKeyAccessServer: kas,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetNamespaceKeyAccessServer(), nil
}

func (h Handler) DeleteKasGrantFromNamespace(ctx context.Context, ns_id string, kas_id string) (*namespaces.NamespaceKeyAccessServer, error) {
	kas := &namespaces.NamespaceKeyAccessServer{
		NamespaceId:       ns_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Namespaces.RemoveKeyAccessServerFromNamespace(ctx, &namespaces.RemoveKeyAccessServerFromNamespaceRequest{
		NamespaceKeyAccessServer: kas,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetNamespaceKeyAccessServer(), nil
}
