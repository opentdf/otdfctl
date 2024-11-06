package profiles

import (
	"errors"
)

// Define constants for the different storage drivers and store keys
const (
	PROFILE_DRIVER_KEYRING   ProfileDriver = "keyring"
	PROFILE_DRIVER_IN_MEMORY ProfileDriver = "in-memory"
	PROFILE_DRIVER_FILE      ProfileDriver = "file"
	PROFILE_DRIVER_DEFAULT                 = PROFILE_DRIVER_FILE

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

// Variadic functions to set different storage drivers
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

func WithFileStore() profileConfigVariadicFunc {
	return func(c profileConfig) profileConfig {
		c.driver = PROFILE_DRIVER_FILE
		return c
	}
}

// newStoreFactory returns a storage interface based on the configured driver
func newStoreFactory(driver ProfileDriver) NewStoreInterface {
	switch driver {
	case PROFILE_DRIVER_KEYRING:
		return NewKeyringStore
	case PROFILE_DRIVER_IN_MEMORY:
		return NewMemoryStore
	case PROFILE_DRIVER_FILE:
		return NewFileStore
	default:
		return nil
	}
}

// New creates a new Profile with the specified configuration options
func New(opts ...profileConfigVariadicFunc) (*Profile, error) {
	var err error

	config := profileConfig{
		driver: PROFILE_DRIVER_DEFAULT,
	}
	for _, opt := range opts {
		config = opt(config)
	}

	// Validate and initialize the store
	newStore := newStoreFactory(config.driver)
	if newStore == nil {
		return nil, errors.New("invalid store driver")
	}

	p := &Profile{
		config: config,
	}

	// Load global configuration
	p.globalStore, err = LoadGlobalConfig(newStore)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// GetGlobalConfig returns the global configuration
func (p *Profile) GetGlobalConfig() *GlobalStore {
	return p.globalStore
}

// AddProfile adds a new profile to the current configuration
func (p *Profile) AddProfile(profileName, endpoint string, tlsNoVerify, setDefault bool) error {
	var err error

	// Check if the profile already exists
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

	// Add profile to global configuration
	if err := p.globalStore.AddProfile(profileName); err != nil {
		return err
	}

	if setDefault || p.globalStore.GetDefaultProfile() == "" {
		return p.globalStore.SetDefaultProfile(profileName)
	}

	return nil
}

// GetCurrentProfile retrieves the current profile store
func (p *Profile) GetCurrentProfile() (*ProfileStore, error) {
	if p.currentProfileStore == nil {
		return nil, errors.New("no current profile set")
	}
	return p.currentProfileStore, nil
}

// GetProfile retrieves a specified profile
func (p *Profile) GetProfile(profileName string) (*ProfileStore, error) {
	if !p.globalStore.ProfileExists(profileName) {
		return nil, errors.New("profile does not exist")
	}
	return LoadProfileStore(newStoreFactory(p.config.driver), profileName)
}

// ListProfiles returns a list of profile names
func (p *Profile) ListProfiles() []string {
	return p.globalStore.ListProfiles()
}

// UseProfile sets the current profile to the specified profile name
func (p *Profile) UseProfile(profileName string) (*ProfileStore, error) {
	var err error

	// If current profile is already set to this, return it
	if p.currentProfileStore != nil && p.currentProfileStore.config.Name == profileName {
		return p.currentProfileStore, nil
	}

	// Set current profile
	p.currentProfileStore, err = p.GetProfile(profileName)
	return p.currentProfileStore, err
}

// UseDefaultProfile sets the current profile to the default profile
func (p *Profile) UseDefaultProfile() (*ProfileStore, error) {
	defaultProfile := p.globalStore.GetDefaultProfile()
	if defaultProfile == "" {
		return nil, errors.New("no default profile set")
	}
	return p.UseProfile(defaultProfile)
}

// SetDefaultProfile sets a specified profile as the default profile
func (p *Profile) SetDefaultProfile(profileName string) error {
	if !p.globalStore.ProfileExists(profileName) {
		return errors.New("profile does not exist")
	}
	return p.globalStore.SetDefaultProfile(profileName)
}

// DeleteProfile removes a profile from storage
func (p *Profile) DeleteProfile(profileName string) error {
	// Check if the profile exists
	if !p.globalStore.ProfileExists(profileName) {
		return errors.New("profile does not exist")
	}

	// Retrieve the profile
	profile, err := LoadProfileStore(newStoreFactory(p.config.driver), profileName)
	if err != nil {
		return err
	}

	// Remove profile from global configuration
	if err := p.globalStore.RemoveProfile(profileName); err != nil {
		return err
	}

	// Delete profile configuration
	return profile.Delete()
}
