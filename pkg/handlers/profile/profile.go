package profile

import (
	"errors"
)

// instead of using the default as the name we'd have a second entry that woul

// otdfctl - global config
//   profiles: []string
//   defaultProfile: string

// otdfctl-<profile> - profile config
//   endpoint: string
//   authCredentials: AuthCredentials
//	 	 authType: string
//     clientCredentials: ClientCredentials

// otdfctl-jake-profile-dev
// otdfctl-jake-profile-staging
// otdfctl-jake-profile-prod

// TODO
// - export profiles to a file (as TDF :D)
// - import profiles from a file (as TDF :D)
// - linux support
// - global logout

const (
	STORE_NAMESPACE   = "otdfctl"
	STORE_KEY_PROFILE = "profile"
	STORE_KEY_GLOBAL  = "global"
)

type Profile struct {
	globalStore         *GlobalStore
	currentProfileStore *ProfileStore
}

type CurrentProfileStore struct {
	*Store

	config ProfileConfig
}

func getStoreKey(n string) string {
	return STORE_KEY_PROFILE + "-" + n
}

func New() (*Profile, error) {
	var err error

	p := &Profile{}

	// load global config
	p.globalStore, err = LoadGlobalConfig()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Profile) AddProfile(profileName string, endpoint string, setDefault bool) error {
	var err error

	// check if profile already exists
	if p.globalStore.ProfileExists(profileName) {
		return errors.New("profile already exists")
	}

	// Create profile store and save
	p.currentProfileStore, err = NewProfileStore(profileName, endpoint)
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

	if setDefault {
		p.globalStore.SetDefaultProfile(profileName)
	}

	return nil
}

func (p *Profile) UseProfile(profileName string) error {
	var err error

	// check if current profile is already set
	if p.currentProfileStore != nil {
		if p.currentProfileStore.config.Name == profileName {
			return nil
		}
	}

	// set current profile
	p.currentProfileStore, err = p.LoadProfile(profileName)
	return err
}

func (p *Profile) UseDefaultProfile() error {
	defaultProfile := p.globalStore.GetDefaultProfile()
	if defaultProfile == "" {
		return errors.New("no default profile set")
	}

	return p.UseProfile(defaultProfile)
}

func (p *Profile) LoadProfile(profileName string) (*ProfileStore, error) {
	if !p.globalStore.ProfileExists(profileName) {
		return nil, errors.New("profile does not exist")
	}

	return LoadProfileStore(profileName)
}

func (p *Profile) DeleteProfile(profileName string) error {
	// check if profile exists
	if !p.globalStore.ProfileExists(profileName) {
		return errors.New("profile does not exist")
	}

	// get profile
	profile, err := LoadProfileStore(profileName)
	if err != nil {
		return err
	}

	// delete profile config
	err = profile.Delete()
	if err != nil {
		return err
	}

	// remove profile from global config
	if err := p.globalStore.RemoveProfile(profileName); err != nil {
		return err
	}

	return nil
}

func (p *Profile) GlobalConfig() *GlobalStore {
	return p.globalStore
}

func (p *Profile) CurrentProfile() (*ProfileStore, error) {
	if p.currentProfileStore == nil {
		return nil, errors.New("no profile loaded")
	}

	return p.currentProfileStore, nil
}
