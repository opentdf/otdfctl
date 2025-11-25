package common

import (
	"errors"
	"fmt"

	"github.com/evertras/bubble-table/table"
	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/opentdf/platform/sdk"
	"github.com/spf13/cobra"
)

// captures all CLI flags that will override pre-specified config values
type ConfigFlagOverrides struct {
	OutputFormatJSON bool
}

var (
	Profile *profiles.Profile

	configFlagOverrides = ConfigFlagOverrides{}

	OtdfctlCfg config.Config
)

// InitProfile initializes the profile store and loads the profile specified in the flags
// if onlyNew is set to true, a new profile will be created and returned
// returns the profile and the current profile store
func InitProfile(c *cli.Cli, onlyNew bool) (*profiles.Profile, *profiles.ProfileStore) {
	var err error
	profileName := c.FlagHelper.GetOptionalString("profile")

	Profile, err = profiles.New()
	if err != nil || Profile == nil {
		c.ExitWithError("Failed to initialize profile store", err)
	}

	// short circuit if onlyNew is set to enable creating a new profile
	if onlyNew && profileName == "" {
		return Profile, nil
	}

	// check if there exists a default profile and warn if not with steps to create one
	if Profile.GetGlobalConfig().GetDefaultProfile() == "" {
		c.ExitWithWarning(fmt.Sprintf("No default profile set. Use `%s profile create <profile> <endpoint>` to create a default profile.", config.AppName))
	}

	if profileName == "" {
		profileName = Profile.GetGlobalConfig().GetDefaultProfile()
	}

	c.Printf("Using profile [%s]\n", profileName)

	// load profile
	cp, err := Profile.UseProfile(profileName)
	if err != nil {
		c.ExitWithError(fmt.Sprintf("Failed to load profile: %s", profileName), err)
	}

	return Profile, cp
}

