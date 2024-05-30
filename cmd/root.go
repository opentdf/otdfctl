/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
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
	defaultCmdName = "otdfctl"
	RootCmd        = &man.Docs.GetDoc("<root>").Command
)

func init() {
	doc := man.Docs.GetDoc("<root>")
	RootCmd = &doc.Command
	RootCmd.PersistentFlags().String(
		doc.GetDocFlag("host").Name,
		doc.GetDocFlag("host").Default,
		doc.GetDocFlag("host").Description,
	)
	RootCmd.PersistentFlags().Bool(
		doc.GetDocFlag("tls-no-verify").Name,
		doc.GetDocFlag("tls-no-verify").DefaultAsBool(),
		doc.GetDocFlag("tls-no-verify").Description,
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
}
