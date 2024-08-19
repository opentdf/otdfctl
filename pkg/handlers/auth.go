package handlers

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/zalando/go-keyring"
	oidcrp "github.com/zitadel/oidc/v3/pkg/client/rp"
	oidcCLI "github.com/zitadel/oidc/v3/pkg/client/rp/cli"
	httphelper "github.com/zitadel/oidc/v3/pkg/http"
	"github.com/zitadel/oidc/v3/pkg/oidc"
	"golang.org/x/oauth2"
)

const (
	OTDFCTL_CLIENT_ID_CACHE_KEY = "OTDFCTL_DEFAULT_CLIENT_ID"
	OTDFCTL_OIDC_TOKEN_KEY      = "OTDFCTL_OIDC_TOKEN"
	authCallbackPath            = "/callback"
	authCodeFlowPort            = "9000"
)

// CheckTokenExpiration checks if an OIDC token has expired.
// Returns true if the token is still valid, false otherwise.
func CheckTokenExpiration(tokenString string) (bool, error) {
	// for simplicity sake, we're skipping the token validation, and just checking the expiration time, if expired we'll get a new token
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if err != nil {
		return false, err // Token could not be parsed
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			return time.Now().Before(expirationTime), nil // Return true if the current time is before the expiration time
		}
	}

	// Return an error if the expiration time could not be found or parsed
	return false, fmt.Errorf("expiration time (exp) claim is missing or invalid")
}

func ClearCachedCredentials(endpoint string) error {
	cachedClientID, err := GetClientIDFromCache(endpoint)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			fmt.Println("No client-id found in the cache to clear.")
		} else {
			return errors.Join(errors.New("failed to retrieve client id from keyring"), err)
		}
	}

	// clear the client ID and secret from the keyring
	err = keyring.Delete(endpoint, cachedClientID)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			fmt.Println("No client secret found in the cache to clear under client-id: ", cachedClientID)
		} else {
			return errors.Join(errors.New("failed to clear client secret from keyring"), err)
		}
	}

	err = keyring.Delete(endpoint, OTDFCTL_CLIENT_ID_CACHE_KEY)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			fmt.Println("No client id found in the cache to clear.")
		} else {
			return errors.Join(errors.New("failed to clear client id from keyring"), err)
		}
	}

	err = keyring.Delete(endpoint, OTDFCTL_OIDC_TOKEN_KEY)
	if err != nil {
		if errors.Is(err, keyring.ErrNotFound) {
			fmt.Println("No token found in the cache to clear.")
		} else {
			return errors.Join(errors.New("failed to clear token from keyring"), err)
		}
	}

	return nil
}

// GetOIDCTokenFromCache retrieves the OIDC token from the keyring.
func GetOIDCTokenFromCache(endpoint string) (string, error) {
	return keyring.Get(endpoint, OTDFCTL_OIDC_TOKEN_KEY)
}

// GetClientIDFromCache retrieves the client ID from the keyring.
func GetClientIDFromCache(endpoint string) (string, error) {
	return keyring.Get(endpoint, OTDFCTL_CLIENT_ID_CACHE_KEY)
}

// GetClientSecretFromCache retrieves the client secret from the keyring.
func GetClientSecretFromCache(endpoint string, clientID string) (string, error) {
	return keyring.Get(endpoint, clientID)
}

// Client ID and Secret for use in the client credentials flow.
type ClientCreds struct {
	ClientID     string `json:"clientId"`
	ClientSecret string `json:"clientSecret"`
}

// Retrieves credentials by reading specified file
func GetClientCredsFromFile(filepath string) (ClientCreds, error) {
	creds := ClientCreds{}
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
func GetClientCredsFromJSON(credsJSON []byte) (ClientCreds, error) {
	creds := ClientCreds{}
	if err := json.Unmarshal(credsJSON, &creds); err != nil {
		return creds, errors.Join(errors.New("failed to decode creds JSON"), err)
	}

	return creds, nil
}

// Retrieves the client secret from the keyring.
func GetClientCredsFromCache(endpoint string) (ClientCreds, error) {
	creds := ClientCreds{}
	// we use the client id to cache the secret, so retrieve it first
	clientID, err := keyring.Get(endpoint, OTDFCTL_CLIENT_ID_CACHE_KEY)
	if err != nil || clientID == "" {
		return creds, errors.Join(errors.New("could not find clientID in OS keyring"), ErrUnauthenticated)
	}

	clientSecret, err := keyring.Get(endpoint, clientID)
	if err != nil {
		return creds, err
	}
	return ClientCreds{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}, nil
}

func GetClientCreds(endpoint string, file string, credsJSON []byte) (ClientCreds, error) {
	if file != "" {
		return GetClientCredsFromFile(file)
	}
	if len(credsJSON) > 0 {
		return GetClientCredsFromJSON(credsJSON)
	}
	return GetClientCredsFromCache(endpoint)
}

// Uses the OAuth2 client credentials flow to obtain a token.
func GetTokenWithClientCreds(ctx context.Context, endpoint string, clientID string, clientSecret string, tlsNoVerify bool) error {
	// TODO improve the way we validate the client credentials
	// sdk, err := NewWithClientCredentials(endpoint, clientID, clientSecret, tlsNoVerify)
	// if err != nil {
	// 	return err
	// }

	// if _, err := sdk.Direct().Authorization.GetDecisions(ctx, &authorization.GetDecisionsRequest{}); err != nil {
	// 	return errors.Join(errors.New("failed to get token with client credentials"), err)
	// }

	if err := keyring.Set(endpoint, clientID, clientSecret); err != nil {
		return fmt.Errorf("failed to store client secret in key: %v", err)
	}
	if err := keyring.Set(endpoint, OTDFCTL_CLIENT_ID_CACHE_KEY, clientID); err != nil {
		return fmt.Errorf("failed to store client ID in keyring: %v", err)
	}

	return nil
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

	// if a login is initiated, clear any existing token from the keyring proactively
	keyring.Delete(platformEndpoint, OTDFCTL_OIDC_TOKEN_KEY)

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
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve well-known token endpoint: %w", err)
	}
	authURL, err := h.Direct().PlatformConfiguration.AuthzEndpoint()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve well-known authz endpoint: %w", err)
	}
	publicClientID, err := h.Direct().PlatformConfiguration.PublicClientID()
	if err != nil {
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
