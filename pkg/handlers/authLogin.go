package handlers

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

const TRUCTL_CLIENT_ID_CACHE_KEY = "TRUCTL_DEFAULT_CLIENT_ID"
const TRUCTL_OIDC_TOKEN_KEY = "TRUCTL_OIDC_TOKEN"

// we're hardcoding this for now, but eventually it will be retrieved from the backend config
// TODO udpate to use the wellknown endpoint for the platform (https://github.com/opentdf/platform/pull/296)
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
	token, err := keyring.Get(TOKEN_URL, TRUCTL_OIDC_TOKEN_KEY)
	if err != nil {
		return "", err
	}
	return token, nil
}

// GetClientIDFromCache retrieves the client ID from the keyring.
func GetClientIDFromCache() (string, error) {
	clientId, err := keyring.Get(TOKEN_URL, TRUCTL_CLIENT_ID_CACHE_KEY)
	if err != nil {
		return "", err
	}
	return clientId, nil
}

// GetClientSecretFromCache retrieves the client secret from the keyring.
func GetClientIdAndSecretFromCache() (string, string, error) {
	// our clientSecret key, is our clientId, so we gotta grab that first
	clientId, err := keyring.Get(TOKEN_URL, TRUCTL_CLIENT_ID_CACHE_KEY)
	if err != nil {
		// we failed to get the clientId for somereason
		return "", "", err
	}

	if clientId == "" {
		return "", "", fmt.Errorf("no clientId found in keyring")
	}

	clientSecret, err := keyring.Get(TOKEN_URL, clientId)
	if err != nil {
		return "", "", err
	}
	return clientSecret, clientId, nil
}

// DEBUG_PrintKeyRingSecrets prints all the secrets in the keyring.
func (h *Handler) DEBUG_PrintKeyRingSecrets() {

	clientId, err := keyring.Get(TOKEN_URL, TRUCTL_CLIENT_ID_CACHE_KEY)
	if err != nil {
		fmt.Println("Failed to retrieve secret from keyring:", err)
		return
	}

	// and our special clientId key, to grab the secret
	secret, errSec := keyring.Get(TOKEN_URL, clientId)
	OIDC_TOKEN, errToken := keyring.Get(TOKEN_URL, TRUCTL_OIDC_TOKEN_KEY)

	if errSec != nil {
		fmt.Println("Failed to retrieve secret from keyring:", err)
		return
	}

	if errToken != nil {
		fmt.Println("Failed to retrieve secret from keyring:", errToken)
		return
	}

	fmt.Println(clientId, ":", secret)
	fmt.Println("Stored OIDC_TOKEN OF:", OIDC_TOKEN)
}

// GetTokenWithClientCredentials uses the OAuth2 client credentials flow to obtain a token.
func (h *Handler) GetTokenWithClientCredentials(clientID, clientSecret, tokenURL string, noCache bool) (*oauth2.Token, error) {
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

	token, err := config.Token(h.ctx)
	if err != nil {
		return nil, err
	}

	// if the users didn't specifically specify not to cache, then we'll cache the clientID, clientSecret, and OIDC_TOKEN in the keyring
	if !noCache {
		// lets store our id and secret in the keyring
		errID := keyring.Set(tokenURL, TRUCTL_CLIENT_ID_CACHE_KEY, clientID)
		err := keyring.Set(tokenURL, clientID, clientSecret)
		// lets also store the oidc token
		errToken := keyring.Set(tokenURL, TRUCTL_OIDC_TOKEN_KEY, token.AccessToken)
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
	h.OIDC_TOKEN = token.AccessToken
	return token, nil
}

// GetTokenWithPasswordFlow creates a custom request to obtain a token using the resource owner password credentials flow.
func (h *Handler) GetTokenWithPasswordFlow(username, password, clientID, clientSecret, tokenURL string, noCache bool) (string, error) {
	errMsg := "Method `GetTokenWithPasswordFlow` is not yet implemented. Please reach out to a Virtru Platform team member to inquire about the status of it."
	fmt.Println(errMsg)
	return "", nil
}
