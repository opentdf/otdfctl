package cli

import (
	"strconv"
	"strings"
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
	KeyIds    []string
}

type SimpleAttributeValue struct {
	Id       string
	FQN      string
	Active   string
	Metadata map[string]string
}

func ConstructMetadata(m *common.Metadata) map[string]string {
	var metadata map[string]string
	if m == nil {
		return metadata
	}
	metadata = map[string]string{
		"Created At": m.GetCreatedAt().AsTime().Format(time.UnixDate),
		"Updated At": m.GetUpdatedAt().AsTime().Format(time.UnixDate),
	}

	labels := []string{}
	if m.Labels != nil {
		for k, v := range m.GetLabels() {
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
	//keyIds := make([]string, len(a.GetKasKeys()))
	//for i, k := range a.GetKasKeys() {
	//	if k.GetKey() != nil && k.GetKey().GetId() != "" {
	//		keyIds[i] = k.GetKey().GetId()
	//	}
	//}

	return SimpleAttribute{
		Id:        a.GetId(),
		Name:      a.GetName(),
		Rule:      handlers.GetAttributeRuleFromAttributeType(a.GetRule()),
		Values:    values,
		Namespace: a.GetNamespace().GetName(),
		Active:    strconv.FormatBool(a.GetActive().GetValue()),
		Metadata:  ConstructMetadata(a.GetMetadata()),
		//KeyIds:    keyIds,
	}
}

func GetSimpleAttributeValue(v *policy.Value) SimpleAttributeValue {
	return SimpleAttributeValue{
		Id:       v.GetId(),
		FQN:      v.GetFqn(),
		Active:   strconv.FormatBool(v.GetActive().GetValue()),
		Metadata: ConstructMetadata(v.GetMetadata()),
	}
}

func GetSimpleRegisteredResourceValues(v []*policy.RegisteredResourceValue) []string {
	values := make([]string, len(v))
	for i, val := range v {
		values[i] = val.GetValue()
	}
	return values
}

func GetSimpleRegisteredResourceActionAttributeValues(v []*policy.RegisteredResourceValue_ActionAttributeValue) []string {
	values := make([]string, len(v))
	sb := new(strings.Builder)

	for i, val := range v {
		action := val.GetAction()
		attrVal := val.GetAttributeValue()

		sb.WriteString(action.GetName())
		sb.WriteString(" -> ")
		sb.WriteString(attrVal.GetFqn())

		values[i] = sb.String()
		sb.Reset()
	}

	return values
}
