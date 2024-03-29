/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/opentdf/tructl/internal/config"
	"github.com/spf13/cobra"
)

var (
	cfgFile   string
	TructlCfg config.Config

	configFlagOverrides = config.ConfigFlagOverrides{}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "tructl",
	Short: "manage Virtru Data Security Platform",
	Long: `
A command line tool to manage Virtru Data Security Platform.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&configFlagOverrides.OutputFormatJSON, "json", false, "output single command in JSON (overrides configured output format)")
	rootCmd.PersistentFlags().String("host", "localhost:8080", "host:port of the Virtru Data Security Platform gRPC server")
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config-file", "", "config file (default is $HOME/.tructl.yaml)")

	cfg, err := config.LoadConfig("tructl")
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	TructlCfg = *cfg
}
