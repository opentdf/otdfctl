package handlers

import (
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/subjectmapping"
)

const (
	SubjectMappingOperatorIn          = "IN"
	SubjectMappingOperatorNotIn       = "NOT_IN"
	SubjectMappingOperatorUnspecified = "UNSPECIFIED"
)

var SubjectMappingOperatorEnumChoices = []string{SubjectMappingOperatorIn, SubjectMappingOperatorNotIn, SubjectMappingOperatorUnspecified}

func (h Handler) GetSubjectMapping(id string) (*policy.SubjectMapping, error) {
	resp, err := h.sdk.SubjectMapping.GetSubjectMapping(h.ctx, &subjectmapping.GetSubjectMappingRequest{
		Id: id,
	})
	return resp.GetSubjectMapping(), err
}

func (h Handler) ListSubjectMappings() ([]*policy.SubjectMapping, error) {
	resp, err := h.sdk.SubjectMapping.ListSubjectMappings(h.ctx, &subjectmapping.ListSubjectMappingsRequest{})

	return resp.GetSubjectMappings(), err
}

// Creates and returns the created subject mapping
func (h Handler) CreateNewSubjectMapping(attrValId string, actions []*policy.Action, existingSCSId string, newScs *subjectmapping.SubjectConditionSetCreate, m *common.MetadataMutable) (*policy.SubjectMapping, error) {
	resp, err := h.sdk.SubjectMapping.CreateSubjectMapping(h.ctx, &subjectmapping.CreateSubjectMappingRequest{
		AttributeValueId:              attrValId,
		Actions:                       actions,
		ExistingSubjectConditionSetId: existingSCSId,
		NewSubjectConditionSet:        newScs,
		Metadata:                      m,
	})
	if err != nil {
		return nil, err
	}
	return h.GetSubjectMapping(resp.GetSubjectMapping().GetId())
}

// Updates and returns the updated subject mapping
func (h Handler) UpdateSubjectMapping(id string, updatedSCSId string, updatedActions []*policy.Action, metadata *common.MetadataMutable, metadataBehavior common.MetadataUpdateEnum) (*policy.SubjectMapping, error) {
	_, err := h.sdk.SubjectMapping.UpdateSubjectMapping(h.ctx, &subjectmapping.UpdateSubjectMappingRequest{
		Id:                     id,
		SubjectConditionSetId:  updatedSCSId,
		Actions:                updatedActions,
		MetadataUpdateBehavior: metadataBehavior,
		Metadata:               metadata,
	})
	if err != nil {
		return nil, err
	}
	return h.GetSubjectMapping(id)
}

func (h Handler) DeleteSubjectMapping(id string) (*policy.SubjectMapping, error) {
	resp, err := h.sdk.SubjectMapping.DeleteSubjectMapping(h.ctx, &subjectmapping.DeleteSubjectMappingRequest{
		Id: id,
	})
	return resp.SubjectMapping, err
}

func GetSubjectMappingOperatorFromChoice(readable string) policy.SubjectMappingOperatorEnum {
	switch readable {
	case SubjectMappingOperatorIn:
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_IN
	case SubjectMappingOperatorNotIn:
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_NOT_IN
	case SubjectMappingOperatorUnspecified:
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_UNSPECIFIED
	default:
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_UNSPECIFIED
	}
}

func GetSubjectMappingOperatorChoiceFromEnum(enum policy.SubjectMappingOperatorEnum) string {
	switch enum {
	case policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_IN:
		return SubjectMappingOperatorIn
	case policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_NOT_IN:
		return SubjectMappingOperatorNotIn
	case policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_UNSPECIFIED:
		return SubjectMappingOperatorUnspecified
	default:
		return SubjectMappingOperatorUnspecified
	}
}
