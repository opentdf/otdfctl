package handlers

import (
	"github.com/opentdf/opentdf-v2-poc/sdk/acre"
	"github.com/opentdf/opentdf-v2-poc/sdk/attributes"
	"github.com/opentdf/opentdf-v2-poc/sdk/common"
)

type ResourceEncoding struct {
	Id          string
	AttributeId string
	Terms       []string
}

func (h *Handler) CreateResourceEncoding(attributeId int, terms []string) (ResourceEncoding, error) {
	synonyms := acre.SynonymRef{
		Ref: &acre.SynonymRef_Synonyms{
			Synonyms: &acre.Synonyms{
				Terms: terms,
			},
		},
	}

	attrs := attributes.AttributeValueReference{
		Ref: &attributes.AttributeValueReference_Descriptor_{
			Descriptor_: &common.ResourceDescriptor{
				Id:        int32(attributeId),
				Name:      "test attr xxxx",
				Namespace: "demo.com",
				Type:      common.PolicyResourceType_POLICY_RESOURCE_TYPE_ATTRIBUTE_DEFINITION,
			},
		},
	}

	_, err := h.sdk.ResourceEncoding.CreateResourceMapping(h.ctx, &acre.CreateResourceMappingRequest{
		Mapping: &acre.ResourceMapping{
			AttributeValueRef: &attrs,
			SynonymRef:        &synonyms,
		},
	})
	if err != nil {
		return ResourceEncoding{}, err
	}

	return ResourceEncoding{}, nil
}

func GetResourceEncoding() {

}

func ListResourceEncodings() {

}

func UpdateResourceEncoding() {

}

func DeleteResourceEncoding() {

}
