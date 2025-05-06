package profiles

import (
	"fmt"
	"strconv"

	"github.com/vmihailenco/msgpack/v5"
	"github.com/zalando/go-keyring"
)

type KeyringStore struct {
	namespace string
	key       string
}

var NewKeyringStore NewStoreInterface = func(namespace string, key string) StoreInterface {
	return &KeyringStore{
		namespace: namespace,
		key:       key,
	}
}

func (k *KeyringStore) Exists() bool {
	s, err := keyring.Get(k.namespace, k.key)
	return err == nil && s != ""
}

func (k *KeyringStore) Get(value interface{}) error {
	s, err := keyring.Get(k.namespace, k.key)
	if err != nil {
		return err
	}
	return msgpack.Unmarshal([]byte(s), value)
}

func (k *KeyringStore) Set(value interface{}) error {
	var refreshTokenRemoved bool
	if c, ok := value.(ProfileConfig); ok && c.AuthCredentials.AccessToken.RefreshToken != "" {
		refreshTokenRemoved = true // remove to save size
		fmt.Print("Minimized size: ")
		c.AuthCredentials.AccessToken.RefreshToken = ""
		value = c
	}
	b, err := msgpack.Marshal(value)
	if err != nil {
		return err
	}
	if refreshTokenRemoved {
		fmt.Printf("%s...", strconv.Itoa(len(b)))
	}
	return keyring.Set(k.namespace, k.key, string(b))
}

func (k *KeyringStore) Delete() error {
	return keyring.Delete(k.namespace, k.key)
}
