package handlers

import (
	"connectrpc.com/connect"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/subjectmapping"
)

const (
	SubjectMappingOperatorIn          = "IN"
	SubjectMappingOperatorNotIn       = "NOT_IN"
	SubjectMappingOperatorInContains  = "IN_CONTAINS"
	SubjectMappingOperatorUnspecified = "UNSPECIFIED"
)

var SubjectMappingOperatorEnumChoices = []string{SubjectMappingOperatorIn, SubjectMappingOperatorNotIn, SubjectMappingOperatorUnspecified}

func (h Handler) GetSubjectMapping(id string) (*policy.SubjectMapping, error) {
	resp, err := h.sdk.SubjectMapping.GetSubjectMapping(h.ctx, &connect.Request[subjectmapping.GetSubjectMappingRequest]{
		Msg: &subjectmapping.GetSubjectMappingRequest{
			Id: id,
		}})
	return resp.Msg.GetSubjectMapping(), err
}

func (h Handler) ListSubjectMappings() ([]*policy.SubjectMapping, error) {
	resp, err := h.sdk.SubjectMapping.ListSubjectMappings(h.ctx, &connect.Request[subjectmapping.ListSubjectMappingsRequest]{
		Msg: &subjectmapping.ListSubjectMappingsRequest{}})

	return resp.Msg.GetSubjectMappings(), err
}

// Creates and returns the created subject mapping
func (h Handler) CreateNewSubjectMapping(attrValId string, actions []*policy.Action, existingSCSId string, newScs *subjectmapping.SubjectConditionSetCreate, m *common.MetadataMutable) (*policy.SubjectMapping, error) {
	resp, err := h.sdk.SubjectMapping.CreateSubjectMapping(h.ctx, &connect.Request[subjectmapping.CreateSubjectMappingRequest]{
		Msg: &subjectmapping.CreateSubjectMappingRequest{
			AttributeValueId:              attrValId,
			Actions:                       actions,
			ExistingSubjectConditionSetId: existingSCSId,
			NewSubjectConditionSet:        newScs,
			Metadata:                      m,
		}})
	if err != nil {
		return nil, err
	}
	return h.GetSubjectMapping(resp.Msg.GetSubjectMapping().GetId())
}

// Updates and returns the updated subject mapping
func (h Handler) UpdateSubjectMapping(id string, updatedSCSId string, updatedActions []*policy.Action, metadata *common.MetadataMutable, metadataBehavior common.MetadataUpdateEnum) (*policy.SubjectMapping, error) {
	_, err := h.sdk.SubjectMapping.UpdateSubjectMapping(h.ctx, &connect.Request[subjectmapping.UpdateSubjectMappingRequest]{
		Msg: &subjectmapping.UpdateSubjectMappingRequest{
			Id:                     id,
			SubjectConditionSetId:  updatedSCSId,
			Actions:                updatedActions,
			MetadataUpdateBehavior: metadataBehavior,
			Metadata:               metadata,
		}})
	if err != nil {
		return nil, err
	}
	return h.GetSubjectMapping(id)
}

func (h Handler) DeleteSubjectMapping(id string) (*policy.SubjectMapping, error) {
	resp, err := h.sdk.SubjectMapping.DeleteSubjectMapping(h.ctx, &connect.Request[subjectmapping.DeleteSubjectMappingRequest]{
		Msg: &subjectmapping.DeleteSubjectMappingRequest{
			Id: id,
		}})
	return resp.Msg.GetSubjectMapping(), err
}

func GetSubjectMappingOperatorFromChoice(readable string) policy.SubjectMappingOperatorEnum {
	switch readable {
	case SubjectMappingOperatorIn:
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_IN
	case SubjectMappingOperatorNotIn:
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_NOT_IN
	case SubjectMappingOperatorInContains:
		return policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_IN_CONTAINS
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
	case policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_IN_CONTAINS:
		return SubjectMappingOperatorInContains
	case policy.SubjectMappingOperatorEnum_SUBJECT_MAPPING_OPERATOR_ENUM_UNSPECIFIED:
		return SubjectMappingOperatorUnspecified
	default:
		return SubjectMappingOperatorUnspecified
	}
}
