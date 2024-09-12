package profiles

import (
	"errors"
)

// TODO:
// - add a version
// - add a migration path

const (
	STORE_KEY_PROFILE = "profile"
	STORE_KEY_GLOBAL  = "global"
)

type profileConfig struct {
	driver ProfileDriver
}

type Profile struct {
	config profileConfig

	globalStore         *GlobalStore
	currentProfileStore *ProfileStore
}

type CurrentProfileStore struct {
	StoreInterface
	ProfileConfig
}

const (
	PROFILE_DRIVER_KEYRING   ProfileDriver = "keyring"
	PROFILE_DRIVER_IN_MEMORY ProfileDriver = "in-memory"
	PROFILE_DRIVER_DEFAULT                 = PROFILE_DRIVER_KEYRING
)

type (
	profileConfigVariadicFunc func(profileConfig) profileConfig
	ProfileDriver             string
)

func WithInMemoryStore() profileConfigVariadicFunc {
	return func(c profileConfig) profileConfig {
		c.driver = PROFILE_DRIVER_IN_MEMORY
		return c
	}
}

func WithKeyringStore() profileConfigVariadicFunc {
	return func(c profileConfig) profileConfig {
		c.driver = PROFILE_DRIVER_KEYRING
		return c
	}
}

func newStoreFactory(driver ProfileDriver) NewStoreInterface {
	switch driver {
	case PROFILE_DRIVER_KEYRING:
		return NewKeyringStore
	case PROFILE_DRIVER_IN_MEMORY:
		return NewMemoryStore
	default:
		return nil
	}
}

// create a new profile and load global config
func New(opts ...profileConfigVariadicFunc) (*Profile, error) {
	var err error

	if testProfile != nil {
		return testProfile, nil
	}

	config := profileConfig{
		driver: PROFILE_DRIVER_DEFAULT,
	}
	for _, opt := range opts {
		config = opt(config)
	}

	// check if the store driver is valid
	newStore := newStoreFactory(config.driver)
	if newStore == nil {
		return nil, errors.New("invalid store driver")
	}

	p := &Profile{
		config: config,
	}

	// load global config
	p.globalStore, err = LoadGlobalConfig(newStore)
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Profile) GetGlobalConfig() *GlobalStore {
	return p.globalStore
}

func (p *Profile) AddProfile(profileName string, endpoint string, tlsNoVerify bool, setDefault bool) error {
	var err error

	// check if profile already exists
	if p.globalStore.ProfileExists(profileName) {
		return errors.New("profile already exists")
	}

	// Create profile store and save
	p.currentProfileStore, err = NewProfileStore(newStoreFactory(p.config.driver), profileName, endpoint, tlsNoVerify)
	if err != nil {
		return err
	}
	if err := p.currentProfileStore.Save(); err != nil {
		return err
	}

	// add profile to global config
	if err := p.globalStore.AddProfile(profileName); err != nil {
		return err
	}

	if setDefault || p.globalStore.GetDefaultProfile() == "" {
		return p.globalStore.SetDefaultProfile(profileName)
	}

	return nil
}

func (p *Profile) GetCurrentProfile() (*ProfileStore, error) {
	if p.currentProfileStore == nil {
		return nil, errors.New("no current profile set")
	}

	return p.currentProfileStore, nil
}

func (p *Profile) GetProfile(profileName string) (*ProfileStore, error) {
	if !p.globalStore.ProfileExists(profileName) {
		return nil, errors.New("profile does not exist")
	}

	return LoadProfileStore(newStoreFactory(p.config.driver), profileName)
}

func (p *Profile) ListProfiles() []string {
	return p.globalStore.ListProfiles()
}

func (p *Profile) UseProfile(profileName string) (*ProfileStore, error) {
	var err error

	// check if current profile is already set
	if p.currentProfileStore != nil {
		if p.currentProfileStore.config.Name == profileName {
			return p.currentProfileStore, nil
		}
	}

	// set current profile
	p.currentProfileStore, err = p.GetProfile(profileName)
	return p.currentProfileStore, err
}

func (p *Profile) UseDefaultProfile() (*ProfileStore, error) {
	defaultProfile := p.globalStore.GetDefaultProfile()
	if defaultProfile == "" {
		return nil, errors.New("no default profile set")
	}

	return p.UseProfile(defaultProfile)
}

func (p *Profile) SetDefaultProfile(profileName string) error {
	if !p.globalStore.ProfileExists(profileName) {
		return errors.New("profile does not exist")
	}

	return p.globalStore.SetDefaultProfile(profileName)
}

func (p *Profile) DeleteProfile(profileName string) error {
	// check if profile exists
	if !p.globalStore.ProfileExists(profileName) {
		return errors.New("profile does not exist")
	}

	// get profile
	profile, err := LoadProfileStore(newStoreFactory(p.config.driver), profileName)
	if err != nil {
		return err
	}

	// remove profile from global config (will error if profile is default)
	if err := p.globalStore.RemoveProfile(profileName); err != nil {
		return err
	}

	// delete profile config
	err = profile.Delete()
	if err != nil {
		return err
	}

	return nil
}
