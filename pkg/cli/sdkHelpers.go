package cli

import (
	"errors"
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

func KeyAlgToEnum(alg string) (policy.Algorithm, error) {
	switch strings.ToLower(alg) {
	case "rsa:2048":
		return policy.Algorithm_ALGORITHM_RSA_2048, nil
	case "rsa:4096":
		return policy.Algorithm_ALGORITHM_RSA_4096, nil
	case "ec:secp256r1":
		return policy.Algorithm_ALGORITHM_EC_P256, nil
	case "ec:secp384r1":
		return policy.Algorithm_ALGORITHM_EC_P384, nil
	case "ec:secp521r1":
		return policy.Algorithm_ALGORITHM_EC_P521, nil
	default:
		return policy.Algorithm_ALGORITHM_UNSPECIFIED, errors.New("invalid algorithm")
	}
}

func KeyEnumToAlg(enum policy.Algorithm) (string, error) {
	switch enum { //nolint:exhaustive // UNSPECIFIED is not needed here
	case policy.Algorithm_ALGORITHM_RSA_2048:
		return "rsa:2048", nil
	case policy.Algorithm_ALGORITHM_RSA_4096:
		return "rsa:4096", nil
	case policy.Algorithm_ALGORITHM_EC_P256:
		return "ec:secp256r1", nil
	case policy.Algorithm_ALGORITHM_EC_P384:
		return "ec:secp384r1", nil
	case policy.Algorithm_ALGORITHM_EC_P521:
		return "ec:secp521r1", nil
	default:
		return "", errors.New("invalid enum algorithm")
	}
}
