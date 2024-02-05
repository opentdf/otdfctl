package handlers

import (
	"fmt"

	"github.com/opentdf/opentdf-v2-poc/sdk/attributes"
	"github.com/opentdf/opentdf-v2-poc/sdk/common"
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

func (h Handler) GetAttribute(id string) (*attributes.Attribute, error) {
	resp, err := h.sdk.Attributes.GetAttribute(h.ctx, &attributes.GetAttributeRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.Attribute, nil
}

func (h Handler) ListAttributes() ([]*attributes.Attribute, error) {
	resp, err := h.sdk.Attributes.ListAttributes(h.ctx, &attributes.ListAttributesRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Attributes, err
}

func (h Handler) CreateAttribute(name string, rule string, values []string, namespace string) (*attributes.Attribute, error) {
	r, err := GetAttributeRuleFromReadableString(rule)
	if err != nil {
		return nil, err
	}

	attrReq := &attributes.CreateAttributeRequest{
		Attribute: &attributes.AttributeCreateUpdate{
			NamespaceId: namespace,
			Name:        name,
			Rule:        r,
		},
	}

	resp, err := h.sdk.Attributes.CreateAttribute(h.ctx, attrReq)
	if err != nil {
		return nil, err
	}

	// short-circuit if there are no values
	if len(values) == 0 {
		return resp.Attribute, nil
	}

	// create attribute values
	attr := resp.Attribute
	attrValues := make([]*attributes.Value, 0, len(values))
	valueErrors := make(map[string]error)
	for _, value := range values {
		v, err := h.CreateAttributeValue(attr.Id, value)
		if err != nil {
			valueErrors[value] = err
		}
		attrValues = append(attrValues, v)
	}

	// return the attribute and any errors
	return &attributes.Attribute{
		Id:        attr.Id,
		Name:      attr.Name,
		Rule:      attr.Rule,
		Values:    attrValues,
		Namespace: attr.Namespace,
	}, &CreateAttributeError{ValueErrors: valueErrors}
}

func (h *Handler) UpdateAttribute(
	id string,
	fns ...func(*common.MetadataMutable) *common.MetadataMutable,
) (*attributes.UpdateAttributeResponse, error) {
	return h.sdk.Attributes.UpdateAttribute(h.ctx, &attributes.UpdateAttributeRequest{
		Id: id,
		Attribute: &attributes.AttributeCreateUpdate{
			Metadata: buildMetadata(&common.MetadataMutable{}, fns...),
		},
	})
}

func (h Handler) DeleteAttribute(id string) (*attributes.Attribute, error) {
	resp, err := h.sdk.Attributes.DeleteAttribute(h.ctx, &attributes.DeleteAttributeRequest{
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

func GetAttributeRuleFromAttributeType(rule attributes.AttributeRuleTypeEnum) string {
	switch rule {
	case attributes.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_ALL_OF:
		return AttributeRuleAllOf
	case attributes.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_ANY_OF:
		return AttributeRuleAnyOf
	case attributes.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_HIERARCHY:
		return AttributeRuleHierarchy
	default:
		return ""
	}
}

func GetAttributeRuleFromReadableString(rule string) (attributes.AttributeRuleTypeEnum, error) {
	switch rule {
	case AttributeRuleAllOf:
		return attributes.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_ALL_OF, nil
	case AttributeRuleAnyOf:
		return attributes.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_ANY_OF, nil
	case AttributeRuleHierarchy:
		return attributes.AttributeRuleTypeEnum_ATTRIBUTE_RULE_TYPE_ENUM_HIERARCHY, nil
	}
	return 0, fmt.Errorf("invalid attribute rule: %s", rule)
}
