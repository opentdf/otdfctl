/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/opentdf/otdfctl/pkg/auth"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/handlers"
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

func InitProfile(c *cli.Cli, onlyNew bool) *profiles.ProfileStore {
	var err error
	profileName := c.FlagHelper.GetOptionalString("profile")

	profile, err = profiles.New()
	if err != nil || profile == nil {
		c.ExitWithError("Failed to initialize profile store", err)
	}

	// short circuit if onlyNew is set to enable creating a new profile
	if onlyNew {
		return nil
	}

	// check if there exists a default profile and warn if not with steps to create one
	if profile.GetGlobalConfig().GetDefaultProfile() == "" {
		c.ExitWithWarning("No default profile set. Use `" + config.AppName + " profile create <profile> <endpoint>` to create a default profile.")
	}
	c.Printf("Using profile [%s]\n", profile.GetGlobalConfig().GetDefaultProfile())

	if profileName == "" {
		profileName = profile.GetGlobalConfig().GetDefaultProfile()
	}

	// load profile
	cp, err := profile.UseProfile(profileName)
	if err != nil {
		c.ExitWithError("Failed to load profile "+profileName, err)
	}

	return cp
}

// instantiates a new handler with authentication via client credentials
// TODO make this a preRun hook
//
//nolint:nestif // separate refactor [https://github.com/opentdf/otdfctl/issues/383]
func NewHandler(c *cli.Cli) handlers.Handler {
	// Non-profile flags
	host := c.FlagHelper.GetOptionalString("host")
	tlsNoVerify := c.FlagHelper.GetOptionalBool("tls-no-verify")
	withClientCreds := c.FlagHelper.GetOptionalString("with-client-creds")
	withClientCredsFile := c.FlagHelper.GetOptionalString("with-client-creds-file")
	var inMemoryProfile bool

	// if global flags are set then validate and create a temporary profile in memory
	var cp *profiles.ProfileStore
	//nolint:nestif // nested if statements are necessary for validation
	if host != "" || tlsNoVerify || withClientCreds != "" || withClientCredsFile != "" {
		err := errors.New(
			"when using global flags --host, --tls-no-verify, --with-client-creds, or --with-client-creds-file, " +
				"profiles will not be used and all required flags must be set",
		)

		// host must be set
		if host == "" {
			cli.ExitWithError("Host must be set", err)
		}

		// either with-client-creds or with-client-creds-file must be set
		if withClientCreds == "" && withClientCredsFile == "" {
			cli.ExitWithError("Either --with-client-creds or --with-client-creds-file must be set", err)
		} else if withClientCreds != "" && withClientCredsFile != "" {
			cli.ExitWithError("Only one of --with-client-creds or --with-client-creds-file can be set", err)
		}

		var cc auth.ClientCredentials
		if withClientCreds != "" {
			cc, err = auth.GetClientCredsFromJSON([]byte(withClientCreds))
		} else {
			cc, err = auth.GetClientCredsFromFile(withClientCredsFile)
		}
		if err != nil {
			cli.ExitWithError("Failed to get client credentials", err)
		}

		inMemoryProfile = true
		profile, err = profiles.New(profiles.WithInMemoryStore())
		if err != nil || profile == nil {
			cli.ExitWithError("Failed to initialize a temporary profile", err)
		}

		if err := profile.AddProfile("temp", host, tlsNoVerify, true); err != nil {
			cli.ExitWithError("Failed to create temporary profile", err)
		}

		// add credentials to the temporary profile
		cp, err = profile.UseProfile("temp")
		if err != nil {
			cli.ExitWithError("Failed to load temporary profile", err)
		}

		// add credentials to the temporary profile
		if err := cp.SetAuthCredentials(profiles.AuthCredentials{
			AuthType:     profiles.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS,
			ClientId:     cc.ClientId,
			ClientSecret: cc.ClientSecret,
		}); err != nil {
			cli.ExitWithError("Failed to set client credentials", err)
		}
		if err := cp.Save(); err != nil {
			cli.ExitWithError("Failed to save profile", err)
		}
	} else {
		cp = InitProfile(c, false)
	}

	if err := auth.ValidateProfileAuthCredentials(c.Context(), cp); err != nil {
		if errors.Is(err, auth.ErrPlatformConfigNotFound) {
			cli.ExitWithError(fmt.Sprintf("Failed to get platform configuration. Is the platform accepting connections at '%s'?", cp.GetEndpoint()), nil)
		}
		if inMemoryProfile {
			cli.ExitWithError("Failed to authenticate with flag-provided client credentials", err)
		}
		if errors.Is(err, auth.ErrProfileCredentialsNotFound) {
			cli.ExitWithWarning("Profile missing credentials. Please login or add client credentials.")
		}

		if errors.Is(err, auth.ErrAccessTokenExpired) {
			cli.ExitWithWarning("Access token expired. Please login again.")
		}
		if errors.Is(err, auth.ErrAccessTokenNotFound) {
			cli.ExitWithWarning("No access token found. Please login or add client credentials.")
		}
		cli.ExitWithError("Failed to get access token", err)
	}

	h, err := handlers.New(handlers.WithProfile(cp))
	if err != nil {
		cli.ExitWithError("Failed to create handler", err)
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
	RootCmd.AddGroup(&cobra.Group{ID: TDF})
}
