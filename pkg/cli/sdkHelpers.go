package cli

import (
	"strconv"
	"time"

	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/platform/protocol/go/common"
	"github.com/opentdf/platform/protocol/go/policy"
)

type SimpleAttribute struct {
	Id        string
	Name      string
	Rule      string
	Values    []string
	Namespace string
	Active    string
	Metadata  map[string]string
}

type SimpleAttributeValue struct {
	Id      string
	FQN     string
	Members []string
	Active  string
}

func ConstructMetadata(m *common.Metadata) map[string]string {
	metadata := map[string]string{
		"Created At": m.CreatedAt.AsTime().Format(time.UnixDate),
		"Updated At": m.UpdatedAt.AsTime().Format(time.UnixDate),
	}

	labels := []string{}
	if m.Labels != nil {
		for k, v := range m.Labels {
			labels = append(labels, k+": "+v)
		}
	}
	metadata["Labels"] = CommaSeparated(labels)
	return metadata
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
		Active:    strconv.FormatBool(a.GetActive().GetValue()),
		Metadata:  ConstructMetadata(a.GetMetadata()),
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
		Active:  strconv.FormatBool(v.Active.GetValue()),
	}
}
