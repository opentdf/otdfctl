package handlers

import (
	"context"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/kasregistry"
)

func (h Handler) CreateKasKey(
	ctx context.Context,
	kasId string,
	keyId string,
	alg policy.Algorithm,
	mode policy.KeyMode,
	pubKeyCtx []byte,
	privKeyCtx []byte,
	providerConfigId string,
	metadata *common.MetadataMutable,
) (*policy.AsymmetricKey, error) {
	req := kasregistry.CreateKeyRequest{
		KasId:            kasId,
		KeyId:            keyId,
		KeyAlgorithm:     alg,
		KeyMode:          mode,
		PublicKeyCtx:     pubKeyCtx,
		PrivateKeyCtx:    privKeyCtx,
		ProviderConfigId: providerConfigId,
		Metadata:         metadata,
	}

	resp, err := h.sdk.KeyAccessServerRegistry.CreateKey(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.GetKey(), nil
}

func (h Handler) GetKasKey(ctx context.Context, id string, keyId string) (*policy.AsymmetricKey, error) {
	req := kasregistry.GetKeyRequest{}
	if id != "" {
		req.Identifier = &kasregistry.GetKeyRequest_Id{
			Id: id,
		}
	} else if keyId != "" {
		req.Identifier = &kasregistry.GetKeyRequest_KeyId{
			KeyId: keyId,
		}
	}

	resp, err := h.sdk.KeyAccessServerRegistry.GetKey(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.GetKey(), nil
}

func (h Handler) UpdateKasKey(ctx context.Context, id string, status policy.KeyStatus, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.AsymmetricKey, error) {
	req := kasregistry.UpdateKeyRequest{
		Id:                     id,
		KeyStatus:              status,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	}

	resp, err := h.sdk.KeyAccessServerRegistry.UpdateKey(ctx, &req)
	if err != nil {
		return nil, err
	}

	return resp.GetKey(), nil
}

func (h Handler) ListKasKeys(
	ctx context.Context,
	limit, offset int32,
	algorithm policy.Algorithm,
	kasId string,
	kasName string,
	kasUri string) ([]*policy.AsymmetricKey, *policy.PageResponse, error) {
	req := kasregistry.ListKeysRequest{
		Pagination: &policy.PageRequest{
			Limit:  limit,
			Offset: offset,
		},
		KeyAlgorithm: algorithm,
	}

	if kasId != "" {
		req.KasFilter = &kasregistry.ListKeysRequest_KasId{
			KasId: kasId,
		}
	} else if kasName != "" {
		req.KasFilter = &kasregistry.ListKeysRequest_KasName{
			KasName: kasName,
		}
	} else if kasUri != "" {
		req.KasFilter = &kasregistry.ListKeysRequest_KasUri{
			KasUri: kasUri,
		}
	}

	resp, err := h.sdk.KeyAccessServerRegistry.ListKeys(ctx, &req)
	if err != nil {
		return nil, nil, err
	}

	return resp.GetKeys(), resp.GetPagination(), nil
}
