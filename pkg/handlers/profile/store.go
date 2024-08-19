package profile

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/zalando/go-keyring"
)

// TODO: update the store to use alternative storage methods besides keyring

type Store struct {
	namespace string
	key       string
}

func NewStore(namespace string, key string) *Store {
	return &Store{
		namespace: namespace,
		key:       key,
	}
}

func (k *Store) Exists() bool {
	s, err := keyring.Get(k.namespace, k.key)
	return err == nil && s != ""
}

func (k *Store) Get(value interface{}) error {
	s, err := keyring.Get(k.namespace, k.key)
	if err != nil {
		return err
	}
	return json.NewDecoder(bytes.NewReader([]byte(s))).Decode(value)
}

func (k *Store) Set(value interface{}) error {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(value); err != nil {
		return err
	}
	return keyring.Set(k.namespace, k.key, b.String())
}

func (k *Store) Delete() error {
	return keyring.Delete(k.namespace, k.key)
}

func fillStruct(m map[string]interface{}, s interface{}) error {
	structValue := reflect.ValueOf(s).Elem()

	for name, value := range m {
		name := strings.ToUpper(name[:1]) + name[1:]
		structFieldValue := structValue.FieldByName(name)

		if !structFieldValue.IsValid() {
			return fmt.Errorf("No such field: %s in obj", name)
		}

		if !structFieldValue.CanSet() {
			return fmt.Errorf("Cannot set %s field value", name)
		}

		val := reflect.ValueOf(value)
		fmt.Print(val)
		if structFieldValue.Type() != val.Type() {
			fmt.Printf("%s, %s", structFieldValue.Type(), val.Type())
			return errors.New("Provided value type didn't match obj field type")
		}

		structFieldValue.Set(val)
	}
	return nil
}
