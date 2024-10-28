/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"

	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/i18n"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/otdfctl/pkg/profiles"
	"github.com/spf13/cobra"
)

var (
	cfgKey              string
	OtdfctlCfg          config.Config
	clientCredsFile     string
	clientCredsJSON     string
	configFlagOverrides = config.ConfigFlagOverrides{}

	profile *profiles.Profile

	RootCmd = &man.Docs.GetDoc("<root>").Command
)

type version struct {
	AppName   string `json:"app_name"`
	Version   string `json:"version"`
	CommitSha string `json:"commit_sha"`
	BuildTime string `json:"build_time"`
}

// InitProfile initializes the profile store and loads the profile specified in the flags
// if onlyNew is set to true, a new profile will be created and returned
// returns the profile and the current profile store
func InitProfile(c *cli.Cli, onlyNew bool) (*profiles.Profile, *profiles.ProfileStore) {
	var err error
	profileName := c.FlagHelper.GetOptionalString("profile")

	profile, err = profiles.New()
	if err != nil || profile == nil {
		c.ExitWithError(i18n.MsgFailedToInitializeProfileStore.String(), err)
	}

	// short circuit if onlyNew is set to enable creating a new profile
	if onlyNew && profileName == "" {
		return profile, nil
	}

	// check if there exists a default profile and warn if not with steps to create one
	if profile.GetGlobalConfig().GetDefaultProfile() == "" {
		c.ExitWithWarning(i18n.MsgNoDefaultProfile.String(config.AppName))
	}

	if profileName == "" {
		profileName = profile.GetGlobalConfig().GetDefaultProfile()
	}

	c.Printf(i18n.MsgUsingProfile.String(profileName))

	// load profile
	cp, err := profile.UseProfile(profileName)
	if err != nil {
		c.ExitWithError(i18n.MsgFailedToLoadProfile.String(profileName), err)
	}

	return profile, cp
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
		err := errors.New(i18n.MsgMixedAuthFlags.String(cli.PrettyList(nonProfileFlags)))

		// host must be set
		if host == "" {
			cli.ExitWithError(i18n.MsgHostMustBeSet.String(), err)
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
			cli.ExitWithError(i18n.MsgOneAuthFlagMustBeSet.String(cli.PrettyList(authFlags)), err)
		} else if authFlagsCounter > 1 {
			cli.ExitWithError(i18n.MsgOnlyOneAuthFlagMustBeSet.String(cli.PrettyList(authFlags)), err)
		}

		inMemoryProfile = true
		profile, err = profiles.New(profiles.WithInMemoryStore())
		if err != nil || profile == nil {
			cli.ExitWithError(i18n.MsgFailedToInitializeInMemoryProfile.String(), err)
		}

		if err := profile.AddProfile("temp", host, tlsNoVerify, true); err != nil {
			cli.ExitWithError(i18n.MsgFailedToCreateInMemoryProfile.String(), err)
		}

		// add credentials to the temporary profile
		cp, err = profile.UseProfile("temp")
		if err != nil {
			cli.ExitWithError(i18n.MsgFailedToLoadInMemoryProfile.String(), err)
		}

		// get credentials from flags
		if withAccessToken != "" {
			claims, err := auth.ParseClaimsJWT(withAccessToken)
			if err != nil {
				cli.ExitWithError(i18n.MsgFailedToGetAccessToken.String(), err)
			}

			if err := cp.SetAuthCredentials(profiles.AuthCredentials{
				AuthType: profiles.PROFILE_AUTH_TYPE_ACCESS_TOKEN,
				AccessToken: profiles.AuthCredentialsAccessToken{
					AccessToken: withAccessToken,
					Expiration:  claims.Expiration,
				},
			}); err != nil {
				cli.ExitWithError(i18n.MsgFailedToSetAccessToken.String(), err)
			}
		} else {
			var cc auth.ClientCredentials
			if withClientCreds != "" {
				cc, err = auth.GetClientCredsFromJSON([]byte(withClientCreds))
			} else if withClientCredsFile != "" {
				cc, err = auth.GetClientCredsFromFile(withClientCredsFile)
			}
			if err != nil {
				cli.ExitWithError(i18n.MsgFailedToGetAccessToken.String(), err)
			}

			// add credentials to the temporary profile
			if err := cp.SetAuthCredentials(profiles.AuthCredentials{
				AuthType:     profiles.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS,
				ClientId:     cc.ClientId,
				ClientSecret: cc.ClientSecret,
			}); err != nil {
				cli.ExitWithError(i18n.MsgFailedToSetClientCreds.String(), err)
			}
		}
		if err := cp.Save(); err != nil {
			cli.ExitWithError(i18n.MsgFailedToSaveProfile.String(), err)
		}
	} else {
		profile, cp = InitProfile(c, false)
	}

	if err := auth.ValidateProfileAuthCredentials(c.Context(), cp); err != nil {
		if errors.Is(err, auth.ErrPlatformConfigNotFound) {
			cli.ExitWithError(i18n.MsgPlatformConfigNotFound.String(cp.GetEndpoint()), nil)
		}
		if inMemoryProfile {
			cli.ExitWithError(i18n.MsgFailedToAuthenticate.String(), err)
		}
		if errors.Is(err, auth.ErrProfileCredentialsNotFound) {
			cli.ExitWithWarning(i18n.MsgProfileMissingCreds.String())
		}

		if errors.Is(err, auth.ErrAccessTokenExpired) {
			cli.ExitWithWarning(i18n.MsgAccessTokenExpired.String())
		}
		if errors.Is(err, auth.ErrAccessTokenNotFound) {
			cli.ExitWithWarning(i18n.MsgAccessTokenNotFound.String())
		}
		cli.ExitWithError(i18n.MsgFailedToGetAccessToken.String(), err)
	}

	h, err := handlers.New(handlers.WithProfile(cp))
	if err != nil {
		cli.ExitWithError(i18n.MsgFailedToCreateHandler.String(), err)
	}
	return h
}

