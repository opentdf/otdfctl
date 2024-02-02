package handlers

import (
	"fmt"
	"slices"
	"strings"

	"github.com/opentdf/opentdf-v2-poc/sdk/common"
	"github.com/opentdf/opentdf-v2-poc/sdk/subjectmapping"
)

const (
	SubjectMappingOperatorIn          = "IN"
	SubjectMappingOperatorNotIn       = "NOT_IN"
	SubjectMappingOperatorUnspecified = "UNSPECIFIED"
)

var SubjectMappingOperatorEnumChoices = []string{SubjectMappingOperatorIn, SubjectMappingOperatorNotIn, SubjectMappingOperatorUnspecified}

func (h Handler) GetSubjectMapping(id string) (*subjectmapping.SubjectMapping, error) {
	resp, err := h.sdk.SubjectMapping.GetSubjectMapping(h.ctx, &subjectmapping.GetSubjectMappingRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.SubjectMapping, nil
}

func (h Handler) ListSubjectMappings() ([]*subjectmapping.SubjectMapping, error) {
	resp, err := h.sdk.SubjectMapping.ListSubjectMappings(h.ctx, &subjectmapping.ListSubjectMappingsRequest{})
	if err != nil {
		return nil, err
	}

	return resp.SubjectMappings, nil
}

func (h Handler) CreateNewSubjectMapping(attributeValueId string, subjectAttribute string, subjectValues []string, operator string, metadata *common.MetadataMutable) (*subjectmapping.SubjectMapping, error) {
	if !slices.Contains(SubjectMappingOperatorEnumChoices, operator) {
		return nil, fmt.Errorf("Invalid operator. Must be one of [%s]" + strings.Join(SubjectMappingOperatorEnumChoices, ", "))
	}

	resp, err := h.sdk.SubjectMapping.CreateSubjectMapping(h.ctx, &subjectmapping.CreateSubjectMappingRequest{
		SubjectMapping: &subjectmapping.SubjectMappingCreateUpdate{
			AttributeValueId: attributeValueId,
			SubjectAttribute: subjectAttribute,
			SubjectValues:    subjectValues,
			Operator:         GetSubjectMappingOperatorFromChoice(operator),
			Metadata:         metadata,
		},
	})
	if err != nil {
		return nil, err
	}

	return resp.SubjectMapping, nil
}

func (h Handler) UpdateSubjectMapping(id string, attributeValueId string, subjectAttribute string, subjectValues []string, operator string, metadata *common.MetadataMutable) (*subjectmapping.SubjectMapping, error) {
	if !slices.Contains(SubjectMappingOperatorEnumChoices, operator) {
		return nil, fmt.Errorf("Invalid operator. Must be one of [%s]" + strings.Join(SubjectMappingOperatorEnumChoices, ", "))
	}

	resp, err := h.sdk.SubjectMapping.UpdateSubjectMapping(h.ctx, &subjectmapping.UpdateSubjectMappingRequest{
		Id: id,
		SubjectMapping: &subjectmapping.SubjectMappingCreateUpdate{
			AttributeValueId: attributeValueId,
			SubjectAttribute: subjectAttribute,
			SubjectValues:    subjectValues,
			Operator:         GetSubjectMappingOperatorFromChoice(operator),
			Metadata:         metadata,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.SubjectMapping, nil
}

func (h Handler) DeleteSubjectMapping(id string) error {
	_, err := h.sdk.SubjectMapping.DeleteSubjectMapping(h.ctx, &subjectmapping.DeleteSubjectMappingRequest{
		Id: id,
	})
	if err != nil {
		return err
	}

	return nil
}

func GetSubjectMappingOperatorFromChoice(readable string) subjectmapping.SubjectMappingOperatorEnum {
	switch readable {
	case SubjectMappingOperatorIn:
		return subjectmapping.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_IN
	case SubjectMappingOperatorNotIn:
		return subjectmapping.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_NOT_IN
	case SubjectMappingOperatorUnspecified:
		return subjectmapping.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_UNSPECIFIED
	default:
		return subjectmapping.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_UNSPECIFIED
	}
}

func GetSubjectMappingOperatorChoiceFromEnum(enum subjectmapping.SubjectMappingOperatorEnum) string {
	switch enum {
	case subjectmapping.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_IN:
		return SubjectMappingOperatorIn
	case subjectmapping.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_NOT_IN:
		return SubjectMappingOperatorNotIn
	case subjectmapping.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_UNSPECIFIED:
		return SubjectMappingOperatorUnspecified
	default:
		return SubjectMappingOperatorUnspecified
	}
}
