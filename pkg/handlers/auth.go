package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const (
	OTDFCTL_CLIENT_ID_CACHE_KEY = "OTDFCTL_DEFAULT_CLIENT_ID"
	OTDFCTL_OIDC_TOKEN_KEY      = "OTDFCTL_OIDC_TOKEN"
)

// TODO: get this dynamically from the platform via SDK or dialing directly: [https://github.com/opentdf/platform/issues/147]
const TOKEN_URL = "http://localhost:8888/auth/realms/opentdf/protocol/openid-connect/token"

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

// GetOIDCTokenFromCache retrieves the OIDC token from the keyring.
func GetOIDCTokenFromCache() (string, error) {
	return keyring.Get(TOKEN_URL, OTDFCTL_OIDC_TOKEN_KEY)
}

// GetClientIDFromCache retrieves the client ID from the keyring.
func GetClientIDFromCache() (string, error) {
	return keyring.Get(TOKEN_URL, OTDFCTL_CLIENT_ID_CACHE_KEY)
}

// GetClientSecretFromCache retrieves the client secret from the keyring.
func GetClientSecretFromCache(clientID string) (string, error) {
	return keyring.Get(TOKEN_URL, clientID)
}

// GetClientSecretFromCache retrieves the client secret from the keyring.
func GetClientIdAndSecretFromCache() (string, string, error) {
	// we use the client id to cache the secret, so retrieve it first
	clientId, err := keyring.Get(TOKEN_URL, OTDFCTL_CLIENT_ID_CACHE_KEY)
	if err != nil || clientId == "" {
		return "", "", ErrUnauthenticated
	}

	clientSecret, err := keyring.Get(TOKEN_URL, clientId)
	if err != nil {
		return "", "", err
	}
	return clientSecret, clientId, nil
}

// GetTokenWithClientCredentials uses the OAuth2 client credentials flow to obtain a token.
func GetTokenWithClientCredentials(ctx context.Context, clientID, clientSecret, tokenURL string, noCache bool) (*oauth2.Token, error) {
	// did the user pass a custom tokenURL?
	if tokenURL == "" {
		// use the default hardcoded constant
		tokenURL = TOKEN_URL
	}

	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     tokenURL,
	}

	token, err := config.Token(ctx)
	if err != nil {
		return nil, err
	}

	// if the users didn't specifically specify not to cache, then we'll cache the clientID, clientSecret, and OIDC_TOKEN in the keyring
	if !noCache {
		// lets store our id and secret in the keyring
		errID := keyring.Set(tokenURL, OTDFCTL_CLIENT_ID_CACHE_KEY, clientID)
		err := keyring.Set(tokenURL, clientID, clientSecret)
		// lets also store the oidc token
		errToken := keyring.Set(tokenURL, OTDFCTL_OIDC_TOKEN_KEY, token.AccessToken)
		if err != nil {
			return nil, err
		}

		if errID != nil {
			return nil, fmt.Errorf("failed to store client ID in keyring: %v", errID)
		}

		if errToken != nil {
			return nil, fmt.Errorf("failed to store OIDC Token in keyring: %v", errToken)
		}
	}
	return token, nil
}
