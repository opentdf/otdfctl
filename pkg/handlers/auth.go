package handlers

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

func GetClientCreds(endpoint string, file string, credsJSON []byte) (ClientCredentials, error) {
	if file != "" {
		return GetClientCredsFromFile(file)
	}
	if len(credsJSON) > 0 {
		return GetClientCredsFromJSON(credsJSON)
	}
	return NewKeyring(endpoint).GetClientCredentials()
}

// Uses the OAuth2 client credentials flow to obtain a token.
func GetTokenWithClientCreds(ctx context.Context, endpoint string, c ClientCredentials, tlsNoVerify bool) (*oauth2.Token, error) {
	h, err := New(endpoint, tlsNoVerify)
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
	cookiehandler := httphelper.NewCookieHandler(hashKey, encryptKey, httphelper.WithUnsecure())

	relyingParty, err := oidcrp.NewRelyingPartyOAuth(conf,
		oidcrp.WithCookieHandler(cookiehandler),
		oidcrp.WithPKCE(cookiehandler),
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

func LoginWithPKCE(host string, tlsNoVerify bool, noCache bool) (*oauth2.Token, error) {
	h, err := New(host, tlsNoVerify)
	if err != nil {
		return nil, fmt.Errorf("failed to create handler: %w", err)
	}

	// retrieve idP well-known configuration values
	tokenURL, err := h.Direct().PlatformConfiguration.TokenEndpoint()
	if err != nil || tokenURL == "" {
		return nil, fmt.Errorf("failed to retrieve well-known token endpoint: %w", err)
	}
	authURL, err := h.Direct().PlatformConfiguration.AuthzEndpoint()
	if err != nil || authURL == "" {
		return nil, fmt.Errorf("failed to retrieve well-known authz endpoint: %w", err)
	}
	publicClientID, err := h.Direct().PlatformConfiguration.PublicClientID()
	if err != nil || publicClientID == "" {
		return nil, fmt.Errorf("failed to retrieve well-known public client ID: %w", err)
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
