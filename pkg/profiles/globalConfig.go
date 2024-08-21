package profiles

import "github.com/opentdf/otdfctl/pkg/config"

// This variable is used to store the version of the profiles. Since the profiles structure might
// change in the future, this variable is used to keep track of the version of the profiles and will
// be used to determine how to handle migration of the profiles.
const PROFILES_VERSION_v1_0 = "1.0"

const PROFILES_VERSION_LATEST = PROFILES_VERSION_v1_0

type GlobalStore struct {
	store StoreInterface

	config GlobalConfig
}

type GlobalConfig struct {
	ProfilesVersion string   `json:"version"`
	Profiles        []string `json:"profiles"`
	DefaultProfile  string   `json:"defaultProfile"`
}

// load global config or create a new one
func LoadGlobalConfig(newStore func(string, string) StoreInterface) (*GlobalStore, error) {
	p := &GlobalStore{
		store: newStore(config.AppName, STORE_KEY_GLOBAL),

		config: GlobalConfig{
			Profiles:       make([]string, 0),
			DefaultProfile: "",
		},
	}

	if p.store.Exists() {
		err := p.store.Get(&p.config)

		// check the version of the profiles
		if p.config.ProfilesVersion != PROFILES_VERSION_LATEST {
			// handle migration of the profiles
			// currently, there is no migration needed
			// so we just set the version to the latest version
			p.config.ProfilesVersion = PROFILES_VERSION_LATEST
			err = p.store.Set(p.config)
			if err != nil {
				return nil, err
			}
		}

		return p, err
	}

	// set the version of the profiles to the latest version
	p.config.ProfilesVersion = PROFILES_VERSION_LATEST
	err := p.store.Set(p.config)
	return p, err
}

func (p *GlobalStore) ProfileExists(profileName string) bool {
	for _, profile := range p.config.Profiles {
		if profile == profileName {
			return true
		}
	}
	return false
}

func (p *GlobalStore) AddProfile(profileName string) error {
	p.config.Profiles = append(p.config.Profiles, profileName)
	return p.store.Set(p.config)
}

func (p *GlobalStore) ListProfiles() []string {
	return p.config.Profiles
}

func (p *GlobalStore) RemoveProfile(profileName string) error {
	for i, profile := range p.config.Profiles {
		if profile == profileName {
			p.config.Profiles = append(p.config.Profiles[:i], p.config.Profiles[i+1:]...)
			return p.store.Set(p.config)
		}
	}
	return nil
}

func (p *GlobalStore) SetDefaultProfile(profileName string) error {
	p.config.DefaultProfile = profileName
	return p.store.Set(p.config)
}

func (p *GlobalStore) GetDefaultProfile() string {
	return p.config.DefaultProfile
}
