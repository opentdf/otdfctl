package profile

import "errors"

const PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS = "client-credentials"

type AuthCredentials struct {
	AuthType          string            `json:"authType"`
	ClientCredentials ClientCredentials `json:"clientCredentials,omitempty"`
}

type ClientCredentials struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

func (p *ProfileStore) GetAuthCredentials() AuthCredentials {
	return p.config.AuthCredentials
}

func (p *ProfileStore) SetAuthCredentials(authCredentials AuthCredentials) error {
	if authCredentials.AuthType != PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS {
		return errors.New("invalid auth type")
	}

	p.config.AuthCredentials = authCredentials
	return p.Save()
}
