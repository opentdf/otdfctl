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
	cfgFile    string
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
	RootCmd.PersistentFlags().StringVar(
		&cfgFile,
		doc.GetDocFlag("config-file").Name,
		doc.GetDocFlag("config-file").Default,
		doc.GetDocFlag("config-file").Description,
	)
	RootCmd.PersistentFlags().String(
		doc.GetDocFlag("log-level").Name,
		doc.GetDocFlag("log-level").Default,
		doc.GetDocFlag("log-level").Description,
	)

	cfg, err := config.LoadConfig("otdfctl")
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	OtdfctlCfg = *cfg
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
