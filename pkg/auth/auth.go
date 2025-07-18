package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-jose/go-jose/v3/jwt"
	"github.com/google/uuid"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/opentdf/otdfctl/pkg/utils"
	"github.com/opentdf/platform/sdk"
	oidcrp "github.com/zitadel/oidc/v3/pkg/client/rp"
	oidcCLI "github.com/zitadel/oidc/v3/pkg/client/rp/cli"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
)

const authCallbackPath = "/callback"

type ClientCredentials struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type platformConfiguration struct {
	issuer        string
	authzEndpoint string
	tokenEndpoint string
}

type oidcClientCredentials struct {
	clientID     string
	clientSecret string
	isPublic     bool
}

type JWTClaims struct {
	Expiration int64 `json:"exp"`
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

// Parse the JSON and return the client ID and secret
func GetClientCredsFromJSON(credsJSON []byte) (ClientCredentials, error) {
	creds := ClientCredentials{}
	if err := json.Unmarshal(credsJSON, &creds); err != nil {
		return creds, errors.Join(errors.New("failed to decode creds JSON"), err)
	}

	return creds, nil
}

func getPlatformConfiguration(endpoint string, tlsNoVerify bool) (platformConfiguration, error) {
	c := platformConfiguration{}

	normalized, err := utils.NormalizeEndpoint(endpoint)
	if err != nil {
		return c, err
	}

	opts := []sdk.Option{sdk.WithConnectionValidation()}
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

	if err != nil {
		return c, errors.Join(err, ErrProfileCredentialsNotFound)
	}

	return c, nil
}

func buildToken(c *profiles.AuthCredentials) *oauth2.Token {
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

func GetSDKAuthOptionFromProfile(profile *profiles.ProfileStore) (sdk.Option, error) {
	c := profile.GetAuthCredentials()

	switch c.AuthType {
	case profiles.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS:
		return sdk.WithClientCredentials(c.ClientId, c.ClientSecret, nil), nil
	case profiles.PROFILE_AUTH_TYPE_ACCESS_TOKEN:
		tokenSource := oauth2.StaticTokenSource(buildToken(&c))
		return sdk.WithOAuthAccessTokenSource(tokenSource), nil
	default:
		return nil, ErrInvalidAuthType
	}
}

func ValidateProfileAuthCredentials(ctx context.Context, profile *profiles.ProfileStore) error {
	c := profile.GetAuthCredentials()
	switch c.AuthType {
	case "":
		return ErrProfileCredentialsNotFound
	case profiles.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS:
		_, err := GetTokenWithClientCreds(ctx, profile.GetEndpoint(), c.ClientId, c.ClientSecret, profile.GetTLSNoVerify())
		if err != nil {
			return err
		}
		return nil
	case profiles.PROFILE_AUTH_TYPE_ACCESS_TOKEN:
		if !buildToken(&c).Valid() {
			return ErrAccessTokenExpired
		}
	default:
		return ErrInvalidAuthType
	}
	return nil
}

func GetTokenWithProfile(ctx context.Context, profile *profiles.ProfileStore) (*oauth2.Token, error) {
	c := profile.GetAuthCredentials()
	switch c.AuthType {
	case profiles.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS:
		return GetTokenWithClientCreds(ctx, profile.GetEndpoint(), c.ClientId, c.ClientSecret, profile.GetTLSNoVerify())
	case profiles.PROFILE_AUTH_TYPE_ACCESS_TOKEN:
		return buildToken(&c), nil
	default:
		return nil, ErrInvalidAuthType
	}
}

// Uses the OAuth2 client credentials flow to obtain a token.
func GetTokenWithClientCreds(ctx context.Context, endpoint string, clientId string, clientSecret string, tlsNoVerify bool) (*oauth2.Token, error) {
	rp, err := newOidcRelyingParty(ctx, endpoint, tlsNoVerify, oidcClientCredentials{
		clientID:     clientId,
		clientSecret: clientSecret,
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

// GetFreePort returns an available TCP port on localhost.
// The function works by asking the operating system to assign
// a free port (by using port 0), then returns that assigned port.
func GetFreePort() (int, error) {
	// Create a listener on localhost with port 0 (OS will assign a free port)
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		return 0, fmt.Errorf("failed to find available port: %w", err)
	}

	// Make sure we release the port when done
	defer listener.Close()

	// Get the address information from the listener
	addr, ok := listener.Addr().(*net.TCPAddr)
	if !ok {
		return 0, fmt.Errorf("failed to get TCP address from listener")
	}

	// Return the port that was assigned
	return addr.Port, nil
}

// Facilitates an auth code PKCE flow to obtain OIDC tokens.
// Spawns a local server to handle the callback and opens a browser window in each respective OS.
func Login(ctx context.Context, platformEndpoint, tokenURL, authURL, publicClientID, authCodeFlowPort string) (*oauth2.Token, error) {
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

	if strings.TrimSpace(authCodeFlowPort) == "" {
		port, err := GetFreePort()
		if err != nil {
			return nil, fmt.Errorf("failed to find available port for auth code flow: %w", err)
		}
		authCodeFlowPort = strconv.Itoa(port)
	}

	conf := &oauth2.Config{
		ClientID:    publicClientID,
		Scopes:      []string{"openid", "profile", "email"},
		RedirectURL: fmt.Sprintf("http://localhost:%s%s", authCodeFlowPort, authCallbackPath),
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
	tok := oidcCLI.CodeFlow[*oidc.IDTokenClaims](ctx, relyingParty, authCallbackPath, authCodeFlowPort, stateProvider)
	return &oauth2.Token{
		AccessToken:  tok.Token.AccessToken,
		TokenType:    tok.Token.TokenType,
		RefreshToken: tok.Token.RefreshToken,
		Expiry:       tok.Token.Expiry,
	}, nil
}

// Logs in using the auth code PKCE flow driven by the platform well-known idP OIDC configuration.
func LoginWithPKCE(ctx context.Context, host, clientID string, tlsNoVerify bool, port string) (*oauth2.Token, error) {
	pc, err := getPlatformConfiguration(host, tlsNoVerify)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform configuration: %w", err)
	}

	tok, err := Login(ctx, host, pc.tokenEndpoint, pc.authzEndpoint, clientID, port)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return tok, nil
}

// Revokes the access token
func RevokeAccessToken(ctx context.Context, endpoint, clientID, refreshToken string, tlsNoVerify bool) error {
	rp, err := newOidcRelyingParty(ctx, endpoint, tlsNoVerify, oidcClientCredentials{
		clientID: clientID,
		isPublic: true,
	})
	if err != nil {
		return err
	}
	return oidcrp.RevokeToken(ctx, rp, refreshToken, "refresh_token")
}

func newOidcRelyingParty(ctx context.Context, endpoint string, tlsNoVerify bool, clientCreds oidcClientCredentials) (oidcrp.RelyingParty, error) {
	if clientCreds.clientID == "" {
		return nil, errors.New("client ID is required")
	}
	if clientCreds.clientSecret == "" && !clientCreds.isPublic {
		return nil, errors.New("client secret is required")
	}
	if clientCreds.clientSecret != "" && clientCreds.isPublic {
		return nil, errors.New("client secret must be empty for public clients")
	}

	pc, err := getPlatformConfiguration(endpoint, tlsNoVerify)
	if err != nil {
		return nil, err
	}

	return oidcrp.NewRelyingPartyOIDC(
		ctx,
		pc.issuer,
		clientCreds.clientID,
		clientCreds.clientSecret,
		"",
		nil,
		oidcrp.WithHTTPClient(utils.NewHttpClient(tlsNoVerify)),
	)
}
