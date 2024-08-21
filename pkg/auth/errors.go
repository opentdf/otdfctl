package auth

import "errors"

var ErrAccessTokenExpired = errors.New("access token expired")
var ErrAccessTokenNotFound = errors.New("no access token found")
var ErrClientCredentialsNotFound = errors.New("client credentials not found")
var ErrInvalidAuthType = errors.New("invalid auth type")
var ErrUnauthenticated = errors.New("not logged in")

var ErrProfileCredentialsNotFound = errors.New("profile missing credentials")
