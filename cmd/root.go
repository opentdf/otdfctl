/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/opentdf/tructl/internal/config"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/cobra"
)

var cfgFile string

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
	singleOutputJSON := rootCmd.PersistentFlags().Bool("json", false, "Output a single command result in JSON format")
	rootCmd.PersistentFlags().String("host", "localhost:9000", "host:port of the Virtru Data Security Platform gRPC server")

	cfg, err := config.LoadConfig("tructl")
	if err != nil {
		fmt.Println("Error loading config:", err)
		os.Exit(1)
	}
	if strings.ToLower(cfg.Output.Format) == cli.OutputJSON || singleOutputJSON {
		cli.OutputFormat = cli.OutputJSON
	}
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.tructl.yaml)")
}
