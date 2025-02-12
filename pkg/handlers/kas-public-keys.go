package handlers

import (
	"errors"

	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
	"github.com/opentdf/platform/protocol/go/policy/kasregistry"
	"github.com/opentdf/platform/protocol/go/policy/unsafe"
	"google.golang.org/grpc/status"
)

func (h Handler) CreatePublicKey(kas, pk, kid, alg string, metadata *common.MetadataMutable) (*policy.Key, error) {
	// Check if alg is valid
	algEnum, err := algToEnum(alg)
	if err != nil {
		return nil, err
	}

	// Key ID can't be more than 32 characters
	if len(kid) > 32 {
		return nil, errors.New("key id must be less than 32 characters")
	}

	// Get KAS UUID if it's not a UUID

	if !utils.IsUUID(kas) {
		k, err := h.GetKasRegistryEntry(kas)
		if err != nil {
			return nil, err
		}
		kas = k.GetId()
	}

	// Create the public key
	resp, err := h.sdk.KeyAccessServerRegistry.CreatePublicKey(h.ctx, &kasregistry.CreatePublicKeyRequest{
		KasId: kas,
		Key: &policy.KasPublicKey{
			Kid: kid,
			Alg: algEnum,
			Pem: pk,
		},
		Metadata: metadata,
	})
	if err != nil {
		s := status.Convert(err)
		return nil, errors.New(s.Message())
	}

	return h.GetPublicKey(resp.GetKey().GetId())
}

func (h Handler) UpdatePublicKey(id string, metadata *common.MetadataMutable, behavior common.MetadataUpdateEnum) (*policy.Key, error) {
	resp, err := h.sdk.KeyAccessServerRegistry.UpdatePublicKey(h.ctx, &kasregistry.UpdatePublicKeyRequest{
		Id:                     id,
		Metadata:               metadata,
		MetadataUpdateBehavior: behavior,
	})
	if err != nil {
		s := status.Convert(err)
		return nil, errors.New(s.Message())
	}

	return resp.GetKey(), nil
}

func (h Handler) GetPublicKey(id string) (*policy.Key, error) {
	resp, err := h.sdk.KeyAccessServerRegistry.GetPublicKey(h.ctx, &kasregistry.GetPublicKeyRequest{
		Identifier: &kasregistry.GetPublicKeyRequest_Id{Id: id},
	})
	if err != nil {
		return nil, err
	}
	return resp.GetKey(), nil
}

func (h Handler) ListPublicKeys(kas string, offset, limit int32) ([]*policy.Key, *policy.PageResponse, error) {
	req := &kasregistry.ListPublicKeysRequest{
		Pagination: &policy.PageRequest{
			Offset: offset,
			Limit:  limit,
		},
	}

	switch {
	case utils.IsUUID(kas):
		req.KasFilter = &kasregistry.ListPublicKeysRequest_KasId{KasId: kas}
	case utils.IsURI(kas):
		req.KasFilter = &kasregistry.ListPublicKeysRequest_KasUri{KasUri: kas}
	case kas != "":
		req.KasFilter = &kasregistry.ListPublicKeysRequest_KasName{KasName: kas}
	}

	resp, err := h.sdk.KeyAccessServerRegistry.ListPublicKeys(h.ctx, req)
	if err != nil {
		return nil, nil, err
	}
	return resp.GetKeys(), resp.GetPagination(), nil
}

func (h Handler) ListPublicKeyMappings(kas, pkID string, offset, limit int32) ([]*kasregistry.ListPublicKeyMappingResponse_PublicKeyMapping, *policy.PageResponse, error) {
	req := &kasregistry.ListPublicKeyMappingRequest{
		PublicKeyId: pkID,
		Pagination: &policy.PageRequest{
			Offset: offset,
			Limit:  limit,
		},
	}

	switch {
	case utils.IsUUID(kas):
		req.KasFilter = &kasregistry.ListPublicKeyMappingRequest_KasId{KasId: kas}
	case utils.IsURI(kas):
		req.KasFilter = &kasregistry.ListPublicKeyMappingRequest_KasUri{KasUri: kas}
	case kas != "":
		req.KasFilter = &kasregistry.ListPublicKeyMappingRequest_KasName{KasName: kas}
	}

	resp, err := h.sdk.KeyAccessServerRegistry.ListPublicKeyMapping(h.ctx, req)
	if err != nil {
		return nil, nil, err
	}

	return resp.GetPublicKeyMappings(), resp.GetPagination(), nil
}

func (h Handler) DeactivatePublicKey(id string) error {
	_, err := h.sdk.KeyAccessServerRegistry.DeactivatePublicKey(h.ctx, &kasregistry.DeactivatePublicKeyRequest{Id: id})
	if err != nil {
		s := status.Convert(err)
		return errors.New(s.Message())
	}
	return nil
}

func (h Handler) ActivatePublicKey(id string) error {
	_, err := h.sdk.KeyAccessServerRegistry.ActivatePublicKey(h.ctx, &kasregistry.ActivatePublicKeyRequest{Id: id})
	if err != nil {
		s := status.Convert(err)
		return errors.New(s.Message())
	}
	return nil
}

func (h Handler) UnsafeDeletePublicKey(id string) error {
	_, err := h.sdk.Unsafe.UnsafeDeletePublicKey(h.ctx, &unsafe.UnsafeDeletePublicKeyRequest{Id: id})
	if err != nil {
		s := status.Convert(err)
		return errors.New(s.Message())
	}
	return nil
}

func algToEnum(alg string) (policy.KasPublicKeyAlgEnum, error) {
	switch alg {
	case "rsa:2048":
		return policy.KasPublicKeyAlgEnum_KAS_PUBLIC_KEY_ALG_ENUM_RSA_2048, nil
	case "rsa:4096":
		return policy.KasPublicKeyAlgEnum_KAS_PUBLIC_KEY_ALG_ENUM_RSA_4096, nil
	case "ec:secp256r1":
		return policy.KasPublicKeyAlgEnum_KAS_PUBLIC_KEY_ALG_ENUM_EC_SECP256R1, nil
	case "ec:secp384r1":
		return policy.KasPublicKeyAlgEnum_KAS_PUBLIC_KEY_ALG_ENUM_EC_SECP384R1, nil
	case "ec:secp521r1":
		return policy.KasPublicKeyAlgEnum_KAS_PUBLIC_KEY_ALG_ENUM_EC_SECP521R1, nil
	default:
		return policy.KasPublicKeyAlgEnum_KAS_PUBLIC_KEY_ALG_ENUM_UNSPECIFIED, errors.New("unsupported algorithm. supported algorithms are rsa:2048, rsa:4096, ec:secp256r1, ec:secp384r1, ec:secp521r1")
	}
}
