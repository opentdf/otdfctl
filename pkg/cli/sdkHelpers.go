package cli

import (
	"strings"

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

type SimpleAttributeValue struct {
	Id      string
	FQN     string
	Members []string
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

func GetSimpleAttributeValue(v *attributes.Value) SimpleAttributeValue {
	members := []string{}

	fqn := strings.Join([]string{"v.Attribute.Namespace.Name", "attr", "v.Attribute.Name", "value", v.Value}, "/")

	return SimpleAttributeValue{
		Id:      v.Id,
		FQN:     fqn,
		Members: members,
	}
}
