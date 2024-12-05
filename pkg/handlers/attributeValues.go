package handlers

import (
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/attributes"
	"github.com/opentdf/platform/protocol/go/policy/unsafe"
)

func (h *Handler) ListAttributeValues(attributeId string, state common.ActiveStateEnum, limit, offset int32) ([]*policy.Value, *policy.PageResponse, error) {
	resp, err := h.sdk.Attributes.ListAttributeValues(h.ctx, &attributes.ListAttributeValuesRequest{
		AttributeId: attributeId,
		State:       state,
		Pagination: &policy.PageRequest{
			Limit:  limit,
			Offset: offset,
		},
	})
	if err != nil {
		return nil, nil, err
	}
	return resp.GetValues(), resp.GetPagination(), err
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
func (h *Handler) UpdateAttributeValue(id string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.Value, error) {
	resp, err := h.sdk.Attributes.UpdateAttributeValue(h.ctx, &attributes.UpdateAttributeValueRequest{
		Id:                     id,
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

// Reactivates and returns reactivated attribute
func (h Handler) UnsafeReactivateAttributeValue(id string) (*policy.Value, error) {
	_, err := h.sdk.Unsafe.UnsafeReactivateAttributeValue(h.ctx, &unsafe.UnsafeReactivateAttributeValueRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return h.GetAttributeValue(id)
}

// Deletes and returns error if deletion failed
func (h Handler) UnsafeDeleteAttributeValue(id, fqn string) error {
	_, err := h.sdk.Unsafe.UnsafeDeleteAttributeValue(h.ctx, &unsafe.UnsafeDeleteAttributeValueRequest{
		Id:  id,
		Fqn: fqn,
	})
	return err
}

// Deletes and returns error if deletion failed
func (h Handler) UnsafeUpdateAttributeValue(id, value string) error {
	req := &unsafe.UnsafeUpdateAttributeValueRequest{
		Id:    id,
		Value: value,
	}

	_, err := h.sdk.Unsafe.UnsafeUpdateAttributeValue(h.ctx, req)
	return err
}
