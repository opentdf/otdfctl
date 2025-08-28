package handlers

import (
	"context"

	"github.com/google/uuid"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/obligations"
)

//
// Obligations
//

func (h Handler) CreateObligation(ctx context.Context, namespace, name string, values []string, metadata *common.MetadataMutable) (*policy.Obligation, error) {
	req := &obligations.CreateObligationRequest{
		Name:     name,
		Values:   values,
		Metadata: metadata,
	}

	_, err := uuid.Parse(namespace)
	if err != nil {
		req.NamespaceIdentifier = &obligations.CreateObligationRequest_Fqn{
			Fqn: namespace,
		}
	} else {
		req.NamespaceIdentifier = &obligations.CreateObligationRequest_Id{
			Id: namespace,
		}
	}

	resp, err := h.sdk.Obligations.CreateObligation(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetObligation(), nil
}

func (h Handler) GetObligation(ctx context.Context, id, fqn string) (*policy.Obligation, error) {
	req := &obligations.GetObligationRequest{}
	if id != "" {
		req.Identifier = &obligations.GetObligationRequest_Id{
			Id: id,
		}
	} else {
		req.Identifier = &obligations.GetObligationRequest_Fqn{
			Fqn: fqn,
		}
	}

	resp, err := h.sdk.Obligations.GetObligation(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetObligation(), nil
}

func (h Handler) ListObligations(ctx context.Context, limit, offset int32, namespace string) ([]*policy.Obligation, *policy.PageResponse, error) {
	req := &obligations.ListObligationsRequest{
		Pagination: &policy.PageRequest{
			Limit:  limit,
			Offset: offset,
		},
	}
	if namespace != "" {
		_, err := uuid.Parse(namespace)
		if err != nil {
			req.NamespaceIdentifier = &obligations.ListObligationsRequest_Fqn{
				Fqn: namespace,
			}
		} else {
			req.NamespaceIdentifier = &obligations.ListObligationsRequest_Id{
				Id: namespace,
			}
		}
	}
	resp, err := h.sdk.Obligations.ListObligations(ctx, req)
	if err != nil {
		return nil, nil, err
	}

	return resp.GetObligations(), resp.GetPagination(), nil
}

func (h Handler) UpdateObligation(ctx context.Context, id, name string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.Obligation, error) {
	res, err := h.sdk.Obligations.UpdateObligation(ctx, &obligations.UpdateObligationRequest{
		Id:                     id,
		Name:                   name,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	})
	if err != nil {
		return nil, err
	}

	return res.GetObligation(), nil
}

func (h Handler) DeleteObligation(ctx context.Context, id, fqn string) error {
	req := &obligations.DeleteObligationRequest{}
	if id != "" {
		req.Identifier = &obligations.DeleteObligationRequest_Id{
			Id: id,
		}
	} else {
		req.Identifier = &obligations.DeleteObligationRequest_Fqn{
			Fqn: fqn,
		}
	}
	_, err := h.sdk.Obligations.DeleteObligation(ctx, req)
	if err != nil {
		return err
	}

	return nil
}

//
// Obligation Values
//

func (h Handler) CreateObligationValue(ctx context.Context, obligation, value string, metadata *common.MetadataMutable) (*policy.ObligationValue, error) {
	req := &obligations.CreateObligationValueRequest{
		Value:    value,
		Metadata: metadata,
	}

	_, err := uuid.Parse(obligation)
	if err != nil {
		req.ObligationIdentifier = &obligations.CreateObligationValueRequest_Fqn{
			Fqn: obligation,
		}
	} else {
		req.ObligationIdentifier = &obligations.CreateObligationValueRequest_Id{
			Id: obligation,
		}
	}

	resp, err := h.sdk.Obligations.CreateObligationValue(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetValue(), nil
}

func (h Handler) GetObligationValue(ctx context.Context, id, fqn string) (*policy.ObligationValue, error) {
	req := &obligations.GetObligationValueRequest{}
	if id != "" {
		req.Identifier = &obligations.GetObligationValueRequest_Id{
			Id: id,
		}
	} else {
		req.Identifier = &obligations.GetObligationValueRequest_Fqn{
			Fqn: fqn,
		}
	}

	resp, err := h.sdk.Obligations.GetObligationValue(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.GetValue(), nil
}

func (h Handler) UpdateObligationValue(ctx context.Context, id, value string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.ObligationValue, error) {
	res, err := h.sdk.Obligations.UpdateObligationValue(ctx, &obligations.UpdateObligationValueRequest{
		Id:                     id,
		Value:                  value,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	})
	if err != nil {
		return nil, err
	}

	return res.GetValue(), nil
}

func (h Handler) DeleteObligationValue(ctx context.Context, id, fqn string) error {
	req := &obligations.DeleteObligationValueRequest{}
	if id != "" {
		req.Identifier = &obligations.DeleteObligationValueRequest_Id{
			Id: id,
		}
	} else {
		req.Identifier = &obligations.DeleteObligationValueRequest_Fqn{
			Fqn: fqn,
		}
	}
	_, err := h.sdk.Obligations.DeleteObligationValue(ctx, req)
	if err != nil {
		return err
	}

	return nil
}
