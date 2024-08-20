package handlers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/google/uuid"
	"github.com/opentdf/otdfctl/pkg/handlers/profile"
	"github.com/zalando/go-keyring"
	oidcrp "github.com/zitadel/oidc/v3/pkg/client/rp"
	oidcCLI "github.com/zitadel/oidc/v3/pkg/client/rp/cli"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
)

const (
	authCallbackPath                   = "/callback"
	authCodeFlowPort                   = "9000"
	OTDFCTL_KEYRING_SERVICE            = "otdfctl"
	OTDFCTL_CLIENT_ID_CACHE_KEY        = "OTDFCTL_DEFAULT_CLIENT_ID"
	OTDFCTL_KEYRING_CLIENT_CREDENTIALS = "OTDFCTL_CLIENT_CREDENTIALS"
	OTDFCTL_OIDC_TOKEN_KEY             = "OTDFCTL_OIDC_TOKEN"
)

type ClientCredentials struct {
	ClientId     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
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

func GetClientCredsFromProfile(p *profile.Profile) (ClientCredentials, error) {
	cp, err := p.GetCurrentProfile()
	if err != nil {
		return ClientCredentials{}, err
	}
	c := cp.GetAuthCredentials()
	if c.AuthType != profile.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS {
		return ClientCredentials{}, errors.New("invalid auth type")
	}

	return ClientCredentials{
		ClientId:     c.ClientCredentials.ClientId,
		ClientSecret: c.ClientCredentials.ClientSecret,
	}, nil
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

func getPlatformIssuer(endpoint string, tlsNoVerify bool) (string, error) {
	// Create a new handler with the provided endpoint and no credentials (empty strings is required by the SDK)
	h, err := NewWithCredentials(endpoint, "", "", tlsNoVerify)
	if err != nil {
		return "", err
	}

	return h.sdk.PlatformIssuer(), nil
}

// Uses the OAuth2 client credentials flow to obtain a token.
func GetTokenWithClientCreds(ctx context.Context, endpoint string, c ClientCredentials, tlsNoVerify bool) (*oauth2.Token, error) {
	issuer, err := getPlatformIssuer(endpoint, tlsNoVerify)
	if err != nil {
		return nil, err
	}

	issuer, err := h.Direct().PlatformConfiguration.Issuer()
	if err != nil {
		return nil, err
	}

	rp, err := oidcrp.NewRelyingPartyOIDC(ctx, issuer, c.ClientId, c.ClientSecret, "", []string{"email"})
	if err != nil {
		return nil, err
	}

	return oidcrp.ClientCredentials(ctx, rp, url.Values{})
}

// Facilitates an auth code PKCE flow to obtain OIDC tokens.
// Spawns a local server to handle the callback and opens a browser window in each respective OS.
func Login(platformEndpoint, tokenURL, authURL, publicClientID string, noPrint bool) (*oauth2.Token, error) {
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
func LoginWithPKCE(host, publicClientID string, tlsNoVerify bool, noCache bool) (*oauth2.Token, error) {
	// retrieve idP well-known configuration values via unauthenticated SDK
	h, err := New(host, tlsNoVerify)
	if err != nil {
		return nil, fmt.Errorf("failed to create handler: %w", err)
	}

	tokenURL, err := h.Direct().PlatformConfiguration.TokenEndpoint()
	if err != nil || tokenURL == "" {
		return nil, fmt.Errorf("failed to retrieve well-known token endpoint: %w", err)
	}
	authURL, err := h.Direct().PlatformConfiguration.AuthzEndpoint()
	if err != nil || authURL == "" {
		return nil, fmt.Errorf("failed to retrieve well-known authz endpoint: %w", err)
	}
	if publicClientID == "" {
		publicClientID, err = h.Direct().PlatformConfiguration.PublicClientID()
		if err != nil || publicClientID == "" {
			return nil, fmt.Errorf("failed to retrieve well-known public client ID: %w", err)
		}
	}

	tok, err := Login(h.platformEndpoint, tokenURL, authURL, publicClientID, noCache)
	if err != nil {
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	if !noCache {
		if err := keyring.Set(h.platformEndpoint, OTDFCTL_OIDC_TOKEN_KEY, tok.AccessToken); err != nil {
			return nil, fmt.Errorf("failed to store token in keyring: %w", err)
		}
	}
	return tok, nil
}
