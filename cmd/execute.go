package cmd

import (
	"os"

	"github.com/opentdf/otdfctl/internal/config"
	"github.com/spf13/cobra"
)

type ExecuteConfig struct {
	configFile string
	configKey  string
	mountTo    *cobra.Command
	renameCmd  *cobra.Command
	cmdName    string
}
type ExecuteOptFunc func(c ExecuteConfig) ExecuteConfig

func WithMountTo(cmd *cobra.Command, renameCmd *cobra.Command) ExecuteOptFunc {
	if cmd == nil {
		panic("cmd is nil")
	}

	return func(c ExecuteConfig) ExecuteConfig {
		c.cmdName = cmd.Use
		if renameCmd.Use != "" {
			c.cmdName = renameCmd.Use
			c.configFile = renameCmd.Use
		}
		c.mountTo = cmd
		c.renameCmd = renameCmd
		return c
	}
}

func Execute(opts ...ExecuteOptFunc) {
	c := ExecuteConfig{}
	for _, opt := range opts {
		c = opt(c)
	}

	cfg, _ := config.LoadConfig(c.configFile, c.configKey)
	// Suppress error for now since config file should be optional
	// if err != nil {
	// 	fmt.Println("Error loading config:", err)
	// 	os.Exit(1)
	// }
	OtdfctlCfg = *cfg

	if c.mountTo != nil {
		MountRoot(c.mountTo, c.renameCmd)
	} else {
		err := RootCmd.Execute()
		if err != nil {
			os.Exit(1)
		}
	}
}
