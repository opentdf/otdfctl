package handlers

import (
	"context"

	"connectrpc.com/connect"
	"github.com/opentdf/platform/protocol/go/policy/attributes"
	"github.com/opentdf/platform/protocol/go/policy/kasregistry"
	"github.com/opentdf/platform/protocol/go/policy/namespaces"
)

func (h Handler) AssignKasGrantToAttribute(ctx context.Context, attr_id string, kas_id string) (*attributes.AttributeKeyAccessServer, error) {
	kas := &attributes.AttributeKeyAccessServer{
		AttributeId:       attr_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Attributes.AssignKeyAccessServerToAttribute(ctx, &connect.Request[attributes.AssignKeyAccessServerToAttributeRequest]{
		Msg: &attributes.AssignKeyAccessServerToAttributeRequest{
			AttributeKeyAccessServer: kas,
		}})
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetAttributeKeyAccessServer(), nil
}

func (h Handler) DeleteKasGrantFromAttribute(ctx context.Context, attr_id string, kas_id string) (*attributes.AttributeKeyAccessServer, error) {
	kas := &attributes.AttributeKeyAccessServer{
		AttributeId:       attr_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Attributes.RemoveKeyAccessServerFromAttribute(ctx, &connect.Request[attributes.RemoveKeyAccessServerFromAttributeRequest]{
		Msg: &attributes.RemoveKeyAccessServerFromAttributeRequest{
			AttributeKeyAccessServer: kas,
		}})
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetAttributeKeyAccessServer(), nil
}

func (h Handler) AssignKasGrantToValue(ctx context.Context, val_id string, kas_id string) (*attributes.ValueKeyAccessServer, error) {
	kas := &attributes.ValueKeyAccessServer{
		ValueId:           val_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Attributes.AssignKeyAccessServerToValue(ctx, &connect.Request[attributes.AssignKeyAccessServerToValueRequest]{
		Msg: &attributes.AssignKeyAccessServerToValueRequest{
			ValueKeyAccessServer: kas,
		}})
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetValueKeyAccessServer(), nil
}

func (h Handler) DeleteKasGrantFromValue(ctx context.Context, val_id string, kas_id string) (*attributes.ValueKeyAccessServer, error) {
	kas := &attributes.ValueKeyAccessServer{
		ValueId:           val_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Attributes.RemoveKeyAccessServerFromValue(ctx, &connect.Request[attributes.RemoveKeyAccessServerFromValueRequest]{Msg: &attributes.RemoveKeyAccessServerFromValueRequest{
		ValueKeyAccessServer: kas,
	}})
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetValueKeyAccessServer(), nil
}

func (h Handler) AssignKasGrantToNamespace(ctx context.Context, ns_id string, kas_id string) (*namespaces.NamespaceKeyAccessServer, error) {
	kas := &namespaces.NamespaceKeyAccessServer{
		NamespaceId:       ns_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Namespaces.AssignKeyAccessServerToNamespace(ctx, &connect.Request[namespaces.AssignKeyAccessServerToNamespaceRequest]{
		Msg: &namespaces.AssignKeyAccessServerToNamespaceRequest{
			NamespaceKeyAccessServer: kas,
		}})
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetNamespaceKeyAccessServer(), nil
}

func (h Handler) DeleteKasGrantFromNamespace(ctx context.Context, ns_id string, kas_id string) (*namespaces.NamespaceKeyAccessServer, error) {
	kas := &namespaces.NamespaceKeyAccessServer{
		NamespaceId:       ns_id,
		KeyAccessServerId: kas_id,
	}
	resp, err := h.sdk.Namespaces.RemoveKeyAccessServerFromNamespace(ctx, &connect.Request[namespaces.RemoveKeyAccessServerFromNamespaceRequest]{
		Msg: &namespaces.RemoveKeyAccessServerFromNamespaceRequest{
			NamespaceKeyAccessServer: kas,
		}})
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetNamespaceKeyAccessServer(), nil
}

func (h Handler) ListKasGrants(ctx context.Context, kas_id, kas_uri string) ([]*kasregistry.KeyAccessServerGrants, error) {
	resp, err := h.sdk.KeyAccessServerRegistry.ListKeyAccessServerGrants(ctx, &connect.Request[kasregistry.ListKeyAccessServerGrantsRequest]{
		Msg: &kasregistry.ListKeyAccessServerGrantsRequest{
			KasId:  kas_id,
			KasUri: kas_uri,
		}})
	if err != nil {
		return nil, err
	}
	return resp.Msg.GetGrants(), nil
}
