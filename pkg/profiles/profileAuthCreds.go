package profiles

const (
	PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS = "client-credentials"
	PROFILE_AUTH_TYPE_ACCESS_TOKEN       = "access-token"
)

type AuthCredentials struct {
	AuthType string `json:"authType"`
	ClientId string `json:"clientId"`
	// Used for client credentials
	ClientSecret string                     `json:"clientSecret,omitempty"`
	AccessToken  AuthCredentialsAccessToken `json:"accessToken,omitempty"`
}

type AuthCredentialsAccessToken struct {
	PublicClientID string `json:"publicClientId"`
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	Expiration     int64  `json:"expiration"`
}

func (p *ProfileStore) GetAuthCredentials() AuthCredentials {
	return p.config.AuthCredentials
}

func (p *ProfileStore) SetAuthCredentials(authCredentials AuthCredentials) error {
	p.config.AuthCredentials = authCredentials
	return p.Save()
}
