package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zalando/go-keyring"
)

const (
	OTDFCTL_CLIENT_ID_CACHE_KEY = "OTDFCTL_DEFAULT_CLIENT_ID"
	OTDFCTL_OIDC_TOKEN_KEY      = "OTDFCTL_OIDC_TOKEN"
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
	// sdk, err := NewWithCredentials(endpoint, clientID, clientSecret, tlsNoVerify)
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
