package handlers

import (
	"context"
	"errors"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/kasregistry"
)

func (h Handler) CreateKasKey(
	ctx context.Context,
	kasID string,
	keyID string,
	alg policy.Algorithm,
	mode policy.KeyMode,
	pubKeyCtx *policy.PublicKeyCtx,
	privKeyCtx *policy.PrivateKeyCtx,
	providerConfigID string,
	metadata *common.MetadataMutable,
) (*policy.KasKey, error) {
	req := kasregistry.CreateKeyRequest{
		KasId:            kasID,
		KeyId:            keyID,
		KeyAlgorithm:     alg,
		KeyMode:          mode,
		PublicKeyCtx:     pubKeyCtx,
		PrivateKeyCtx:    privKeyCtx,
		ProviderConfigId: providerConfigID,
		Metadata:         metadata,
	}

	resp, err := h.sdk.KeyAccessServerRegistry.CreateKey(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.GetKasKey(), nil
}

func (h Handler) GetKasKey(ctx context.Context, id string, key *kasregistry.KasKeyIdentifier) (*policy.KasKey, error) {
	req := kasregistry.GetKeyRequest{}
	switch {
	case id != "":
		req.Identifier = &kasregistry.GetKeyRequest_Id{
			Id: id,
		}
	case key != nil:
		req.Identifier = &kasregistry.GetKeyRequest_Key{
			Key: key,
		}
	default:
		return nil, errors.New("id or key must be provided")
	}

	resp, err := h.sdk.KeyAccessServerRegistry.GetKey(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.GetKasKey(), nil
}

func (h Handler) UpdateKasKey(ctx context.Context, id string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.KasKey, error) {
	req := kasregistry.UpdateKeyRequest{
		Id:                     id,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	}

	resp, err := h.sdk.KeyAccessServerRegistry.UpdateKey(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.GetKasKey(), nil
}

func (h Handler) ListKasKeys(
	ctx context.Context,
	limit, offset int32,
	algorithm policy.Algorithm,
	identifier KasIdentifier) ([]*policy.KasKey, *policy.PageResponse, error) {
	req := kasregistry.ListKeysRequest{
		Pagination: &policy.PageRequest{
			Limit:  limit,
			Offset: offset,
		},
		KeyAlgorithm: algorithm,
	}

	switch {
	case identifier.ID != "":
		req.KasFilter = &kasregistry.ListKeysRequest_KasId{
			KasId: identifier.ID,
		}
	case identifier.Name != "":
		req.KasFilter = &kasregistry.ListKeysRequest_KasName{
			KasName: identifier.Name,
		}
	case identifier.URI != "":
		req.KasFilter = &kasregistry.ListKeysRequest_KasUri{
			KasUri: identifier.URI,
		}
	}

	resp, err := h.sdk.KeyAccessServerRegistry.ListKeys(ctx, &req)
	if err != nil {
		return nil, nil, err
	}

	return resp.GetKasKeys(), resp.GetPagination(), nil
}
