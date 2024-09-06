package handlers

import (
	"fmt"
	"log"
	"reflect"

	"github.com/golang-jwt/jwt"
	"github.com/itchyny/gojq"
	"github.com/opentdf/platform/protocol/go/policy"
)

// Recursively process json into a list of jq syntax selectors and their values when applying the jq selector to the input json
func ProcessSubjectContext(subject interface{}, currSelector string, result []*policy.SubjectProperty) ([]*policy.SubjectProperty, error) {
	if currSelector == "" {
		currSelector = "'"
	}
	currType := reflect.TypeOf(subject)

	//nolint:exhaustive // default handles unspecified types as desired
	switch currType.Kind() {
	// maps (structs not supported): add the key to the selector then call on all values
	case reflect.Map:
		for _, key := range reflect.ValueOf(subject).MapKeys() {
			newSelector := fmt.Sprintf("%s.%s", currSelector, key)
			newValue := reflect.ValueOf(subject).MapIndex(key).Interface()
			if r, err := ProcessSubjectContext(newValue, newSelector, result); err != nil {
				return nil, err
			} else {
				result = r
			}
		}
	// lists: invoke on all array values with index added to selector
	case reflect.Array, reflect.Slice:
		for i := 0; i < reflect.ValueOf(subject).Len(); i++ {
			// exists at specific index
			idxSelector := fmt.Sprintf("%s[%d]", currSelector, i)
			newValue := reflect.ValueOf(subject).Index(i).Interface()
			if r, err := ProcessSubjectContext(newValue, idxSelector, result); err != nil {
				return nil, err
			} else {
				result = r
			}
			// if primitive, add selector for if it exists at any index
			if isPrimitive(reflect.TypeOf(newValue).Kind()) {
				anySelector := currSelector + ` | map(tostring) | any(index("` + fmt.Sprintf("%v", newValue) + `"))`
				if r, err := ProcessSubjectContext(true, anySelector, result); err != nil {
					return nil, err
				} else {
					result = r
				}
			}
		}

	// primitives: add the selector and value to the list
	case reflect.String,
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		result = append(result, &policy.SubjectProperty{ExternalSelectorValue: currSelector + "'", ExternalValue: fmt.Sprintf("%v", subject)})

	default:
		return nil, fmt.Errorf("unsupported type %v", currType.Kind())
	}

	return result, nil
}

func isPrimitive(t reflect.Kind) bool {
	//nolint:exhaustive // all primitives are covered
	switch t {
	case reflect.String,
		reflect.Bool,
		reflect.Int,
		reflect.Int8,
		reflect.Int16,
		reflect.Int32,
		reflect.Int64,
		reflect.Uint,
		reflect.Uint8,
		reflect.Uint16,
		reflect.Uint32,
		reflect.Uint64,
		reflect.Uintptr,
		reflect.Float32,
		reflect.Float64,
		reflect.Complex64,
		reflect.Complex128:
		return true
	default:
		return false
	}
}

func TestSubjectContext(subject interface{}, selectors []string) ([]*policy.SubjectProperty, error) {
	// genericize type to avoid panic parsing jwt.MapClaims in gojq
	var sub any
	if _, ok := subject.(jwt.MapClaims); ok {
		subj := make(map[string]interface{})
		subject, subOk := subject.(jwt.MapClaims)
		if !subOk {
			return nil, fmt.Errorf("failed to convert subject to jwt.MapClaims")
		}
		for k, v := range subject {
			subj[k] = v
		}
		sub = subj
	} else {
		sub = subject
	}

	found := []*policy.SubjectProperty{}

	for _, s := range selectors {
		query, err := gojq.Parse(s)
		if err != nil {
			log.Fatalln(err)
		}
		iter := query.Run(sub) // or query.RunWithContext
		for {
			v, ok := iter.Next()
			if !ok {
				break
			}
			if err, ok = v.(error); ok {
				//nolint:errorlint // halt error is a type
				if err, errOk := err.(*gojq.HaltError); errOk && err.Value() == nil {
					break
				}
				// ignore error: we don't have a match but that is not an error state in this case
			} else {
				if v != nil {
					found = append(found, &policy.SubjectProperty{ExternalSelectorValue: "'" + s + "'", ExternalValue: fmt.Sprintf("%v", v)})
				}
			}
		}
	}
	return found, nil
}
