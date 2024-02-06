package cli

import (
	"github.com/opentdf/opentdf-v2-poc/sdk/attributes"
	"github.com/opentdf/tructl/pkg/handlers"
)

type SimpleAttribute struct {
	Id        string
	Name      string
	Rule      string
	Values    []string
	Namespace string
}

func GetSimpleAttribute(a *attributes.Attribute) SimpleAttribute {
	values := []string{}
	for _, v := range a.Values {
		values = append(values, v.Value)
	}

	return SimpleAttribute{
		Id:        a.Id,
		Name:      a.Name,
		Rule:      handlers.GetAttributeRuleFromAttributeType(a.Rule),
		Values:    values,
		Namespace: a.Namespace.Name,
	}
}
