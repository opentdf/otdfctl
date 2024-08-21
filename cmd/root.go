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

func InitProfile(cmd *cobra.Command) *profiles.ProfileStore {
	flag := cli.NewFlagHelper(cmd)
	profileName := flag.GetOptionalString("profile")

	if profile == nil {
		cli.ExitWithError("Profile not loaded", nil)
	}

	if profileName == "" {
		profileName = profile.GetGlobalConfig().GetDefaultProfile()
	}
	if err := profile.UseProfile(profileName); err != nil {
		cli.ExitWithError("Failed to load profile "+profileName, err)
	}

	cp, err := profile.GetCurrentProfile()
	if err != nil {
		cli.ExitWithError("Failed to get profile "+profileName, err)
	}

	return cp
}

// instantiates a new handler with authentication via client credentials
// TODO make this a preRun hook
func NewHandler(cmd *cobra.Command) handlers.Handler {
	// TODO add support for without profile

	cp := InitProfile(cmd)

	if err := auth.ValidateProfileAuthCredentials(cmd.Context(), cp); err != nil {
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
		flaghelper := cli.NewFlagHelper(cmd)

		if flaghelper.GetOptionalBool("version") {
			fmt.Println(config.AppName + " version " + config.Version + " (" + config.BuildTime + ") " + config.CommitSha)
			return
		}

		cmd.Help()
	}))

	RootCmd = &rootCmd.Command

	RootCmd.Flags().Bool(
		rootCmd.GetDocFlag("version").Name,
		rootCmd.GetDocFlag("version").DefaultAsBool(),
		rootCmd.GetDocFlag("version").Description,
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
	RootCmd.AddGroup(&cobra.Group{ID: "tdf"})
}
