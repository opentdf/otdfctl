package profile

import (
	"bytes"
	"encoding/json"

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
