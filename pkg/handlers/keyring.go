package handlers

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/zalando/go-keyring"
)

type Keyring struct {
	endpoint string
}

func NewKeyring(endpoint string) *Keyring {
	return &Keyring{
		endpoint: endpoint,
	}
}

func (k *Keyring) getService() string {
	return OTDFCTL_KEYRING_SERVICE + "-" + k.endpoint
}

func (k *Keyring) get(key string, value *map[string]interface{}) error {
	s, err := keyring.Get(k.getService(), key)
	if err != nil {
		return err
	}
	return json.NewDecoder(bytes.NewReader([]byte(s))).Decode(value)
}

func (k *Keyring) set(key string, value any) error {
	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(value); err != nil {
		return err
	}
	return keyring.Set(k.getService(), key, b.String())
}

func (k *Keyring) delete(key string) error {
	return keyring.Delete(k.getService(), key)
}

func (k *Keyring) GetClientCredentials() (ClientCredentials, error) {
	var v map[string]interface{}
	var c ClientCredentials

	if err := k.get(OTDFCTL_KEYRING_CLIENT_CREDENTIALS, &v); err != nil {
		return c, err
	} else if v == nil {
		return c, errors.New("client credentials not found")
	}
	if _, ok := v["clientId"]; !ok {
		return c, errors.New("client_id not found")
	}
	c.ClientId = v["clientId"].(string)

	if _, ok := v["clientSecret"]; !ok {
		return c, errors.New("client_secret not found")
	}
	c.ClientSecret = v["clientSecret"].(string)

	return c, nil
}

func (k *Keyring) SetClientCredentials(c ClientCredentials) error {
	return k.set(OTDFCTL_KEYRING_CLIENT_CREDENTIALS, c)
}

func (k *Keyring) DeleteClientCredentials() error {
	return k.delete(OTDFCTL_KEYRING_CLIENT_CREDENTIALS)
}
