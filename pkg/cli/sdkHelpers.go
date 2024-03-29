package cli

import (
	"github.com/opentdf/platform/protocol/go/policy"
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

func GetSimpleAttribute(a *policy.Attribute) SimpleAttribute {
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

func GetSimpleAttributeValue(v *policy.Value) SimpleAttributeValue {
	memberIds := []string{}
	for _, m := range v.Members {
		memberIds = append(memberIds, m.Id)
	}
	return SimpleAttributeValue{
		Id:      v.Id,
		FQN:     v.Fqn,
		Members: memberIds,
	}
}
