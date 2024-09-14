package handlers

import (
	"connectrpc.com/connect"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/attributes"
	"github.com/opentdf/platform/protocol/go/policy/unsafe"
)

func (h *Handler) ListAttributeValues(attributeId string, state common.ActiveStateEnum) ([]*policy.Value, error) {
	resp, err := h.sdk.Attributes.ListAttributeValues(h.ctx, &connect.Request[attributes.ListAttributeValuesRequest]{
		Msg: &attributes.ListAttributeValuesRequest{AttributeId: attributeId, State: state}})
	if err != nil {
		return nil, err
	}
	return resp.Msg.GetValues(), err
}

// Creates and returns the created value
func (h *Handler) CreateAttributeValue(attributeId string, value string, metadata *common.MetadataMutable) (*policy.Value, error) {
	resp, err := h.sdk.Attributes.CreateAttributeValue(h.ctx, &connect.Request[attributes.CreateAttributeValueRequest]{
		Msg: &attributes.CreateAttributeValueRequest{
			AttributeId: attributeId,
			Value:       value,
			Metadata:    metadata,
		}})
	if err != nil {
		return nil, err
	}

	return h.GetAttributeValue(resp.Msg.GetValue().GetId())
}

func (h *Handler) GetAttributeValue(id string) (*policy.Value, error) {
	resp, err := h.sdk.Attributes.GetAttributeValue(h.ctx, &connect.Request[attributes.GetAttributeValueRequest]{
		Msg: &attributes.GetAttributeValueRequest{
			Id: id,
		}})
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetValue(), nil
}

// Updates and returns updated value
func (h *Handler) UpdateAttributeValue(id string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.Value, error) {
	resp, err := h.sdk.Attributes.UpdateAttributeValue(h.ctx, &connect.Request[attributes.UpdateAttributeValueRequest]{
		Msg: &attributes.UpdateAttributeValueRequest{
			Id:                     id,
			Metadata:               metadata,
			MetadataUpdateBehavior: behavior,
		}})
	if err != nil {
		return nil, err
	}

	return h.GetAttributeValue(resp.Msg.GetValue().GetId())
}

// Deactivates and returns deactivated value
func (h *Handler) DeactivateAttributeValue(id string) (*policy.Value, error) {
	_, err := h.sdk.Attributes.DeactivateAttributeValue(h.ctx, &connect.Request[attributes.DeactivateAttributeValueRequest]{
		Msg: &attributes.DeactivateAttributeValueRequest{
			Id: id,
		}})
	if err != nil {
		return nil, err
	}
	return h.GetAttributeValue(id)
}

// Reactivates and returns reactivated attribute
func (h Handler) UnsafeReactivateAttributeValue(id string) (*policy.Value, error) {
	_, err := h.sdk.Unsafe.UnsafeReactivateAttributeValue(h.ctx, &connect.Request[unsafe.UnsafeReactivateAttributeValueRequest]{
		Msg: &unsafe.UnsafeReactivateAttributeValueRequest{
			Id: id,
		}})
	if err != nil {
		return nil, err
	}
	return h.GetAttributeValue(id)
}

// Deletes and returns error if deletion failed
func (h Handler) UnsafeDeleteAttributeValue(id, fqn string) error {
	_, err := h.sdk.Unsafe.UnsafeDeleteAttributeValue(h.ctx, &connect.Request[unsafe.UnsafeDeleteAttributeValueRequest]{
		Msg: &unsafe.UnsafeDeleteAttributeValueRequest{
			Id:  id,
			Fqn: fqn,
		}})
	return err
}

// Deletes and returns error if deletion failed
func (h Handler) UnsafeUpdateAttributeValue(id, value string) error {
	req := &unsafe.UnsafeUpdateAttributeValueRequest{
		Id:    id,
		Value: value,
	}

	_, err := h.sdk.Unsafe.UnsafeUpdateAttributeValue(h.ctx, &connect.Request[unsafe.UnsafeUpdateAttributeValueRequest]{Msg: req})
	return err
}