func init() {
	rootCmd := man.Docs.GetCommand("<root>", man.WithRun(func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)

		if c.Flags.GetOptionalBool("version") {
			v := version{
				AppName:   config.AppName,
				Version:   config.Version,
				CommitSha: config.CommitSha,
				BuildTime: config.BuildTime,
			}

			c.Println(config.AppName + " version " + config.Version + " (" + config.BuildTime + ") " + config.CommitSha)
			c.ExitWithJSON(v)
			return
		}

		//nolint:errcheck // error does not need to be checked
		cmd.Help()
	}))

	RootCmd = &rootCmd.Command

	RootCmd.Flags().Bool(
		rootCmd.GetDocFlag("version").Name,
		rootCmd.GetDocFlag("version").DefaultAsBool(),
		rootCmd.GetDocFlag("version").Description,
	)

	RootCmd.PersistentFlags().Bool(
		rootCmd.GetDocFlag("debug").Name,
		rootCmd.GetDocFlag("debug").DefaultAsBool(),
		rootCmd.GetDocFlag("debug").Description,
	)

	RootCmd.PersistentFlags().Bool(
		rootCmd.GetDocFlag("json").Name,
		rootCmd.GetDocFlag("json").DefaultAsBool(),
		rootCmd.GetDocFlag("json").Description,
	)

	RootCmd.PersistentFlags().String(
		rootCmd.GetDocFlag("profile").Name,
		rootCmd.GetDocFlag("profile").Default,
		rootCmd.GetDocFlag("profile").Description,
	)

	RootCmd.PersistentFlags().String(
		rootCmd.GetDocFlag("host").Name,
		rootCmd.GetDocFlag("host").Default,
		rootCmd.GetDocFlag("host").Description,
	)
	RootCmd.PersistentFlags().Bool(
		rootCmd.GetDocFlag("tls-no-verify").Name,
		rootCmd.GetDocFlag("tls-no-verify").DefaultAsBool(),
		rootCmd.GetDocFlag("tls-no-verify").Description,
	)
	RootCmd.PersistentFlags().String(
		rootCmd.GetDocFlag("log-level").Name,
		rootCmd.GetDocFlag("log-level").Default,
		rootCmd.GetDocFlag("log-level").Description,
	)
	RootCmd.PersistentFlags().StringVar(
		&clientCredsFile,
		rootCmd.GetDocFlag("with-client-creds-file").Name,
		rootCmd.GetDocFlag("with-client-creds-file").Default,
		rootCmd.GetDocFlag("with-client-creds-file").Description,
	)
	RootCmd.PersistentFlags().StringVar(
		&clientCredsJSON,
		rootCmd.GetDocFlag("with-client-creds").Name,
		rootCmd.GetDocFlag("with-client-creds").Default,
		rootCmd.GetDocFlag("with-client-creds").Description,
	)
	RootCmd.PersistentFlags().String(
		rootCmd.GetDocFlag("with-access-token").Name,
		rootCmd.GetDocFlag("with-access-token").Default,
		rootCmd.GetDocFlag("with-access-token").Description,
	)
	RootCmd.AddGroup(&cobra.Group{ID: TDF})
}
