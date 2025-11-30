package profiles

import (
	"errors"
	"runtime"

	osprofiles "github.com/jrschumacher/go-osprofiles"
	osplatform "github.com/jrschumacher/go-osprofiles/pkg/platform"
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/utils"
)

type ProfileDriver string

const (
	PROFILE_DRIVER_KEYRING     ProfileDriver = "keyring"
	PROFILE_DRIVER_IN_MEMORY   ProfileDriver = "in-memory"
	PROFILE_DRIVER_FILE_SYSTEM ProfileDriver = "filesystem"
	PROFILE_DRIVER_DEFAULT                   = PROFILE_DRIVER_FILE_SYSTEM
)

type OtdfctlProfileStore struct {
	store    osprofiles.ProfileStore
	config   *ProfileConfig // Pointer to the store.Profile field. Use for sets/reads.
	profiler osprofiles.Profiler
}

type ProfileConfig struct {
	Name            string          `json:"profile"`
	Endpoint        string          `json:"endpoint"`
	TlsNoVerify     bool            `json:"tlsNoVerify"`
	AuthCredentials AuthCredentials `json:"authCredentials"`
}

func (pc *ProfileConfig) GetName() string {
	return pc.Name
}

func createProfiler(profileType ProfileDriver) (*osprofiles.Profiler, error) {
	switch profileType {
	case PROFILE_DRIVER_IN_MEMORY:
		return osprofiles.New(config.AppName, osprofiles.WithInMemoryStore())
	case PROFILE_DRIVER_KEYRING:
		return osprofiles.New(config.AppName, osprofiles.WithKeyringStore())
	default:
		platform, err := osplatform.NewPlatform(config.ServicePublisher, config.AppName, runtime.GOOS)
		if err != nil {
			return nil, errors.Join(ErrCreatingPlatform, err)
		}
		return osprofiles.New(config.AppName, osprofiles.WithFileStore(platform.UserAppConfigDirectory()))
	}
}

func NewOtdfctlProfileStore(profileType ProfileDriver, profileName string, endpoint string, tlsNoVerify, setDefault bool) (*OtdfctlProfileStore, error) {
	profiler, err := createProfiler(profileType)
	if err != nil {
		return nil, err
	}

	u, err := utils.NormalizeEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	p := &ProfileConfig{
		Name:        profileName,
		Endpoint:    u.String(),
		TlsNoVerify: tlsNoVerify,
	}
	err = profiler.AddProfile(p, setDefault)
	if err != nil {
		return nil, err
	}

	store, err := osprofiles.GetProfile[*ProfileConfig](profiler, profileName)
	if err != nil {
		return nil, err
	}

	// Cast Profile to ProfileConfig
	pc, ok := store.Profile.(*ProfileConfig)
	if !ok {
		return nil, errors.Join(ErrProfileIncorrectType, err)
	}

	return &OtdfctlProfileStore{
		store:    *store,
		config:   pc,
		profiler: *profiler,
	}, nil
}

func LoadOtdfctlProfileStore(profiler *osprofiles.Profiler, profileName string) (*OtdfctlProfileStore, error) {
	store, err := osprofiles.GetProfile[*ProfileConfig](profiler, profileName)
	if err != nil {
		return nil, err
	}

	pc, ok := store.Profile.(*ProfileConfig)
	if !ok {
		return nil, errors.Join(ErrProfileIncorrectType, err)
	}

	return &OtdfctlProfileStore{
		store:    *store,
		config:   pc,
		profiler: *profiler,
	}, nil
}

func (p *OtdfctlProfileStore) GetEndpoint() string {
	return p.config.Endpoint
}

func (p *OtdfctlProfileStore) SetEndpoint(endpoint string) error {
	u, err := utils.NormalizeEndpoint(endpoint)
	if err != nil {
		return err
	}

	p.config.Endpoint = u.String()
	return p.store.Save()
}

func (p *OtdfctlProfileStore) GetTLSNoVerify() bool {
	return p.config.TlsNoVerify
}

func (p *OtdfctlProfileStore) SetTLSNoVerify(tlsNoVerify bool) error {
	p.config.TlsNoVerify = tlsNoVerify
	return p.store.Save()
}
