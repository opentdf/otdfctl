package profiles

import (
	"errors"
	"runtime"
	"strings"

	osprofiles "github.com/jrschumacher/go-osprofiles"
	osplatform "github.com/jrschumacher/go-osprofiles/pkg/platform"
	"github.com/opentdf/otdfctl/pkg/config"
)

type ProfileDriver string

const (
	PROFILE_DRIVER_KEYRING     ProfileDriver = "keyring"
	PROFILE_DRIVER_IN_MEMORY   ProfileDriver = "in-memory"
	PROFILE_DRIVER_FILE_SYSTEM ProfileDriver = "filesystem"
	PROFILE_DRIVER_UNKNOWN     ProfileDriver = "unknown"
	PROFILE_DRIVER_DEFAULT                   = PROFILE_DRIVER_FILE_SYSTEM
)

func newFileStoreProfiler() (*osprofiles.Profiler, error) {
	platform, err := osplatform.NewPlatform(config.ServicePublisher, config.AppName, runtime.GOOS)
	if err != nil {
		return nil, errors.Join(ErrCreatingPlatform, err)
	}
	profiler, err := osprofiles.New(config.AppName, osprofiles.WithFileStore(platform.UserAppConfigDirectory()))
	if err != nil {
		return nil, errors.Join(ErrCreatingNewProfile, err)
	}
	return profiler, nil
}

func NewProfiler(store string) (*osprofiles.Profiler, error) {
	driverType, err := ToProfileDriver(store)
	if err != nil {
		return nil, err
	}

	return CreateProfiler(driverType)
}

func ToProfileDriver(driverType string) (ProfileDriver, error) {
	normalizedType := strings.ToLower(strings.TrimSpace(driverType))
	switch normalizedType {
	case string(PROFILE_DRIVER_IN_MEMORY):
		return PROFILE_DRIVER_IN_MEMORY, nil
	case string(PROFILE_DRIVER_KEYRING):
		return PROFILE_DRIVER_KEYRING, nil
	case string(PROFILE_DRIVER_FILE_SYSTEM):
		return PROFILE_DRIVER_FILE_SYSTEM, nil
	default:
		return PROFILE_DRIVER_UNKNOWN, ErrUnknownProfileDriverType
	}
}

func CreateProfiler(driverType ProfileDriver) (*osprofiles.Profiler, error) {
	switch driverType {
	case PROFILE_DRIVER_IN_MEMORY:
		return osprofiles.New(config.AppName, osprofiles.WithInMemoryStore())
	case PROFILE_DRIVER_KEYRING:
		return osprofiles.New(config.AppName, osprofiles.WithKeyringStore())
	case PROFILE_DRIVER_FILE_SYSTEM:
		return newFileStoreProfiler()
	default:
		return nil, ErrUnknownProfileDriverType
	}
}