// instantiates a new handler with authentication via client credentials
// TODO make this a preRun hook
//
//nolint:nestif // separate refactor [https://github.com/opentdf/otdfctl/issues/383]
func NewHandler(c *cli.Cli) handlers.Handler {
	// if global flags are set then validate and create a temporary profile in memory
	var cp *profiles.ProfileStore

	// Non-profile flags
	host := c.FlagHelper.GetOptionalString("host")
	tlsNoVerify := c.FlagHelper.GetOptionalBool("tls-no-verify")
	withClientCreds := c.FlagHelper.GetOptionalString("with-client-creds")
	withClientCredsFile := c.FlagHelper.GetOptionalString("with-client-creds-file")
	withAccessToken := c.FlagHelper.GetOptionalString("with-access-token")
	var inMemoryProfile bool

	authFlags := []string{"--with-access-token", "--with-client-creds", "--with-client-creds-file"}
	nonProfileFlags := append([]string{"--host", "--tls-no-verify"}, authFlags...)
	hasNonProfileFlags := host != "" || tlsNoVerify || withClientCreds != "" || withClientCredsFile != "" || withAccessToken != ""

	//nolint:nestif // nested if statements are necessary for validation
	if hasNonProfileFlags {
		err := fmt.Errorf("when using global flags %s, profiles will not be used and all required flags must be set", cli.PrettyList(nonProfileFlags))

		// host must be set
		if host == "" {
			cli.ExitWithError("Host must be set", err)
		}

		authFlagsCounter := 0
		if withAccessToken != "" {
			authFlagsCounter++
		}
		if withClientCreds != "" {
			authFlagsCounter++
		}
		if withClientCredsFile != "" {
			authFlagsCounter++
		}
		if authFlagsCounter == 0 {
			cli.ExitWithError(fmt.Sprintf("One of %s must be set", cli.PrettyList(authFlags)), err)
		} else if authFlagsCounter > 1 {
			cli.ExitWithError(fmt.Sprintf("Only one of %s must be set", cli.PrettyList(authFlags)), err)
		}

		inMemoryProfile = true
		Profile, err = profiles.New(profiles.WithInMemoryStore())
		if err != nil || Profile == nil {
			cli.ExitWithError("Failed to initialize in-memory profile", err)
		}

		if err := Profile.AddProfile("temp", host, tlsNoVerify, true); err != nil {
			cli.ExitWithError("Failed to create in-memory profile", err)
		}

		// add credentials to the temporary profile
		cp, err = Profile.UseProfile("temp")
		if err != nil {
			cli.ExitWithError("Failed to load in-memory profile", err)
		}

		// get credentials from flags
		if withAccessToken != "" {
			claims, err := auth.ParseClaimsJWT(withAccessToken)
			if err != nil {
				cli.ExitWithError("Failed to get access token", err)
			}

			if err := cp.SetAuthCredentials(profiles.AuthCredentials{
				AuthType: profiles.AuthTypeAccessToken,
				AccessToken: profiles.AuthCredentialsAccessToken{
					AccessToken: withAccessToken,
					Expiration:  claims.Expiration,
				},
			}); err != nil {
				cli.ExitWithError("Failed to set access token", err)
			}
		} else {
			var cc auth.ClientCredentials
			if withClientCreds != "" {
				cc, err = auth.GetClientCredsFromJSON([]byte(withClientCreds))
			} else if withClientCredsFile != "" {
				cc, err = auth.GetClientCredsFromFile(withClientCredsFile)
			}
			if err != nil {
				cli.ExitWithError("Failed to get client credentials", err)
			}

			// add credentials to the temporary profile
			if err := cp.SetAuthCredentials(profiles.AuthCredentials{
				AuthType:     profiles.AuthTypeClientCredentials,
				ClientID:     cc.ClientID,
				ClientSecret: cc.ClientSecret,
			}); err != nil {
				cli.ExitWithError("Failed to set client credentials", err)
			}
		}
		if err := cp.Save(); err != nil {
			cli.ExitWithError("Failed to save profile", err)
		}
	} else {
		Profile, cp = InitProfile(c, false)
	}

	if err := auth.ValidateProfileAuthCredentials(c.Context(), cp); err != nil {
		if errors.Is(err, sdk.ErrPlatformUnreachable) {
			cli.ExitWithError(fmt.Sprintf("Failed to connect to the platform. Is the platform accepting connections at '%s'?", cp.GetEndpoint()), nil)
		}
		if errors.Is(err, sdk.ErrPlatformConfigFailed) {
			cli.ExitWithError(fmt.Sprintf("Failed to get the platform configuration. Is the platform serving a well-known configuration at '%s'?", cp.GetEndpoint()), nil)
		}
		if inMemoryProfile {
			cli.ExitWithError("Failed to authenticate with flag-provided client credentials.", err)
		}
		if errors.Is(err, auth.ErrProfileCredentialsNotFound) {
			cli.ExitWithWarning("Profile missing credentials. Please login or add client credentials.")
		}

		if errors.Is(err, auth.ErrAccessTokenExpired) {
			cli.ExitWithWarning("Access token expired. Please login or add flag-provided credentials.")
		}
		if errors.Is(err, auth.ErrAccessTokenNotFound) {
			cli.ExitWithWarning("No access token found. Please login or add flag-provided credentials.")
		}
		cli.ExitWithError("Failed to get access token.", err)
	}

	h, err := handlers.New(handlers.WithProfile(cp))
	if err != nil {
		cli.ExitWithError("Unexpected error", err)
	}
	return h
}

// HandleSuccess prints a success message according to the configured format (styled table or JSON)
func HandleSuccess(command *cobra.Command, id string, t table.Model, policyObject interface{}) {
	c := cli.New(command, []string{})
	jsonFlag := c.Flags.GetOptionalBool("json")
	if OtdfctlCfg.Output.Format == config.OutputJSON || configFlagOverrides.OutputFormatJSON || jsonFlag {
		c.ExitWithJSON(policyObject)
	}
	cli.PrintSuccessTable(command, id, t)
}
