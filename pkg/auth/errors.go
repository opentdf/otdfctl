package auth

import "errors"

var (
	ErrAccessTokenExpired        = errors.New("access token expired")
	ErrAccessTokenNotFound       = errors.New("no access token found")
	ErrClientCredentialsNotFound = errors.New("client credentials not found")
	ErrInvalidAuthType           = errors.New("invalid auth type")
	ErrUnauthenticated           = errors.New("not logged in")
)

var (
	ErrProfileCredentialsNotFound = errors.New("profile missing credentials")
	ErrPlatformConfigNotFound     = errors.New("platform configuration not found")
)
