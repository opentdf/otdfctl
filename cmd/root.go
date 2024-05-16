/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/opentdf/otdfctl/internal/config"
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
	doc := man.Docs.GetDoc("<root>")
	RootCmd = &doc.Command
	RootCmd.PersistentFlags().String(
		doc.GetDocFlag("config-file").Name,
		doc.GetDocFlag("config-file").Default,
		doc.GetDocFlag("config-file").Description,
	)
	RootCmd.PersistentFlags().String(
		doc.GetDocFlag("config-key").Name,
		doc.GetDocFlag("config-key").Default,
		doc.GetDocFlag("config-key").Description,
	)
	RootCmd.PersistentFlags().String(
		doc.GetDocFlag("host").Name,
		doc.GetDocFlag("host").Default,
		doc.GetDocFlag("host").Description,
	)
	RootCmd.PersistentFlags().Bool(
		doc.GetDocFlag("insecure").Name,
		doc.GetDocFlag("insecure").DefaultAsBool(),
		doc.GetDocFlag("insecure").Description,
	)
	RootCmd.PersistentFlags().String(
		doc.GetDocFlag("log-level").Name,
		doc.GetDocFlag("log-level").Default,
		doc.GetDocFlag("log-level").Description,
	)
	RootCmd.PersistentFlags().StringVar(
		&clientCredsFile,
		doc.GetDocFlag("with-client-creds-file").Name,
		doc.GetDocFlag("with-client-creds-file").Default,
		doc.GetDocFlag("with-client-creds-file").Description,
	)
	RootCmd.PersistentFlags().StringVar(
		&clientCredsJSON,
		doc.GetDocFlag("with-client-creds").Name,
		doc.GetDocFlag("with-client-creds").Default,
		doc.GetDocFlag("with-client-creds").Description,
	)
	RootCmd.AddGroup(&cobra.Group{ID: "tdf"})
	RootCmd.PersistentPostRunE = func(cmd *cobra.Command, args []string) error {
		configFile, _ := cmd.Flags().GetString("config-file")
		configKey, _ := cmd.Flags().GetString("config-key")

		cfg, err := config.LoadConfig(configFile, configKey)
		if err != nil {
			return fmt.Errorf("issue loading config: %w", err)
		}
		OtdfctlCfg = *cfg
		return nil
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// The config file and key are defaulted to otdfctl.yaml.
func Execute() {
	RootCmd.Execute()
}
