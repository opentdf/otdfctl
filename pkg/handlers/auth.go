package handlers

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zalando/go-keyring"
	"golang.org/x/oauth2"
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

type AuthorizationCodePKCE struct {
	Oauth2Config *oauth2.Config
	Token        *oauth2.Token
}

type OpenTdfTokenSource struct {
	OpenTdfToken *oauth2.Token
}

const (
	opentdfPublicClientID = "opentdf-public"
	authCodeFlowPort      = "9000"
)

func (acp *AuthorizationCodePKCE) Login(platformEndpoint, tokenURL, authURL string, noPrint bool) (*oauth2.Token, error) {
	var (
		token *oauth2.Token
		err   error
	)

	// if a login is initiated, clear any existing token from the keyring proactively
	keyring.Delete(platformEndpoint, OTDFCTL_OIDC_TOKEN_KEY)

	conf := &oauth2.Config{
		ClientID:    opentdfPublicClientID,
		Scopes:      []string{"openid", "profile", "email"},
		RedirectURL: fmt.Sprintf("http://localhost:%s/callback", authCodeFlowPort),
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
	acp.Oauth2Config = conf

	// Create a HTTP server to handle the callback ":9000"
	srv := &http.Server{Addr: ":9000"}
	stop := make(chan os.Signal, 1)

	// Generate a code verifier and code challenge.
	verifier, err := generateCodeVerifier()
	if err != nil {
		return nil, fmt.Errorf("failed to generate code verifier: %v", err)
	}
	challenge := generateCodeChallenge(verifier)

	// Start a web server to handle the OAuth2 callback.
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		// Get the authorization code from the query parameters.
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing authorization code", http.StatusBadRequest)
			return
		}

		// Exchange the authorization code for an access token.
		token, err = conf.Exchange(context.Background(), code, oauth2.SetAuthURLParam("code_verifier", verifier))
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to exchange authorization code: %v", err), http.StatusInternalServerError)
			return
		}

		// Let the user know the flow was successful.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode("Return to the CLI to continue. You may close this tab.")

		// Send a value to the stop channel to simulate the SIGINT signal.
		stop <- syscall.SIGINT
	})
	url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("code_challenge", challenge), oauth2.SetAuthURLParam("code_challenge_method", "S256"), oauth2.SetAuthURLParam("audience", "http://localhost:8080"))

	// avoid printing the help directions if not caching the token to avoid breaking scripts
	if !noPrint {
		fmt.Print("Open the following URL in a browser if it did not automatically open for you: ", url)
	}
	openBrowser(url)

	// Start the HTTP server in a separate goroutine.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(fmt.Errorf("failed to start HTTP server: %w", err))
		}
	}()

	// Wait for a SIGINT or SIGTERM signal to shutdown the server.
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Failed to shutdown HTTP server gracefully: %v", err)
		return nil, err
	}
	acp.Token = token
	return token, nil
}

func (acp *AuthorizationCodePKCE) Client() (*http.Client, error) {
	token, err := acp.Oauth2Config.TokenSource(context.Background(), acp.Token).Token()
	if err != nil {
		return nil, err
	}
	return acp.Oauth2Config.Client(context.Background(), token), nil
}

func openBrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		return fmt.Errorf("failed to open browser: %v", err)
	}

	return nil
}

func generateCodeVerifier() (string, error) {
	const codeVerifierLength = 32 // You can adjust the length of the code verifier as needed
	randomBytes := make([]byte, codeVerifierLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", fmt.Errorf("failed to generate code verifier: %v", err)
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func generateCodeChallenge(verifier string) string {
	hash := sha256.Sum256([]byte(verifier))
	return base64.RawURLEncoding.EncodeToString(hash[:])
}

func (ots *OpenTdfTokenSource) Token() (*oauth2.Token, error) {
	return ots.OpenTdfToken, nil
}

func LoginWithPKCE(host string, tlsNoVerify bool, noCache bool) (*oauth2.Token, error) {
	h, err := New(host, tlsNoVerify)
	if err != nil {
		return nil, fmt.Errorf("failed to create handler: %w", err)
	}
	tokenURL, err := h.Direct().PlatformTokenEndpoint()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve well-known token endpoint: %w", err)
	}
	authURL, err := h.Direct().PlatformAuthzEndpoint()
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve well-known authz endpoint: %w", err)
	}

	acp := new(AuthorizationCodePKCE)

	tok, err := acp.Login(h.platformEndpoint, tokenURL, authURL, noCache)
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

func buildTokenSource(token string) oauth2.TokenSource {
	return &OpenTdfTokenSource{
		OpenTdfToken: &oauth2.Token{
			AccessToken: token,
		},
	}
}