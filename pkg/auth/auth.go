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

const (
	authCallbackPath = "/callback"
	authCodeFlowPort = "9000"
)

type ClientCredentials struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

type platformConfiguration struct {
	issuer         string
	authzEndpoint  string
	tokenEndpoint  string
	publicClientID string
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

func GetClientCreds(endpoint string, file string, credsJSON []byte) (ClientCredentials, error) {
	if file != "" {
		return GetClientCredsFromFile(file)
	}
	if len(credsJSON) > 0 {
		return GetClientCredsFromJSON(credsJSON)
	}
	return ClientCredentials{}, errors.New("no client credentials provided")
}

func getPlatformConfiguration(endpoint string, tlsNoVerify bool) (platformConfiguration, error) {
	c := platformConfiguration{}

	e, err := utils.NormalizeEndpoint(endpoint)
	if err != nil {
		return c, errors.Join(ErrNoPlatformConfiguration, err)
	}

	opts := []sdk.Option{}
	if tlsNoVerify {
		opts = append(opts, sdk.WithInsecureSkipVerifyConn())
	}

	if e.Scheme == "http" {
		opts = append(opts, sdk.WithInsecurePlaintextConn())
	}

	s, err := sdk.New(e.String(), opts...)
	if err != nil {
		return c, errors.Join(ErrNoPlatformConfiguration, err)
	}

	errs := []error{}
	c.issuer, err = s.PlatformConfiguration.Issuer()
	if err != nil {
		errs = append(errs, err)
	}

	c.authzEndpoint, err = s.PlatformConfiguration.AuthzEndpoint()
	if err != nil {
		errs = append(errs, err)
	}

	c.tokenEndpoint, err = s.PlatformConfiguration.TokenEndpoint()
	if err != nil {
		errs = append(errs, err)
	}

	// TODO fix error
	// c.publicClientID, err = s.PlatformConfiguration.PublicClientID()
	// if err != nil {
	// 	errs = append(errs, err)
	// }

	if len(errs) > 0 {
		errs = append([]error{ErrNoPlatformConfiguration}, errs...)
		return c, errors.Join(errs...)
	}

	return c, nil
}

// func GetAccessTokenFromProfile() {}
func GetSDKAuthOptionFromProfile(profile *profiles.ProfileStore) (sdk.Option, error) {
	c := profile.GetAuthCredentials()

	switch c.AuthType {
	case profiles.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS:
		return sdk.WithClientCredentials(c.ClientId, c.ClientSecret, nil), nil
	// case profiles.PROFILE_AUTH_TYPE_ACCESS_TOKEN:
	// 	return sdk.WithOAuthAccessTokenSource(o.authClientCredentials.AccessToken.AccessToken), nil
	default:
		return nil, ErrInvalidAuthType
	}
}

func ValidateProfileAuthCredentials(ctx context.Context, profile *profiles.ProfileStore) error {
	c := profile.GetAuthCredentials()
	switch c.AuthType {
	case profiles.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS:
		_, err := GetTokenWithClientCreds(ctx, profile.GetEndpoint(), c.ClientId, c.ClientSecret, profile.GetTLSNoVerify())
		if err != nil {
			return err
		}
		return nil
	// case profiles.PROFILE_AUTH_TYPE_ACCESS_TOKEN:
	// 	return sdk.WithOAuthAccessTokenSource(o.authClientCredentials.AccessToken.AccessToken), nil
	default:
		return ErrInvalidAuthType
	}
}

func GetTokenWithProfile(ctx context.Context, profile *profiles.ProfileStore) (*oauth2.Token, error) {
	c := profile.GetAuthCredentials()
	switch c.AuthType {
	case profiles.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS:
		return GetTokenWithClientCreds(ctx, profile.GetEndpoint(), c.ClientId, c.ClientSecret, profile.GetTLSNoVerify())
	// case profiles.PROFILE_AUTH_TYPE_ACCESS_TOKEN:
	// 	return sdk.WithOAuthAccessTokenSource(o.authClientCredentials.AccessToken.AccessToken), nil
	default:
		return nil, ErrInvalidAuthType
	}
}

// Uses the OAuth2 client credentials flow to obtain a token.
func GetTokenWithClientCreds(ctx context.Context, endpoint string, clientId string, clientSecret string, tlsNoVerify bool) (*oauth2.Token, error) {
	pc, err := getPlatformConfiguration(endpoint, tlsNoVerify)
	if err != nil {
		return nil, err
	}

	rp, err := oidcrp.NewRelyingPartyOIDC(ctx, pc.issuer, clientId, clientSecret, "", []string{"email"})
	if err != nil {
		return nil, err
	}

	return oidcrp.ClientCredentials(ctx, rp, url.Values{})
}

// Facilitates an auth code PKCE flow to obtain OIDC tokens.
// Spawns a local server to handle the callback and opens a browser window in each respective OS.
func Login(platformEndpoint, tokenURL, authURL, publicClientID string) (*oauth2.Token, error) {
	// Generate random hash and encryption keys for cookie handling
	hashKey := make([]byte, 16)
	encryptKey := make([]byte, 16)

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
		RedirectURL: fmt.Sprintf("http://localhost:%s%s", authCodeFlowPort, authCallbackPath),
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}

	ctx := context.Background()
	cookiehandler := httphelper.NewCookieHandler(hashKey, encryptKey)

	relyingParty, err := oidcrp.NewRelyingPartyOAuth(conf,
		// allow cookie handling for PKCE
		oidcrp.WithCookieHandler(cookiehandler),
		// use PKCE
		oidcrp.WithPKCE(cookiehandler),
		// allow IAT claim offset of 5 seconds
		oidcrp.WithVerifierOpts(oidcrp.WithIssuedAtOffset(5*time.Second)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create relying party: %v", err)
	}
	stateProvider := func() string {
		return uuid.New().String()
	}
	tok := oidcCLI.CodeFlow[*oidc.IDTokenClaims](ctx, relyingParty, authCallbackPath, authCodeFlowPort, stateProvider)
	return tok.Token, nil
}

// Logs in using the auth code PKCE flow driven by the platform well-known idP OIDC configuration.
func LoginWithPKCE(host, publicClientID string, tlsNoVerify bool) (*oauth2.Token, error) {
	pc, err := getPlatformConfiguration(host, tlsNoVerify)
	if err != nil {
		return nil, fmt.Errorf("failed to get platform configuration: %w", err)
	}

	tok, err := Login(host, pc.tokenEndpoint, pc.authzEndpoint, pc.publicClientID)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	return tok, nil
}