package profiles

import (
	"errors"
)

// TODO:
// - add a version
// - add a migration path

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

func (p *Profile) GetGlobalConfig() *GlobalStore {
	return p.globalStore
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

	return LoadProfileStore(profileName)
}

func (p *Profile) ListProfiles() []string {
	return p.globalStore.ListProfiles()
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
	p.currentProfileStore, err = p.GetProfile(profileName)
	return err
}

func (p *Profile) UseDefaultProfile() error {
	defaultProfile := p.globalStore.GetDefaultProfile()
	if defaultProfile == "" {
		return errors.New("no default profile set")
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