package profiles

import (
	"context"

	"github.com/opentdf/otdfctl/internal/auth"
	"github.com/opentdf/platform/sdk"
	"golang.org/x/oauth2"
)

type platformConfiguration struct {
	issuer         string
	authzEndpoint  string
	tokenEndpoint  string
	publicClientID string
}

func GetSDKAuthOptionFromProfile(p *ProfileCLI) (sdk.Option, error) {
	c := p.GetAuthCredentials()

	switch c.AuthType {
	case auth.AUTH_TYPE_CLIENT_CREDENTIALS:
		return sdk.WithClientCredentials(c.ClientID, c.ClientSecret, nil), nil
	case auth.AUTH_TYPE_ACCESS_TOKEN:
		tokenSource := oauth2.StaticTokenSource(auth.BuildToken(&c))
		return sdk.WithOAuthAccessTokenSource(tokenSource), nil
	default:
		return nil, ErrInvalidAuthType
	}
}

func ValidateProfileAuthCredentials(ctx context.Context, p *ProfileCLI) error {
	c := p.GetAuthCredentials()
	switch c.AuthType {
	case "":
		return ErrProfileCredentialsNotFound
	case auth.AUTH_TYPE_CLIENT_CREDENTIALS:
		_, err := auth.GetTokenWithClientCreds(ctx, p.GetEndpoint(), c.ClientID, c.ClientSecret, p.GetTLSNoVerify())
		if err != nil {
			return err
		}
		return nil
	case auth.AUTH_TYPE_ACCESS_TOKEN:
		if !auth.BuildToken(&c).Valid() {
			return ErrAccessTokenExpired
		}
	default:
		return ErrInvalidAuthType
	}
	return nil
}

func GetTokenWithProfile(ctx context.Context, p *ProfileCLI) (*oauth2.Token, error) {
	c := p.GetAuthCredentials()
	switch c.AuthType {
	case auth.AUTH_TYPE_CLIENT_CREDENTIALS:
		return auth.GetTokenWithClientCreds(ctx, p.GetEndpoint(), c.ClientID, c.ClientSecret, p.GetTLSNoVerify())
	case auth.AUTH_TYPE_ACCESS_TOKEN:
		return auth.BuildToken(&c), nil
	default:
		return nil, ErrInvalidAuthType
	}
}
