package handlers

import (
	"fmt"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/attributes"
	"github.com/opentdf/platform/protocol/go/policy/unsafe"
)

// TODO: Might be useful to map out the attribute rule definitions for help text in the CLI and TUI

const (
	AttributeRuleAllOf     = "ALL_OF"
	AttributeRuleAnyOf     = "ANY_OF"
	AttributeRuleHierarchy = "HIERARCHY"
)

type CreateAttributeError struct {
	ValueErrors map[string]error

	Err error
}

func (e *CreateAttributeError) Error() string {
	if e.ValueErrors != nil {
		return "Error creating attribute with values" + fmt.Sprintf("%v", e.ValueErrors)
	}

	return "Error creating attribute"
}

func (h Handler) GetAttribute(id string) (*policy.Attribute, error) {
	resp, err := h.sdk.Attributes.GetAttribute(h.ctx, &attributes.GetAttributeRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Attribute, nil
}

func (h Handler) ListAttributes(state common.ActiveStateEnum) ([]*policy.Attribute, error) {
	resp, err := h.sdk.Attributes.ListAttributes(h.ctx, &attributes.ListAttributesRequest{State: state})
	if err != nil {
		return nil, err
	}
	return resp.Attributes, err
}

// Creates and returns the created attribute
func (h Handler) CreateAttribute(name string, rule string, namespace string, values []string, metadata *common.MetadataMutable) (*policy.Attribute, error) {
	r, err := GetAttributeRuleFromReadableString(rule)
	if err != nil {
		return nil, err
	}

	attrReq := &attributes.CreateAttributeRequest{
		NamespaceId: namespace,
		Name:        name,
		Rule:        r,
		Metadata:    metadata,
		Values:      values,
	}

	resp, err := h.sdk.Attributes.CreateAttribute(h.ctx, attrReq)
	if err != nil {
		return nil, err
	}

	return h.GetAttribute(resp.GetAttribute().GetId())
}

// Updates and returns updated attribute
func (h *Handler) UpdateAttribute(
	id string,
	metadata *common.MetadataMutable,
	behavior common.MetadataUpdateEnum,
) (*policy.Attribute, error) {
	_, err := h.sdk.Attributes.UpdateAttribute(h.ctx, &attributes.UpdateAttributeRequest{
		Id:                     id,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	})
	if err != nil {
		return nil, err
	}
	return h.GetAttribute(id)
}

// Deactivates and returns deactivated attribute
func (h Handler) DeactivateAttribute(id string) (*policy.Attribute, error) {
	_, err := h.sdk.Attributes.DeactivateAttribute(h.ctx, &attributes.DeactivateAttributeRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return h.GetAttribute(id)
}

// Reactivates and returns reactivated attribute
func (h Handler) UnsafeReactivateAttribute(id string) (*policy.Attribute, error) {
	_, err := h.sdk.Unsafe.UnsafeReactivateAttribute(h.ctx, &unsafe.UnsafeReactivateAttributeRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return h.GetAttribute(id)
}

// Deletes and returns error if deletion failed
func (h Handler) UnsafeDeleteAttribute(id, fqn string) error {
	_, err := h.sdk.Unsafe.UnsafeDeleteAttribute(h.ctx, &unsafe.UnsafeDeleteAttributeRequest{
		Id:  id,
		Fqn: fqn,
	})
	return err
}

// Deletes and returns error if deletion failed
func (h Handler) UnsafeUpdateAttribute(id, name, rule string, values_order []string) error {
	req := &unsafe.UnsafeUpdateAttributeRequest{
		Id:   id,
		Name: name,
	}

	if rule != "" {
		r, err := GetAttributeRuleFromReadableString(rule)
		if err != nil {
			return fmt.Errorf("invalid attribute rule: %s", rule)
		}
		req.Rule = r
	}
	if len(values_order) > 0 {
		req.ValuesOrder = values_order
	}

	_, err := h.sdk.Unsafe.UnsafeUpdateAttribute(h.ctx, req)
	return err
}

func GetAttributeFqn(namespace string, name string) string {
	return fmt.Sprintf("https://%s/attr/%s", namespace, name)
}

func GetAttributeRuleOptions() []string {
	return []string{
		AttributeRuleAllOf,
		AttributeRuleAnyOf,
		AttributeRuleHierarchy,
	}
}

// Provides the un-prefixed human-readable attribute rule
func GetAttributeRuleFromAttributeType(rule policy.AttributeRuleTypeEnum) string {
	switch rule {
	case policy.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_ALL_OF:
		return AttributeRuleAllOf
	case policy.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_ANY_OF:
		return AttributeRuleAnyOf
	case policy.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_HIERARCHY:
		return AttributeRuleHierarchy
	default:
		return ""
	}
}

func GetAttributeRuleFromReadableString(rule string) (policy.AttributeRuleTypeEnum, error) {
	switch rule {
	case AttributeRuleAllOf:
		return policy.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_ALL_OF, nil
	case AttributeRuleAnyOf:
		return policy.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_ANY_OF, nil
	case AttributeRuleHierarchy:
		return policy.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_HIERARCHY, nil
	}
	return 0, fmt.Errorf("invalid attribute rule: %s, must be one of [%s, %s, %s]", rule, AttributeRuleAllOf, AttributeRuleAnyOf, AttributeRuleHierarchy)
}
