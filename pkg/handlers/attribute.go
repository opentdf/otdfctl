package handlers

import (
	"fmt"

	attributesv1 "github.com/opentdf/opentdf-v2-poc/gen/attributes/v1"
	commonv1 "github.com/opentdf/opentdf-v2-poc/gen/common/v1"
	"github.com/opentdf/tructl/pkg/grpc"
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

func GetAttribute(id int) (Attribute, error) {
	client := attributesv1.NewAttributesServiceClient(grpc.Conn)
	resp, err := client.GetAttribute(grpc.Context, &attributesv1.GetAttributeRequest{
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

func ListAttributes() ([]Attribute, error) {
	client := attributesv1.NewAttributesServiceClient(grpc.Conn)
	resp, err := client.ListAttributes(grpc.Context, &attributesv1.ListAttributesRequest{})
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

func CreateAttribute(name string, rule string, values []string, namespace string, description string) (Attribute, error) {
	var attrValues []*attributesv1.AttributeDefinitionValue
	for _, v := range values {
		if v != "" {
			attrValues = append(attrValues, &attributesv1.AttributeDefinitionValue{Value: v})
		}
	}

	client := attributesv1.NewAttributesServiceClient(grpc.Conn)
	_, err := client.CreateAttribute(grpc.Context, &attributesv1.CreateAttributeRequest{
		Definition: &attributesv1.AttributeDefinition{
			Name:   name,
			Rule:   GetAttributeRuleFromReadableString(rule),
			Values: attrValues,
			Descriptor_: &commonv1.ResourceDescriptor{
				Namespace:   namespace,
				Name:        name,
				Type:        commonv1.PolicyResourceType_POLICY_RESOURCE_TYPE_ATTRIBUTE_DEFINITION,
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

func DeleteAttribute(id int) error {
	client := attributesv1.NewAttributesServiceClient(grpc.Conn)

	_, err := client.DeleteAttribute(grpc.Context, &attributesv1.DeleteAttributeRequest{
		Id: int32(id),
	})

	return err
}

func GetAttributeFqn(resp *attributesv1.AttributeDefinition) string {
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

func GetAttributeRuleFromAttributeType(rule attributesv1.AttributeDefinition_AttributeRuleType) string {
	switch rule {
	case attributesv1.AttributeDefinition_ATTRIBUTE_RULE_TYPE_ALL_OF:
		return AttributeRuleAllOf
	case attributesv1.AttributeDefinition_ATTRIBUTE_RULE_TYPE_ANY_OF:
		return AttributeRuleAnyOf
	case attributesv1.AttributeDefinition_ATTRIBUTE_RULE_TYPE_HIERARCHICAL:
		return AttributeRuleHierarchy
	case attributesv1.AttributeDefinition_ATTRIBUTE_RULE_TYPE_UNSPECIFIED:
		return AttributeRuleUnspecified
	default:
		return ""
	}
}

func GetAttributeRuleFromReadableString(rule string) attributesv1.AttributeDefinition_AttributeRuleType {
	switch rule {
	case AttributeRuleAllOf:
		return attributesv1.AttributeDefinition_ATTRIBUTE_RULE_TYPE_ALL_OF
	case AttributeRuleAnyOf:
		return attributesv1.AttributeDefinition_ATTRIBUTE_RULE_TYPE_ANY_OF
	case AttributeRuleHierarchy:
		return attributesv1.AttributeDefinition_ATTRIBUTE_RULE_TYPE_HIERARCHICAL
	case AttributeRuleUnspecified:
		return attributesv1.AttributeDefinition_ATTRIBUTE_RULE_TYPE_UNSPECIFIED
	}
	return attributesv1.AttributeDefinition_ATTRIBUTE_RULE_TYPE_UNSPECIFIED
}
