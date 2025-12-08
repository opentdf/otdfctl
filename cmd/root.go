package cmd

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/opentdf/otdfctl/cmd/auth"
	configCmd "github.com/opentdf/otdfctl/cmd/config"
	"github.com/opentdf/otdfctl/cmd/dev"
	"github.com/opentdf/otdfctl/cmd/policy"
	"github.com/opentdf/otdfctl/cmd/tdf"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/man"
	"github.com/opentdf/platform/sdk"
	"github.com/spf13/cobra"
)

var (
	clientCredsFile string
	clientCredsJSON string

	RootCmd = &man.Docs.GetDoc("<root>").Command
)

type version struct {
	AppName       string `json:"app_name"`
	Version       string `json:"version"`
	CommitSha     string `json:"commit_sha"`
	BuildTime     string `json:"build_time"`
	SDKVersion    string `json:"sdk_version"`
	SchemaVersion string `json:"schema_version"`
}

func init() {
	rootCmd := man.Docs.GetCommand("<root>", man.WithRun(func(cmd *cobra.Command, args []string) {
		c := cli.New(cmd, args)

		if c.Flags.GetOptionalBool("version") {
			v := version{
				AppName:       config.AppName,
				Version:       config.Version,
				CommitSha:     config.CommitSha,
				BuildTime:     config.BuildTime,
				SDKVersion:    sdk.Version,
				SchemaVersion: sdk.TDFSpecVersion,
			}

			version := fmt.Sprintf("%s version %s (%s) %s", config.AppName, config.Version, config.BuildTime, config.CommitSha)
			slog.Debug(version)
			c.ExitWithStyled(version)
			c.ExitWithJSON(v)
			return
		}

		//nolint:errcheck // error does not need to be checked
		cmd.Help()
	}))

	RootCmd = &rootCmd.Command

	// Run logger setup for all commands
	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		c := cli.New(cmd, args)

		// log-level from flag will take precedence over env var
		if logLevelStr := c.Flags.GetOptionalString("log-level"); logLevelStr != "" {
			l := new(slog.LevelVar)
			if err := l.UnmarshalText([]byte(logLevelStr)); err != nil {
				return fmt.Errorf("invalid log level: %s", logLevelStr)
			}
			logger := slog.New(slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
				Level: l,
			}))

			slog.SetDefault(logger)
		}
		return nil
	}

	RootCmd.AddCommand(
		// config
		configCmd.Cmd,
		// tdf
		tdf.EncryptCmd,
		tdf.DecryptCmd,
		tdf.InspectCmd,
		// auth
		auth.Cmd,
		// policy
		policy.Cmd,
		// dev
		dev.Cmd,
	)

	RootCmd.Flags().Bool(
		rootCmd.GetDocFlag("version").Name,
		rootCmd.GetDocFlag("version").DefaultAsBool(),
		rootCmd.GetDocFlag("version").Description,
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
	RootCmd.PersistentFlags().String(
		rootCmd.GetDocFlag("with-access-token").Name,
		rootCmd.GetDocFlag("with-access-token").Default,
		rootCmd.GetDocFlag("with-access-token").Description,
	)
	RootCmd.AddGroup(&cobra.Group{ID: tdf.GroupID})

	// Initialize all subcommands that have been refactored to use explicit initialization
	auth.InitCommands()
	configCmd.InitCommands()
	policy.InitCommands()
	dev.InitCommands()
	tdf.InitEncryptCommand()
	tdf.InitDecryptCommand()
	tdf.InitInspectCommand()
	InitProfileCommands()

	// Add interactive command
	RootCmd.AddCommand(newInteractiveCmd())
}
