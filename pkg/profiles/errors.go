package profiles

import "errors"

var (
	ErrDeletingDefaultProfile = errors.New("cannot delete the default profile")

	ErrAccessTokenExpired        = errors.New("access token expired")
	ErrAccessTokenNotFound       = errors.New("no access token found")
	ErrClientCredentialsNotFound = errors.New("client credentials not found")
	ErrInvalidAuthType           = errors.New("invalid auth type")
	ErrUnauthenticated           = errors.New("not logged in")

	ErrProfileCredentialsNotFound = errors.New("profile missing credentials")

	ErrStoredProfileInvalid = errors.New("stored profile is invalid")

	// internal error
	errStoredProfileWrongType = errors.New("stored profile is not of type ProfileCLI")
)
