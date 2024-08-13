package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/creasty/defaults"
	"github.com/spf13/viper"
)

var AppName = "otdfctl"

var Version = "0.0.0"
var BuildTime = "1970-01-01T00:00:00Z"
var CommitSha = "0000000"

type Output struct {
	Format string `yaml:"format" default:"styled"`
}

type Config struct {
	Output Output `yaml:"output"`
}

// captures all CLI flags that will override pre-specified config values
type ConfigFlagOverrides struct {
	OutputFormatJSON bool
}

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	OutputJSON   = "json"
	OutputStyled = "styled"

	ErrLoadingConfig Error = "error loading config"
)

// Load config with viper.
// TODO force creation of the config in the `~/.config/otdfctl` directory
// TODO the config file in gh is config.yaml -- might want to emulate this
func LoadConfig(file string, key string) (*Config, error) {
	// default the config values if not passed in
	if file == "" && key == "" {
		key = "otdfctl"
		slog.Debug("LoadConfig: file and key not provided, using default file", "config file", file)
	} else {
		slog.Debug("LoadConfig", "config file", file, "config key", key)
	}

	config := &Config{}
	homedir, err := os.UserHomeDir()
	if err != nil {
		return nil, errors.Join(err, ErrLoadingConfig)
	}
	viper.AddConfigPath(fmt.Sprintf("%s/."+key, homedir))
	viper.AddConfigPath("." + key)
	viper.AddConfigPath(".")
	viper.SetConfigName(key)
	viper.SetConfigType("yaml")

	viper.SetEnvPrefix(key)
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	// Allow for a custom config file to be passed in
	// This takes precedence over the AddConfigPath/SetConfigName
	if file != "" {
		viper.SetConfigFile(file)
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, errors.Join(err, ErrLoadingConfig)
	}

	if err := defaults.Set(config); err != nil {
		return nil, errors.Join(err, ErrLoadingConfig)
	}

	err = viper.Unmarshal(config)
	if err != nil {
		return nil, errors.Join(err, ErrLoadingConfig)
	}

	return config, nil
}

func UpdateOutputFormat(cfgKey, format string) {
	v := viper.GetViper()
	format = strings.ToLower(format)
	formatter := "output.format"
	if cfgKey != "" {
		formatter = cfgKey + "." + formatter
	}
	if format == OutputJSON {
		v.Set(formatter, OutputJSON)
	} else {
		v.Set(formatter, OutputStyled)
	}
	viper.WriteConfig()
}
