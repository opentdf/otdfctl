package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"time"

	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/google/uuid"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/sdk"
	oidcrp "github.com/zitadel/oidc/v3/pkg/client/rp"
	oidcCLI "github.com/zitadel/oidc/v3/pkg/client/rp/cli"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
)

var (
	ErrParsingAccessToken     = errors.New("failed to parse access token")
	ErrPlatformConfigNotFound = errors.New("platform configuration not found")
)

const (
	AuthCallbackPath             = "/callback"
	AuthCodeFlowPort             = "9000"
	AUTH_TYPE_CLIENT_CREDENTIALS = "client-credentials"
	AUTH_TYPE_ACCESS_TOKEN       = "access-token"
)

type ClientCredentials struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type OidcClientCredentials struct {
	ClientID     string
	ClientSecret string
	isPublic     bool
}

type JWTClaims struct {
	Expiration int64 `json:"exp"`
}

type AuthCredentials struct {
	AuthType string `json:"authType"`
	ClientID string `json:"clientId"`
	// Used for client credentials
	ClientSecret string                      `json:"clientSecret,omitempty"`
	AccessToken  *AuthCredentialsAccessToken `json:"accessToken,omitempty"`
}

type AuthCredentialsAccessToken struct {
	PublicClientID string `json:"publicClientId"`
	AccessToken    string `json:"accessToken"`
	RefreshToken   string `json:"refreshToken"`
	Expiration     int64  `json:"expiration"`
}

// Parse the JSON and return the client ID and secret
func GetClientCredsFromJSON(credsJSON []byte) (ClientCredentials, error) {
	creds := ClientCredentials{}
	if err := json.Unmarshal(credsJSON, &creds); err != nil {
		return creds, errors.Join(errors.New("failed to decode creds JSON"), err)
	}

	return creds, nil
}

func BuildToken(c *AuthCredentials) *oauth2.Token {
	return &oauth2.Token{
		AccessToken:  c.AccessToken.AccessToken,
		Expiry:       time.Unix(c.AccessToken.Expiration, 0),
		RefreshToken: c.AccessToken.RefreshToken,
	}
}

func ParseClaimsJWT(accessToken string) (JWTClaims, error) {
	c := JWTClaims{}
	jwt, err := jwt.ParseSigned(accessToken)
	if err != nil {
		return c, errors.Join(ErrParsingAccessToken, err)
	}
	if err := jwt.UnsafeClaimsWithoutVerification(&c); err != nil {
		return c, errors.Join(ErrParsingAccessToken, err)
	}
	return c, nil
}

// Retrieves credentials by reading specified file
func GetClientCredsFromFile(filepath string) (ClientCredentials, error) {
	creds := ClientCredentials{}
	f, err := os.Open(filepath)
	if err != nil {
		return creds, errors.Join(errors.New("failed to open creds file"), err)
	}
	defer f.Close()

	if err := json.NewDecoder(f).Decode(&creds); err != nil {
		return creds, errors.Join(errors.New("failed to decode creds file"), err)
	}

	return creds, nil
}

// Uses the OAuth2 client credentials flow to obtain a token.
func GetTokenWithClientCreds(ctx context.Context, endpoint string, clientID string, clientSecret string, tlsNoVerify bool) (*oauth2.Token, error) {
	rp, err := newOidcRelyingParty(ctx, endpoint, tlsNoVerify, OidcClientCredentials{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	})
	if err != nil {
		return nil, err
	}
	return oidcrp.ClientCredentials(ctx, rp, url.Values{})
}

const (
	keyLength       = 16
	fiveSecDuration = 5 * time.Second
)

