package handlers

import (
	"fmt"

	"github.com/opentdf/opentdf-v2-poc/sdk/attributes"
	"github.com/opentdf/opentdf-v2-poc/sdk/common"
)

// TODO: Might be useful to map out the attribute rule definitions for help text in the CLI and TUI

const (
	AttributeRuleAllOf       = "ALL_OF"
	AttributeRuleAnyOf       = "ANY_OF"
	AttributeRuleHierarchy   = "HIERARCHY"
	AttributeRuleUnspecified = "UNSPECIFIED"
)

type Attribute struct {
	Name        string
	Rule        string
	Values      []string
	Namespace   string
	Description string
	Fqn         string
}

func (h Handler) GetAttribute(id int) (Attribute, error) {
	resp, err := h.sdk.Attributes.GetAttribute(h.ctx, &attributes.GetAttributeRequest{
		Id: int32(id),
	})
	if err != nil {
		return Attribute{}, err
	}

	values := []string{}
	for _, v := range resp.Definition.Values {
		values = append(values, v.Value)
	}

	return Attribute{
		Name:        resp.Definition.Name,
		Rule:        GetAttributeRuleFromAttributeType(resp.Definition.Rule),
		Values:      values,
		Namespace:   resp.Definition.Descriptor_.Namespace,
		Description: resp.Definition.Descriptor_.Description,
		Fqn:         GetAttributeFqn(resp.Definition),
	}, nil
}

func (h Handler) ListAttributes() ([]Attribute, error) {
	resp, err := h.sdk.Attributes.ListAttributes(h.ctx, &attributes.ListAttributesRequest{})
	if err != nil {
		return nil, err
	}

	var attrs []Attribute
	for _, attr := range resp.Definitions {
		values := []string{}
		for _, v := range attr.Values {
			values = append(values, v.Value)
		}
		attrs = append(attrs, Attribute{
			Name:        attr.Name,
			Rule:        GetAttributeRuleFromAttributeType(attr.Rule),
			Values:      values,
			Namespace:   attr.Descriptor_.Namespace,
			Description: attr.Descriptor_.Description,
		})
	}

	return attrs, err
}

func (h Handler) CreateAttribute(name string, rule string, values []string, namespace string, description string) (Attribute, error) {
	var attrValues []*attributes.AttributeDefinitionValue
	for _, v := range values {
		if v != "" {
			attrValues = append(attrValues, &attributes.AttributeDefinitionValue{Value: v})
		}
	}

	_, err := h.sdk.Attributes.CreateAttribute(h.ctx, &attributes.CreateAttributeRequest{
		Definition: &attributes.AttributeDefinition{
			Name:   name,
			Rule:   GetAttributeRuleFromReadableString(rule),
			Values: attrValues,
			Descriptor_: &common.ResourceDescriptor{
				Namespace:   namespace,
				Name:        name,
				Type:        common.PolicyResourceType_POLICY_RESOURCE_TYPE_ATTRIBUTE_DEFINITION,
				Description: description,
			},
		},
	})
	if err != nil {
		return Attribute{}, err
	}

	return Attribute{
		Name:        name,
		Rule:        rule,
		Values:      values,
		Namespace:   namespace,
		Description: description,
	}, nil
}

func (h Handler) DeleteAttribute(id int) error {
	_, err := h.sdk.Attributes.DeleteAttribute(h.ctx, &attributes.DeleteAttributeRequest{
		Id: int32(id),
	})

	return err
}

func GetAttributeFqn(resp *attributes.AttributeDefinition) string {
	return fmt.Sprintf("https://%s/attr/%s", resp.Descriptor_.Namespace, resp.Name)
}

func GetAttributeRuleOptions() []string {
	return []string{
		AttributeRuleAllOf,
		AttributeRuleAnyOf,
		AttributeRuleHierarchy,
		AttributeRuleUnspecified,
	}
}

func GetAttributeRuleFromAttributeType(rule attributes.AttributeDefinition_AttributeRuleType) string {
	switch rule {
	case attributes.AttributeDefinition_ATTRIBUTE_RULE_TYPE_ALL_OF:
		return AttributeRuleAllOf
	case attributes.AttributeDefinition_ATTRIBUTE_RULE_TYPE_ANY_OF:
		return AttributeRuleAnyOf
	case attributes.AttributeDefinition_ATTRIBUTE_RULE_TYPE_HIERARCHICAL:
		return AttributeRuleHierarchy
	case attributes.AttributeDefinition_ATTRIBUTE_RULE_TYPE_UNSPECIFIED:
		return AttributeRuleUnspecified
	default:
		return ""
	}
}

func GetAttributeRuleFromReadableString(rule string) attributes.AttributeDefinition_AttributeRuleType {
	switch rule {
	case AttributeRuleAllOf:
		return attributes.AttributeDefinition_ATTRIBUTE_RULE_TYPE_ALL_OF
	case AttributeRuleAnyOf:
		return attributes.AttributeDefinition_ATTRIBUTE_RULE_TYPE_ANY_OF
	case AttributeRuleHierarchy:
		return attributes.AttributeDefinition_ATTRIBUTE_RULE_TYPE_HIERARCHICAL
	case AttributeRuleUnspecified:
		return attributes.AttributeDefinition_ATTRIBUTE_RULE_TYPE_UNSPECIFIED
	}
	return attributes.AttributeDefinition_ATTRIBUTE_RULE_TYPE_UNSPECIFIED
}
