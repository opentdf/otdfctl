package profiles

import (
	"errors"

	osprofiles "github.com/jrschumacher/go-osprofiles"
	"github.com/opentdf/otdfctl/pkg/config"
)

type ProfileManager struct {
	profiler *osprofiles.Profiler
}

// NewProfileManager sets up profiler to CRUD ProfileCLI instances to the configured storage location
func NewProfileManager(cfg *config.Config, isInMemoryProfile bool) (*ProfileManager, error) {
	var (
		profiler  *osprofiles.Profiler
		err       error
		namespace = config.AppName
	)

	// allow override
	if isInMemoryProfile {
		cfg.ProfileStoreType = config.ProfileStoreInMemory
	}

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
func (p ProfileManager) AddProfile(profile *ProfileCLI, setDefault bool) error {
	return p.profiler.AddProfile(profile, setDefault)
}

// GetCurrentProfile returns the current stored default profile
func (p ProfileManager) GetCurrentProfile() (*ProfileCLI, error) {
	profileName := osprofiles.GetGlobalConfig(p.profiler).GetDefaultProfile()
	store, err := osprofiles.GetProfile[*ProfileCLI](p.profiler, profileName)
	if err != nil {
		return nil, err
	}
	return isProfileCLI(store.Profile)
}

// GetProfile returns the profile stored under the specified name
func (p ProfileManager) GetProfile(name string) (*ProfileCLI, error) {
	store, err := osprofiles.GetProfile[*ProfileCLI](p.profiler, name)
	if err != nil {
		return nil, err
	}
	return isProfileCLI(store.Profile)
}

// ListProfiles returns a list of all stored profiles
func (p ProfileManager) ListProfiles() ([]*ProfileCLI, error) {
	var profiles []*ProfileCLI
	for _, profileName := range osprofiles.ListProfiles(p.profiler) {
		store, err := osprofiles.GetProfile[*ProfileCLI](p.profiler, profileName)
		if err != nil {
			return nil, err
		}
		profile, ok := store.Profile.(*ProfileCLI)
		if !ok {
			return nil, errStoredProfileWrongType
		}
		profiles = append(profiles, profile)
	}
	return profiles, nil
}

// UseProfile sets the current profile to the specified profile name
func (p ProfileManager) UseProfile(name string) (*ProfileCLI, error) {
	store, err := osprofiles.UseProfile[*ProfileCLI](p.profiler, name)
	if err != nil {
		return nil, err
	}
	return isProfileCLI(store.Profile)
}

// SetDefaultProfile sets the specified profile name the default profile in the store
func (p ProfileManager) SetDefaultProfile(name string) error {
	return osprofiles.SetDefaultProfile(p.profiler, name)
}

// UpdateProfile updates the specified profile saved to the store
func (p ProfileManager) UpdateProfile(profile *ProfileCLI) error {
	return osprofiles.UpdateCurrentProfile(p.profiler, profile)
}

// DeleteProfile removes the specified profile from the store
func (p ProfileManager) DeleteProfile(name string) error {
	return osprofiles.DeleteProfile[*ProfileCLI](p.profiler, name)
}

func isProfileCLI(p osprofiles.NamedProfile) (*ProfileCLI, error) {
	profile, ok := p.(*ProfileCLI)
	if !ok {
		return nil, errStoredProfileWrongType
	}
	return profile, nil
}
