package profiles

import "github.com/opentdf/otdfctl/pkg/utils"

type ProfileStore struct {
	*Store

	config ProfileConfig
}

type ProfileConfig struct {
	Name            string          `json:"profile"`
	Endpoint        string          `json:"endpoint"`
	tlsNoVerify     bool            `json:"tlsNoVerify"`
	AuthCredentials AuthCredentials `json:"authCredentials"`
}

func NewProfileStore(profileName string, endpoint string) (*ProfileStore, error) {
	if err := validateProfileName(profileName); err != nil {
		return nil, err
	}
	if _, err := utils.NormalizeEndpoint(endpoint); err != nil {
		return nil, err
	}

	p := &ProfileStore{
		Store: NewStore(STORE_NAMESPACE, getStoreKey(profileName)),
		config: ProfileConfig{
			Name:     profileName,
			Endpoint: endpoint,
		},
	}
	return p, nil
}

func LoadProfileStore(profileName string) (*ProfileStore, error) {
	if err := validateProfileName(profileName); err != nil {
		return nil, err
	}

	p := &ProfileStore{
		Store: NewStore(STORE_NAMESPACE, getStoreKey(profileName)),
	}
	return p, p.Get()
}

func (p *ProfileStore) Get() error {
	return p.Store.Get(&p.config)
}

func (p *ProfileStore) Save() error {
	return p.Store.Set(p.config)
}

func (p *ProfileStore) Delete() error {
	return p.Store.Delete()
}

// Profile Name
func (p *ProfileStore) GetProfileName() string {
	return p.config.Name
}

// Endpoint
func (p *ProfileStore) GetEndpoint() string {
	return p.config.Endpoint
}

func (p *ProfileStore) SetEndpoint(endpoint string) error {
	if _, err := utils.NormalizeEndpoint(endpoint); err != nil {
		return err
	}
	p.config.Endpoint = endpoint
	return p.Save()
}

// TLS No Verify
func (p *ProfileStore) GetTLSNoVerify() bool {
	return p.config.tlsNoVerify
}

func (p *ProfileStore) SetTLSNoVerify(tlsNoVerify bool) error {
	p.config.tlsNoVerify = tlsNoVerify
	return p.Save()
}
