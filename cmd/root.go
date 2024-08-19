/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/pkg/handlers/profile"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var (
	cfgKey              string
	OtdfctlCfg          config.Config
	clientCredsFile     string
	clientCredsJSON     string
	configFlagOverrides = config.ConfigFlagOverrides{}

	profileStore *profile.Profile

	RootCmd = &man.Docs.GetDoc("<root>").Command
)

// instantiates a new handler with authentication via client credentials
func NewHandler(cmd *cobra.Command) handlers.Handler {
	flag := cli.NewFlagHelper(cmd)
	host := flag.GetRequiredString("host")
	tlsNoVerify := flag.GetOptionalBool("tls-no-verify")
	clientCredsFile := flag.GetOptionalString("with-client-creds-file")
	clientCredsJSON := flag.GetOptionalString("with-client-creds")
	profileName := flag.GetOptionalString("profile")

	// use the profile
	if profileStore != nil {
		if err := profileStore.UseProfile(profileName); err != nil {
			cli.ExitWithError("Failed to load profile", err)
		}
	}

	var authCredentials profile.AuthCredentials
	if profileStore != nil {
		cp, err := profileStore.GetCurrentProfile()
		if err != nil {
			cli.ExitWithError("Failed to get current profile", err)
		}
		authCredentials = cp.GetAuthCredentials()
	} else {
		creds, err := handlers.GetClientCreds(host, clientCredsFile, []byte(clientCredsJSON))
		if err != nil {
			cli.ExitWithError("Failed to get client credentials", err)
		}

		authCredentials = profile.AuthCredentials{
			AuthType: profile.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS,
			ClientCredentials: profile.ClientCredentials{
				ClientId:     creds.ClientId,
				ClientSecret: creds.ClientSecret,
			},
		}
	}

	h := handlers.Handler{}
	if authCredentials.AuthType == profile.PROFILE_AUTH_TYPE_CLIENT_CREDENTIALS {
		var err error
		clientCredentials := authCredentials.ClientCredentials
		h, err = handlers.NewWithCredentials(host, clientCredentials.ClientId, clientCredentials.ClientSecret, tlsNoVerify)
		if err != nil {
			if errors.Is(err, handlers.ErrUnauthenticated) {
				cli.ExitWithError(fmt.Sprintf("Not logged in. Please authenticate via CLI auth flow(s) before using command (%s %s)", cmd.Parent().Use, cmd.Use), err)
			}
			cli.ExitWithError("Failed to connect to server", err)
		}
	} else {
		cli.ExitWithError("Invalid auth type", errors.New("invalid auth type"))
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
