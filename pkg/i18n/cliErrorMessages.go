package i18n

import "fmt"

type message struct {
	en      string
	locales map[string]string
}

func (m *message) String(args ...interface{}) string {
	s := m.locales[Locale]
	if s != "" {
		return fmt.Sprintf(s, args...)
	}
	return fmt.Sprintf(m.en, args...)
}

func New(enText string, locales map[string]string) *message {
	return &message{
		en:      enText,
		locales: locales,
	}
}

//////////////////////////////////////////
// Errors
//////////////////////////////////////////

// New Handler errors
var (
	MsgUsingProfile                   = New("Using profile [%s]\n", map[string]string{})
	MsgFailedToInitializeProfileStore = New("Failed to initialize profile store", map[string]string{})
	MsgFailedToLoadProfile            = New("Failed to load profile '%s'", map[string]string{})
	MsgNoDefaultProfile               = New("No default profile set. Use `%s profile create <profile> <endpoint>` to create a default profile.", map[string]string{})
	MsgHostMustBeSet                  = New("Host must be set", map[string]string{})
	MsgMixedAuthFlags                 = New("when using global flags %s, profiles will not be used and all required flags must be set", map[string]string{})
	MsgOneAuthFlagMustBeSet           = New("One of %s must be set", map[string]string{})
	MsgOnlyOneAuthFlagMustBeSet       = New("Only one of %s must be set", map[string]string{})
	MsgPlatformConfigNotFound         = New("Failed to get platform configuration. Is the platform accepting connections at '%s'?", map[string]string{})
	MsgFailedToAuthenticate           = New("Failed to authenticate with flag-provided credentials", map[string]string{})
	MsgProfileMissingCreds            = New("Profile missing credentials. Please login or add flag-provided credentials", map[string]string{})
	MsgAccessTokenExpired             = New("Access token expired. Please login again", map[string]string{})
	MsgAccessTokenNotFound            = New("No access token found. Please login or add flag-provided credentials", map[string]string{})
	MsgFailedToGetAccessToken         = New("Failed to get access token", map[string]string{})
	MsgFailedToCreateHandler          = New("Failed to create handler", map[string]string{})

	MsgFailedToInitializeInMemoryProfile = New("Failed to initialize a temporary profile", map[string]string{})
	MsgFailedToCreateInMemoryProfile     = New("Failed to create a temporary profile", map[string]string{})
	MsgFailedToLoadInMemoryProfile       = New("Failed to load temporary profile", map[string]string{})

	MsgFailedToSetAccessToken = New("Failed to set access token", map[string]string{})
	MsgFailedToGetClientCreds = New("Failed to get client credentials", map[string]string{})
	MsgFailedToSetClientCreds = New("Failed to set client credentials", map[string]string{})
	MsgFailedToSaveProfile    = New("Failed to save profile", map[string]string{})
)
