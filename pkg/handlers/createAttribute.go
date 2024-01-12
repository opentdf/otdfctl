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

func UpdateAttribute(
	Id int32,
	name string,
	rule string,
	values []string,
	groupBy []string,
	resourceId int32,
	resourceVersion int32,
	resourceName string,
	resourceNamespace string,
	resourceFqn string,
	resourceDescription string,
	resourceDependencies []string,
) (*attributesv1.UpdateAttributeResponse, error) {
	var attrValues []*attributesv1.AttributeDefinitionValue
	for _, v := range values {
		if v != "" {
			attrValues = append(attrValues, &attributesv1.AttributeDefinitionValue{Value: v})
		}
	}

	var attrGroupBy []*attributesv1.AttributeDefinitionValue
	for _, v := range groupBy {
		if v != "" {
			attrGroupBy = append(attrGroupBy, &attributesv1.AttributeDefinitionValue{Value: v})
		}
	}

	var dependencies []*commonv1.ResourceDependency
	for _, v := range resourceDependencies {
		if v != "" {
			dependencies = append(dependencies, &commonv1.ResourceDependency{Namespace: v})
		}
	}

	client := attributesv1.NewAttributesServiceClient(grpc.Conn)
	return client.UpdateAttribute(grpc.Context, &attributesv1.UpdateAttributeRequest{
		Id: Id,
		Definition: &attributesv1.AttributeDefinition{
			Name:    name,
			Rule:    GetAttributeRuleFromReadableString(rule),
			Values:  attrValues,
			GroupBy: attrGroupBy,
			Descriptor_: &commonv1.ResourceDescriptor{
				Type:         commonv1.PolicyResourceType_POLICY_RESOURCE_TYPE_ATTRIBUTE_DEFINITION,
				Id:           resourceId,
				Version:      resourceVersion,
				Name:         resourceName,
				Namespace:    resourceNamespace,
				Fqn:          resourceFqn,
				Description:  resourceDescription,
				Dependencies: dependencies,
			},
		},
	})
}

// TODO: do we implement all methods for attribute groups as well, or attributes alone?

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
