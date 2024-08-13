/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/internal/config"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/spf13/cobra"
)

var (
	cfgKey          string
	OtdfctlCfg      config.Config
	clientCredsFile string
	clientCredsJSON string

	configFlagOverrides = config.ConfigFlagOverrides{}
)

// RootCmd represents the base command when called without any subcommands.
var (
	RootCmd = &man.Docs.GetDoc("<root>").Command
)

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
