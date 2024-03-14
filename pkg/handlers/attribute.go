package handlers

import (
	"fmt"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/attributes"
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

func (h Handler) ListAttributes() ([]*policy.Attribute, error) {
	resp, err := h.sdk.Attributes.ListAttributes(h.ctx, &attributes.ListAttributesRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Attributes, err
}

// TODO: allow creation of a value with the attribute simultaneously?
func (h Handler) CreateAttribute(name string, rule string, namespace string, metadata *common.MetadataMutable) (*policy.Attribute, error) {
	r, err := GetAttributeRuleFromReadableString(rule)
	if err != nil {
		return nil, err
	}

	attrReq := &attributes.CreateAttributeRequest{
		NamespaceId: namespace,
		Name:        name,
		Rule:        r,
		Metadata:    metadata,
	}

	resp, err := h.sdk.Attributes.CreateAttribute(h.ctx, attrReq)
	if err != nil {
		return nil, err
	}

	attr := resp.Attribute

	return &policy.Attribute{
		Id:        attr.Id,
		Name:      attr.Name,
		Rule:      attr.Rule,
		Namespace: attr.Namespace,
	}, nil
}

// TODO: verify updation behavior
func (h *Handler) UpdateAttribute(
	id string,
	metadata *common.MetadataMutable,
	behavior common.MetadataUpdateEnum,
) (*attributes.UpdateAttributeResponse, error) {
	return h.sdk.Attributes.UpdateAttribute(h.ctx, &attributes.UpdateAttributeRequest{
		Id:                     id,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	})
}

func (h Handler) DeactivateAttribute(id string) (*policy.Attribute, error) {
	resp, err := h.sdk.Attributes.DeactivateAttribute(h.ctx, &attributes.DeactivateAttributeRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}
	return resp.Attribute, err
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
	return 0, fmt.Errorf("invalid attribute rule: %s", rule)
}
