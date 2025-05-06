package profiles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

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
	return json.NewDecoder(bytes.NewReader([]byte(s))).Decode(value)
}

func (k *KeyringStore) Set(value interface{}) error {
	var refreshTokenRemoved bool
	if c, ok := value.(ProfileConfig); ok && c.AuthCredentials.AccessToken.RefreshToken != "" {
		refreshTokenRemoved = true // remove to save size
		fmt.Print("Removing refresh token, size: ")
		c.AuthCredentials.AccessToken.RefreshToken = ""
		value = c
	}
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(value); err != nil {
		return err
	}
	if refreshTokenRemoved {
		fmt.Println(strconv.Itoa(b.Len()))
	}
	return keyring.Set(k.namespace, k.key, b.String())
}

func (k *KeyringStore) Delete() error {
	return keyring.Delete(k.namespace, k.key)
}
