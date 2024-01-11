package handlers

import (
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

func ListAttributes() (*attributesv1.ListAttributesResponse, error) {
	client := attributesv1.NewAttributesServiceClient(grpc.Conn)
	return client.ListAttributes(grpc.Context, &attributesv1.ListAttributesRequest{})
}

func CreateAttribute(name string, rule string, values []string, namespace string, description string) (*attributesv1.CreateAttributeResponse, error) {
	var attrValues []*attributesv1.AttributeDefinitionValue
	for _, v := range values {
		if v != "" {
			attrValues = append(attrValues, &attributesv1.AttributeDefinitionValue{Value: v})
		}
	}

	client := attributesv1.NewAttributesServiceClient(grpc.Conn)
	return client.CreateAttribute(grpc.Context, &attributesv1.CreateAttributeRequest{
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
