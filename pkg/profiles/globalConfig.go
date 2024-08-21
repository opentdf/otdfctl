package profiles

type GlobalStore struct {
	*Store

	config GlobalConfig
}

type GlobalConfig struct {
	Profiles       []string `json:"profiles"`
	DefaultProfile string   `json:"defaultProfile"`
}

// load global config or create a new one
func LoadGlobalConfig() (*GlobalStore, error) {
	p := &GlobalStore{
		Store: NewStore(STORE_NAMESPACE, STORE_KEY_GLOBAL),

		config: GlobalConfig{
			Profiles:       make([]string, 0),
			DefaultProfile: "",
		},
	}

	if p.Exists() {
		err := p.Get(&p.config)
		return p, err
	}

	err := p.Set(p.config)
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
	return p.Set(p.config)
}

func (p *GlobalStore) ListProfiles() []string {
	return p.config.Profiles
}

func (p *GlobalStore) RemoveProfile(profileName string) error {
	for i, profile := range p.config.Profiles {
		if profile == profileName {
			p.config.Profiles = append(p.config.Profiles[:i], p.config.Profiles[i+1:]...)
			return p.Set(p.config)
		}
	}
	return nil
}

func (p *GlobalStore) SetDefaultProfile(profileName string) error {
	p.config.DefaultProfile = profileName
	return p.Set(p.config)
}

func (p *GlobalStore) GetDefaultProfile() string {
	return p.config.DefaultProfile
}
