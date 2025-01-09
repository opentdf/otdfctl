package profiles

import (
	"errors"

	osprofiles "github.com/jrschumacher/go-osprofiles"
	"github.com/opentdf/otdfctl/pkg/config"
)

type ProfileManager struct {
	profiler *osprofiles.Profiler
}

// TODO: errors should be named

// NewProfileManager sets up profiler to CRUD ProfileCLI instances to the configured storage location
func NewProfileManager(cfg *config.Config, isInMemoryProfile bool) (*ProfileManager, error) {
	var (
		profiler  *osprofiles.Profiler
		err       error
		namespace = config.AppName
	)

	// allow override
	if isInMemoryProfile {
		profiler, err = osprofiles.New(namespace, osprofiles.WithInMemoryStore())
	}

	// defer to config
	switch cfg.ProfileStoreType {
	case config.ProfileStoreInMemory:
		profiler, err = osprofiles.New(namespace, osprofiles.WithInMemoryStore())
	case config.ProfileStoreFile:
		profiler, err = osprofiles.New(namespace, osprofiles.WithFileStore(cfg.ProfileStoreDir))
	case config.ProfileStoreNativeKeyring:
		profiler, err = osprofiles.New(namespace, osprofiles.WithKeyringStore())
	default:
		err = errors.New("Invalid storage location")
	}

	return &ProfileManager{
		profiler: profiler,
	}, err
}

// AddProfile adds a new profile to the profile store
func (p ProfileManager) AddProfile(profile ProfileCLI, setDefault bool) error {
	return p.profiler.AddProfile(profile, setDefault)
}

// GetCurrentProfile returns the current stored default profile
func (p ProfileManager) GetCurrentProfile() (*ProfileCLI, error) {
	profileName := p.profiler.GetGlobalConfig().GetDefaultProfile()
	store, err := p.profiler.GetProfile(profileName)
	if err != nil {
		return nil, err
	}
	profile, ok := store.Profile.(*ProfileCLI)
	if !ok {
		return nil, errors.New("Profile is not of type ProfileCLI")
	}
	return profile, nil
}

// GetProfile returns the profile stored under the specified name
func (p ProfileManager) GetProfile(name string) (*ProfileCLI, error) {
	store, err := p.profiler.GetProfile(name)
	if err != nil {
		return nil, err
	}
	profile, ok := store.Profile.(*ProfileCLI)
	if !ok {
		return nil, errors.New("Profile is not of type ProfileCLI")
	}
	return profile, nil
}

// ListProfiles returns a list of all stored profiles
func (p ProfileManager) ListProfiles() ([]*ProfileCLI, error) {
	var profiles []*ProfileCLI
	for _, profileName := range p.profiler.ListProfiles() {
		store, err := p.profiler.GetProfile(profileName)
		if err != nil {
			return nil, err
		}
		profile, ok := store.Profile.(*ProfileCLI)
		if !ok {
			return nil, errors.New("Profile is not of type ProfileCLI")
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// UseProfile sets the current profile to the specified profile name
func (p ProfileManager) UseProfile(name string) (*ProfileCLI, error) {
	store, err := p.profiler.UseProfile(name)
	if err != nil {
		return nil, err
	}
	profile, ok := store.Profile.(*ProfileCLI)
	if !ok {
		return nil, errors.New("Profile is not of type ProfileCLI")
	}
	return profile, nil
}

// SetDefaultProfile sets the specified profile name the default profile in the store
func (p ProfileManager) SetDefaultProfile(name string) error {
	return p.profiler.SetDefaultProfile(name)
}

// UpdateProfile updates the specified profile saved to the store
func (p ProfileManager) UpdateProfile(profile *ProfileCLI) error {
	profileStore, err := p.profiler.GetProfile(profile.GetName())
	if err != nil {
		return err
	}
	profileStore.Profile = profile
	return profileStore.Save()
}

// DeleteProfile removes the specified profile from the store
func (p ProfileManager) DeleteProfile(name string) error {
	return p.profiler.DeleteProfile(name)
}
