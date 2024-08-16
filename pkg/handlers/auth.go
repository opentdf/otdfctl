package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"os"

	oidcrp "github.com/zitadel/oidc/v3/pkg/client/rp"
	"golang.org/x/oauth2"
)

const (
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

	rp, err := oidcrp.NewRelyingPartyOIDC(ctx, issuer, c.ClientId, c.ClientSecret, "", []string{"email"})
	if err != nil {
		return nil, err
	}

	return oidcrp.ClientCredentials(ctx, rp, url.Values{})
}
