package config

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/creasty/defaults"
	profiles "github.com/jrschumacher/go-osprofiles/pkg/platform"
	"github.com/spf13/viper"
)

var (
	// AppName is the name of the application
	// Note: use caution when renaming as it is used in various places within the CLI including for
	// config file naming and in the profile store
	AppName = "otdfctl"

	Version   = "0.0.0"
	BuildTime = "1970-01-01T00:00:00Z"
	CommitSha = "0000000"

	// Test mode is used to determine if the application is running in test mode
	//   "true" = running in test mode
	TestMode = ""

	// Test terminal size is a runtime env var to allow for testing of terminal output
	TEST_TERMINAL_WIDTH = "TEST_TERMINAL_WIDTH"
)

// Profile storage drivers
const (
	_ = iota
	ProfileStoreInMemory
	ProfileStoreFile
	ProfileStoreNativeKeyring
)

type Output struct {
	Format string `yaml:"format" default:"styled"`
}

type Config struct {
	Output Output `yaml:"output"`
	// TODO: make this actually configurable by flags/env
	// ProfileStoreType is the type of profile store to use
	ProfileStoreType int
	// Configured file system directory to store profiles (if using file system profile driver)
	ProfileStoreDir string
}

// captures all CLI flags that will override pre-specified config values
type ConfigFlagOverrides struct {
	OutputFormatJSON bool
	// TODO: add profile store type override
}

const (
	OutputJSON   = "json"
	OutputStyled = "styled"
)

var ErrLoadingConfig = errors.New("error loading config")

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

	platOS, err := profiles.NewPlatform(AppName, runtime.GOOS)
	if err != nil {
		return nil, errors.Join(err, ErrLoadingConfig)
	}
	config.ProfileStoreDir = platOS.GetDataDirectory()

	// TODO: do not hardcode!
	config.ProfileStoreType = ProfileStoreNativeKeyring

	// Override platform-native profile store directory if set in the environment
	profilePathKey := strings.ToUpper(fmt.Sprintf("%s_PROFILE_STORE_DIR", AppName))
	profilePath := os.Getenv(profilePathKey)
	if profilePath != "" {
		config.ProfileStoreDir = profilePath
	}

	return config, nil
}

func UpdateOutputFormat(cfgKey, format string) error {
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
	return viper.WriteConfig()
}
