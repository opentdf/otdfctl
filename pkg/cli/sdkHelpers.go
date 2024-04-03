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
	Active    string
}

type SimpleAttributeValue struct {
	Id      string
	FQN     string
	Members []string
}

func GetSimpleAttribute(a *policy.Attribute) SimpleAttribute {
	values := []string{}
	for _, v := range a.GetValues() {
		values = append(values, v.GetValue())
	}

	return SimpleAttribute{
		Id:        a.GetId(),
		Name:      a.GetName(),
		Rule:      handlers.GetAttributeRuleFromAttributeType(a.GetRule()),
		Values:    values,
		Namespace: a.GetNamespace().GetName(),
		Active:    a.GetActive().String(),
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
