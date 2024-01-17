package handlers

import (
	resourceencodingsv1 "github.com/opentdf/opentdf-v2-poc/gen/acre/v1"
	attributesv1 "github.com/opentdf/opentdf-v2-poc/gen/attributes/v1"
	commonv1 "github.com/opentdf/opentdf-v2-poc/gen/common/v1"
	"github.com/opentdf/tructl/pkg/grpc"
)

func CreateResourceEncoding() {
	client := resourceencodingsv1.NewResourceEncodingsServiceClient(grpc.Conn)

	synonyms := resourceencodingsv1.SynonymRef{
		Synonyms: []*resourceencodingsv1.Synonym{
			{
				Terms: []string{
					"test",
				},
			},
		},
	}

	attrs := []*attributesv1.AttributeValueRef{
		{
			Descriptor_: &commonv1.Descriptor{
				Id:   0,
				Name: "test attr",
				Type: &commonv1.PolicyResourceType_POLICY_RESOURCE_TYPE_ATTRIBUTE_DEFINITION,
			},
		},
	}

	resp, err := client.CreateResourceMapping(grpc.Context, &resourceencodingsv1.CreateResourceMappingRequest{
		Mapping: &resourceencodingsv1.ResourceMapping{
			AttributeValueRef: attrs,
			SynonymRef:        synonyms,
		},
	})
	if err != nil {
		panic(err)
	}
	return resp
}

func GetResourceEncoding() {

}

func ListResourceEncodings() {

}

func UpdateResourceEncoding() {

}

func DeleteResourceEncoding() {

}
