/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/opentdf/otdfctl/internal/config"
	"github.com/opentdf/otdfctl/pkg/man"
)

var (
	cfgKey     string
	OtdfctlCfg config.Config

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
		doc.GetDocFlag("host").Name,
		doc.GetDocFlag("host").Default,
		doc.GetDocFlag("host").Description,
	)
	RootCmd.PersistentFlags().String(
		doc.GetDocFlag("log-level").Name,
		doc.GetDocFlag("log-level").Default,
		doc.GetDocFlag("log-level").Description,
	)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// The config file and key are defaulted to otdfctl.yaml.
func Execute() {
	ExecuteWithBootstrap("", "")
}

// Execute adds all child commands to the root command and sets flags appropriately.
// It also allows the config file & key to be bootstrapped for wrapping the CLI.
func ExecuteWithBootstrap(configFile, configKey string) {
	cfgKey = configKey
	cfg, err := config.LoadConfig(configFile, configKey)
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	OtdfctlCfg = *cfg
	err = RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
