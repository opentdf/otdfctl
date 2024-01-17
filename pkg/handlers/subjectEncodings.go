package handlers

import (
	acse "github.com/opentdf/opentdf-v2-poc/sdk/acse"
	attributes "github.com/opentdf/opentdf-v2-poc/sdk/attributes"
	common "github.com/opentdf/opentdf-v2-poc/sdk/common"
	"github.com/opentdf/tructl/pkg/grpc"
)

const (
	SubjectMappingOperator_IN          = "IN"
	SubjectMappingOperator_NOT_IN      = "NOT_IN"
	SubjectMappingOperator_UNSPECIFIED = "UNSPECIFIED"
)

type SubjectMapping struct {
	Name          string
	SubjectAttr   string
	SubjectValues []string
	Operator      string // human-readable
}

func (h Handler) CreateSubjectMapping(mapping SubjectMapping, description string, resourceDeps []string, attrRefName string, attrRefLabels map[string]string) error {
	// Hierarchy: prefer 1st name, 2nd labels in order
	ref := &attributes.AttributeValueReference{}
	if attrRefName != "" {
		ref.Ref = &common.ResourceSelector_Name{Name: attrRefName}
	} else if len(attrRefLabels) > 0 {
		ref.Ref = &attributes.AttributeValueReference_AttributeValue{
			AttributeValue: &attributes.AttributeValue{
				Descriptor_: &common.ResourceDescriptor{
					Name: "resource-selector-labels",
				},
			},
		}
	}

	_, err := h.sdk.SubjectEncoding.CreateSubjectMapping(grpc.Context, &acse.CreateSubjectMappingRequest{
		SubjectMapping: &acse.SubjectMapping{
			Descriptor_:       &common.ResourceDescriptor{Name: mapping.Name},
			SubjectAttribute:  mapping.SubjectAttr,
			SubjectValues:     mapping.SubjectValues,
			Operator:          GetSubjectMappingOperatorFromReadableOperatorString(mapping.Operator),
			AttributeValueRef: ref,
		},
	})
	return err
}

func (h Handler) GetSubjectMapping(id int) (SubjectMapping, error) {
	resp, err := h.sdk.SubjectEncoding.GetSubjectMapping(grpc.Context, &acse.GetSubjectMappingRequest{
		Id: int32(id),
	})
	if err != nil {
		return SubjectMapping{}, err
	}

	return SubjectMapping{
		Name:          resp.SubjectMapping.Descriptor_.Name,
		SubjectAttr:   resp.SubjectMapping.SubjectAttribute,
		SubjectValues: resp.SubjectMapping.SubjectValues,
		Operator:      GetSubjectMappingOperatorFromIota(resp.SubjectMapping.Operator),
	}, nil
}

func (h Handler) ListSubjectMappings(attrRefName string, resourceSelectorLabels map[string]string) ([]SubjectMapping, error) {
	// Hierarchy: prefer 1st name, 2nd labels in order
	s := &common.ResourceSelector{}
	if attrRefName != "" {
		s.Selector = &common.ResourceSelector_Name{Name: attrRefName}
	} else if len(resourceSelectorLabels) > 0 {
		s.Selector = &common.ResourceSelector_LabelSelector_{
			LabelSelector: &common.ResourceSelector_LabelSelector{
				Labels: resourceSelectorLabels,
			},
		}
	}

	resp, err := h.sdk.SubjectEncoding.ListSubjectMappings(grpc.Context, &acse.ListSubjectMappingsRequest{
		Selector: s,
	})
	if err != nil {
		return nil, err
	}

	var mappings []SubjectMapping
	for _, m := range resp.SubjectMappings {
		mappings = append(mappings, SubjectMapping{
			Name:          m.Descriptor_.Name,
			SubjectAttr:   m.SubjectAttribute,
			SubjectValues: m.SubjectValues,
			Operator:      GetSubjectMappingOperatorFromIota(m.Operator),
		})
	}

	return mappings, nil
}

func (h Handler) DeleteSubjectMapping(id int) error {
	_, err := h.sdk.SubjectEncoding.DeleteSubjectMapping(grpc.Context, &acse.DeleteSubjectMappingRequest{
		Id: int32(id),
	})
	return err
}

func GetSubjectMappingOperatorFromIota(operator acse.SubjectMapping_Operator) string {
	switch operator {
	case acse.SubjectMapping_OPERATOR_IN:
		return SubjectMappingOperator_IN
	case acse.SubjectMapping_OPERATOR_NOT_IN:
		return SubjectMappingOperator_NOT_IN
	default:
		return SubjectMappingOperator_UNSPECIFIED
	}
}

func GetSubjectMappingOperatorFromReadableOperatorString(operator string) acse.SubjectMapping_Operator {
	switch operator {
	case SubjectMappingOperator_IN:
		return acse.SubjectMapping_OPERATOR_IN
	case SubjectMappingOperator_NOT_IN:
		return acse.SubjectMapping_OPERATOR_NOT_IN
	default:
		return acse.SubjectMapping_OPERATOR_UNSPECIFIED
	}
}

func GetSubjectMappingOperatorOptions() []string {
	return []string{
		SubjectMappingOperator_IN,
		SubjectMappingOperator_NOT_IN,
		SubjectMappingOperator_UNSPECIFIED,
	}
}
