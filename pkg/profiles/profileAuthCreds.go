package profiles

const (
	PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS = "client-credentials"
	PROFILE_AUTH_TYPE_ACCESS_TOKEN       = "access-token"
)

type AuthCredentials struct {
	AuthType string `json:"authType" msgpack:"t"`
	ClientId string `json:"clientId,omitempty" msgpack:"cid,omitempty"`
	// Used for client credentials
	ClientSecret string                     `json:"clientSecret,omitempty" msgpack:"cs,omitempty"`
	AccessToken  AuthCredentialsAccessToken `json:"accessToken,omitempty" msgpack:"ac,omitempty"`
}

type AuthCredentialsAccessToken struct {
	PublicClientID string `json:"publicClientId" msgpack:"pcid"`
	AccessToken    string `json:"accessToken" msgpack:"at"`
	RefreshToken   string `json:"refreshToken,omitempty" msgpack:"rt,omitempty"`
	Expiration     int64  `json:"expiration" msgpack:"e"`
}

func (p *ProfileStore) GetAuthCredentials() AuthCredentials {
	return p.config.AuthCredentials
}

func (p *ProfileStore) SetAuthCredentials(authCredentials AuthCredentials) error {
	p.config.AuthCredentials = authCredentials
	return p.Save()
}
