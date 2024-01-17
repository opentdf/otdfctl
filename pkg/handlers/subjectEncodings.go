package handlers

import (
	acsev1 "github.com/opentdf/opentdf-v2-poc/gen/acse/v1"
	commonv1 "github.com/opentdf/opentdf-v2-poc/gen/common/v1"
	attributesv1 "github.com/opentdf/opentdf-v2-poc/gen/attributes/v1"

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

func CreateSubjectMapping(mapping SubjectMapping, description string, resourceDeps []string, attrRefName string, attrRefLabels map[string]string) error {
	client := acsev1.NewSubjectEncodingServiceClient(grpc.Conn)

	// Hierarchy: prefer 1st name, 2nd labels in order
	ref := &attributesv1.AttributeValueReference{}
	if attrRefName != "" {
		ref.Ref = &commonv1.ResourceSelector_Name{Name: attrRefName}
	} else if len(attrRefLabels) > 0 {
		ref.Ref = &attributesv1.AttributeValueReference_AttributeValue{
			AttributeValue: &attributesv1.AttributeValue{
				Descriptor_: &commonv1.ResourceDescriptor{
					Name: "resource-selector-labels",
		},
		}
	}

	_, err := client.CreateSubjectMapping(grpc.Context, &acsev1.CreateSubjectMappingRequest{
		SubjectMapping: &acsev1.SubjectMapping{
			Descriptor_:      &commonv1.ResourceDescriptor{Name: mapping.Name},
			SubjectAttribute: mapping.SubjectAttr,
			SubjectValues:    mapping.SubjectValues,
			Operator:         GetSubjectMappingOperatorFromReadableOperatorString(mapping.Operator),
			AttributeValueRef: ref,
		},
	})
	return err
}

func GetSubjectMapping(id int) (SubjectMapping, error) {
	client := acsev1.NewSubjectEncodingServiceClient(grpc.Conn)
	resp, err := client.GetSubjectMapping(grpc.Context, &acsev1.GetSubjectMappingRequest{
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

func ListSubjectMappings(attrRefName string, resourceSelectorLabels map[string]string) ([]SubjectMapping, error) {
	client := acsev1.NewSubjectEncodingServiceClient(grpc.Conn)

	// Hierarchy: prefer 1st name, 2nd labels in order
	s := &commonv1.ResourceSelector{}
	if attrRefName != "" {
		s.Selector = &commonv1.ResourceSelector_Name{Name: attrRefName}
	} else if len(resourceSelectorLabels) > 0 {
		s.Selector = &commonv1.ResourceSelector_LabelSelector_{
			LabelSelector: &commonv1.ResourceSelector_LabelSelector{
				Labels: resourceSelectorLabels,
			},
		}
	}

	resp, err := client.ListSubjectMappings(grpc.Context, &acsev1.ListSubjectMappingsRequest{
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

func DeleteSubjectMapping(id int) error {
	client := acsev1.NewSubjectEncodingServiceClient(grpc.Conn)
	_, err := client.DeleteSubjectMapping(grpc.Context, &acsev1.DeleteSubjectMappingRequest{
		Id: int32(id),
	})
	return err
}

func GetSubjectMappingOperatorFromIota(operator acsev1.SubjectMapping_Operator) string {
	switch operator {
	case acsev1.SubjectMapping_OPERATOR_IN:
		return SubjectMappingOperator_IN
	case acsev1.SubjectMapping_OPERATOR_NOT_IN:
		return SubjectMappingOperator_NOT_IN
	default:
		return SubjectMappingOperator_UNSPECIFIED
	}
}

func GetSubjectMappingOperatorFromReadableOperatorString(operator string) acsev1.SubjectMapping_Operator {
	switch operator {
	case SubjectMappingOperator_IN:
		return acsev1.SubjectMapping_OPERATOR_IN
	case SubjectMappingOperator_NOT_IN:
		return acsev1.SubjectMapping_OPERATOR_NOT_IN
	default:
		return acsev1.SubjectMapping_OPERATOR_UNSPECIFIED
	}
}

func GetSubjectMappingOperatorOptions() []string {
	return []string{
		SubjectMappingOperator_IN,
		SubjectMappingOperator_NOT_IN,
		SubjectMappingOperator_UNSPECIFIED,
	}
}