// Facilitates an auth code PKCE flow to obtain OIDC tokens.
// Spawns a local server to handle the callback and opens a browser window in each respective OS.
func Login(ctx context.Context, platformEndpoint, tokenURL, authURL, publicClientID string) (*oauth2.Token, error) {
	// Generate random hash and encryption keys for cookie handling
	hashKey := make([]byte, keyLength)
	encryptKey := make([]byte, keyLength)

	_, err := rand.Read(hashKey)
	if err != nil {
		return nil, err
	}

	_, err = rand.Read(encryptKey)
	if err != nil {
		return nil, err
	}

	conf := &oauth2.Config{
		ClientID:    publicClientID,
		Scopes:      []string{"openid", "profile", "email"},
		RedirectURL: fmt.Sprintf("http://localhost:%s%s", AuthCodeFlowPort, AuthCallbackPath),
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	cookiehandler := httphelper.NewCookieHandler(hashKey, encryptKey)

	relyingParty, err := oidcrp.NewRelyingPartyOAuth(conf,
		// allow cookie handling for PKCE
		oidcrp.WithCookieHandler(cookiehandler),
		// use PKCE
		oidcrp.WithPKCE(cookiehandler),
		// allow IAT claim offset of 5 seconds
		oidcrp.WithVerifierOpts(oidcrp.WithIssuedAtOffset(fiveSecDuration)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create relying party: %w", err)
	}
	stateProvider := func() string {
		return uuid.New().String()
	}
	tok := oidcCLI.CodeFlow[*oidc.IDTokenClaims](ctx, relyingParty, AuthCallbackPath, AuthCodeFlowPort, stateProvider)
	return &oauth2.Token{
		AccessToken:  tok.Token.AccessToken,
		TokenType:    tok.Token.TokenType,
		RefreshToken: tok.Token.RefreshToken,
		Expiry:       tok.Token.Expiry,
	}, nil
}

// Logs in using the auth code PKCE flow driven by the platform well-known idP OIDC configuration.
func LoginWithPKCE(ctx context.Context, host, publicClientID string, tlsNoVerify bool) (*oauth2.Token, string, error) {
	pc, err := getPlatformConfiguration(host, publicClientID, tlsNoVerify)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get platform configuration: %w", err)
	}

	tok, err := Login(ctx, host, pc.tokenEndpoint, pc.authzEndpoint, pc.publicClientID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to login: %w", err)
	}

	return tok, pc.publicClientID, nil
}

// Revokes the access token
func RevokeAccessToken(ctx context.Context, endpoint, publicClientID, refreshToken string, tlsNoVerify bool) error {
	rp, err := newOidcRelyingParty(ctx, endpoint, tlsNoVerify, OidcClientCredentials{
		ClientID: publicClientID,
		isPublic: true,
	})
	if err != nil {
		return err
	}
	return oidcrp.RevokeToken(ctx, rp, refreshToken, "refresh_token")
}

func newOidcRelyingParty(ctx context.Context, endpoint string, tlsNoVerify bool, clientCreds OidcClientCredentials) (oidcrp.RelyingParty, error) {
	if clientCreds.ClientID == "" {
		return nil, errors.New("client ID is required")
	}
	if clientCreds.ClientSecret == "" && !clientCreds.isPublic {
		return nil, errors.New("client secret is required")
	}
	if clientCreds.ClientSecret != "" && clientCreds.isPublic {
		return nil, errors.New("client secret must be empty for public clients")
	}

	var pcClient string
	if clientCreds.isPublic {
		pcClient = clientCreds.ClientID
	}

	pc, err := getPlatformConfiguration(endpoint, pcClient, tlsNoVerify)
	if err != nil {
		if errors.Is(err, sdk.ErrPlatformConfigFailed) {
			return nil, ErrPlatformConfigNotFound
		}
		return nil, err
	}

	return oidcrp.NewRelyingPartyOIDC(
		ctx,
		pc.issuer,
		clientCreds.ClientID,
		clientCreds.ClientSecret,
		"",
		nil,
		oidcrp.WithHTTPClient(utils.NewHttpClient(tlsNoVerify)),
	)
}

type platformConfiguration struct {
	issuer         string
	authzEndpoint  string
	tokenEndpoint  string
	publicClientID string
}

func getPlatformConfiguration(endpoint, publicClientID string, tlsNoVerify bool) (platformConfiguration, error) {
	c := platformConfiguration{}

	normalized, err := utils.NormalizeEndpoint(endpoint)
	if err != nil {
		return c, err
	}

	opts := []sdk.Option{}
	if tlsNoVerify {
		opts = append(opts, sdk.WithInsecureSkipVerifyConn())
	}

	if normalized.Scheme == "http" {
		opts = append(opts, sdk.WithInsecurePlaintextConn())
	}

	s, err := sdk.New(normalized.String(), opts...)
	if err != nil {
		return c, err
	}

	var e error
	c.issuer, e = s.PlatformConfiguration.Issuer()
	if e != nil {
		err = errors.Join(err, sdk.ErrPlatformIssuerNotFound)
	}

	c.authzEndpoint, e = s.PlatformConfiguration.AuthzEndpoint()
	if e != nil {
		err = errors.Join(err, sdk.ErrPlatformAuthzEndpointNotFound)
	}

	c.tokenEndpoint, e = s.PlatformConfiguration.TokenEndpoint()
	if e != nil {
		err = errors.Join(err, sdk.ErrPlatformTokenEndpointNotFound)
	}

	c.publicClientID = publicClientID
	if c.publicClientID == "" {
		c.publicClientID, e = s.PlatformConfiguration.PublicClientID()
		if e != nil {
			err = errors.Join(err, sdk.ErrPlatformPublicClientIDNotFound)
		}
	}

	if err != nil {
		return c, errors.Join(err, ErrPlatformConfigNotFound)
	}

	return c, nil
}
