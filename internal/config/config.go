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
func LoadConfig(key string) (*Config, error) {
	if key == "" {
		key = "otdfctl"
		slog.Debug("LoadConfig: key not provided, using default", "config file", key)
	} else {
		slog.Debug("LoadConfig", "config file", key)
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

func UpdateOutputFormat(format string) {
	v := viper.GetViper()
	format = strings.ToLower(format)
	if format == OutputJSON {
		v.Set("output.format", OutputJSON)
	} else {
		v.Set("output.format", OutputStyled)
	}
	viper.WriteConfig()
}
