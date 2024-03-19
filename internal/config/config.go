package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/creasty/defaults"
	"github.com/opentdf/tructl/pkg/cli"
	"github.com/spf13/viper"
)

type Output struct {
	Format string `yaml:"format" default:"styled"`
}

type Config struct {
	Output Output `yaml:"output"`
}

type Error string

func (e Error) Error() string {
	return string(e)
}

const (
	ErrLoadingConfig Error = "error loading config"
)

// Load config with viper.
func LoadConfig(key string) (*Config, error) {
	if key == "" {
		key = "tructl"
		slog.Info("LoadConfig: key not provided, using default", "config", key)
	} else {
		slog.Info("LoadConfig", "config", key)
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
	if format == cli.OutputJSON {
		v.Set("output.format", cli.OutputJSON)
	} else {
		v.Set("output.format", cli.OutputStyled)
	}
	viper.WriteConfig()
}
