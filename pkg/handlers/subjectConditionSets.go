package handlers

import (
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/subjectmapping"
)

func (h Handler) GetSubjectConditionSet(id string) (*policy.SubjectConditionSet, error) {
	resp, err := h.sdk.SubjectMapping.GetSubjectConditionSet(h.ctx, &subjectmapping.GetSubjectConditionSetRequest{
		Id: id,
	})
	if err != nil {
		return nil, err
	}

	return resp.GetSubjectConditionSet(), nil
}

func (h Handler) ListSubjectConditionSets() ([]*policy.SubjectConditionSet, error) {
	resp, err := h.sdk.SubjectMapping.ListSubjectConditionSets(h.ctx, &subjectmapping.ListSubjectConditionSetsRequest{})
	if err != nil {
		return nil, err
	}
	return resp.GetSubjectConditionSets(), err
}

// Creates and returns the created subject condition set
func (h Handler) CreateSubjectConditionSet(ss []*policy.SubjectSet, metadata *common.MetadataMutable) (*policy.SubjectConditionSet, error) {
	resp, err := h.sdk.SubjectMapping.CreateSubjectConditionSet(h.ctx, &subjectmapping.CreateSubjectConditionSetRequest{
		SubjectConditionSet: &subjectmapping.SubjectConditionSetCreate{
			SubjectSets: ss,
			Metadata:    metadata,
		},
	})
	if err != nil {
		return nil, err
	}
	return h.GetSubjectConditionSet(resp.GetSubjectConditionSet().GetId())
}

// Updates and returns the updated subject condition set
func (h Handler) UpdateSubjectConditionSet(id string, ss []*policy.SubjectSet, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.SubjectConditionSet, error) {
	_, err := h.sdk.SubjectMapping.UpdateSubjectConditionSet(h.ctx, &subjectmapping.UpdateSubjectConditionSetRequest{
		Id:                     id,
		SubjectSets:            ss,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	})
	if err != nil {
		return nil, err
	}
	return h.GetSubjectConditionSet(id)
}

func (h Handler) DeleteSubjectConditionSet(id string) error {
	_, err := h.sdk.SubjectMapping.DeleteSubjectConditionSet(h.ctx, &subjectmapping.DeleteSubjectConditionSetRequest{
		Id: id,
	})
	return err
}

func (h Handler) PruneSubjectConditionSets() ([]*policy.SubjectConditionSet, error) {
	rsp, err := h.sdk.SubjectMapping.DeleteAllUnmappedSubjectConditionSets(h.ctx, &subjectmapping.DeleteAllUnmappedSubjectConditionSetsRequest{})
	if err != nil {
		return nil, err
	}
	return rsp.GetSubjectConditionSets(), nil
}
