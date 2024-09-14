package handlers

import (
	"connectrpc.com/connect"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/subjectmapping"
)

func (h Handler) GetSubjectConditionSet(id string) (*policy.SubjectConditionSet, error) {
	resp, err := h.sdk.SubjectMapping.GetSubjectConditionSet(h.ctx, &connect.Request[subjectmapping.GetSubjectConditionSetRequest]{
		Msg: &subjectmapping.GetSubjectConditionSetRequest{
			Id: id,
		}})
	if err != nil {
		return nil, err
	}

	return resp.Msg.GetSubjectConditionSet(), nil
}

func (h Handler) ListSubjectConditionSets() ([]*policy.SubjectConditionSet, error) {
	resp, err := h.sdk.SubjectMapping.ListSubjectConditionSets(h.ctx, &connect.Request[subjectmapping.ListSubjectConditionSetsRequest]{
		Msg: &subjectmapping.ListSubjectConditionSetsRequest{}})
	if err != nil {
		return nil, err
	}
	return resp.Msg.GetSubjectConditionSets(), err
}

// Creates and returns the created subject condition set
func (h Handler) CreateSubjectConditionSet(ss []*policy.SubjectSet, metadata *common.MetadataMutable) (*policy.SubjectConditionSet, error) {
	resp, err := h.sdk.SubjectMapping.CreateSubjectConditionSet(h.ctx, &connect.Request[subjectmapping.CreateSubjectConditionSetRequest]{
		Msg: &subjectmapping.CreateSubjectConditionSetRequest{
			SubjectConditionSet: &subjectmapping.SubjectConditionSetCreate{
				SubjectSets: ss,
				Metadata:    metadata,
			},
		}})
	if err != nil {
		return nil, err
	}
	return h.GetSubjectConditionSet(resp.Msg.GetSubjectConditionSet().GetId())
}

// Updates and returns the updated subject condition set
func (h Handler) UpdateSubjectConditionSet(id string, ss []*policy.SubjectSet, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.SubjectConditionSet, error) {
	_, err := h.sdk.SubjectMapping.UpdateSubjectConditionSet(h.ctx, &connect.Request[subjectmapping.UpdateSubjectConditionSetRequest]{
		Msg: &subjectmapping.UpdateSubjectConditionSetRequest{
			Id:                     id,
			SubjectSets:            ss,
			Metadata:               metadata,
			MetadataUpdateBehavior: behavior,
		}})
	if err != nil {
		return nil, err
	}
	return h.GetSubjectConditionSet(id)
}

func (h Handler) DeleteSubjectConditionSet(id string) error {
	_, err := h.sdk.SubjectMapping.DeleteSubjectConditionSet(h.ctx, &connect.Request[subjectmapping.DeleteSubjectConditionSetRequest]{
		Msg: &subjectmapping.DeleteSubjectConditionSetRequest{
			Id: id,
		}})
	return err
}
