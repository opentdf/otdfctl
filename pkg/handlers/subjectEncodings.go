package handlers

import (
	acsev1 "github.com/opentdf/opentdf-v2-poc/gen/acse/v1"
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

func ListSubjectMappings() ([]SubjectMapping, error) {
	client := acsev1.NewSubjectEncodingServiceClient(grpc.Conn)
	resp, err := client.ListSubjectMappings(grpc.Context, &acsev1.ListSubjectMappingsRequest{
		// TODO: selector?
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
