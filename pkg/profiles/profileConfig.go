package profiles

import (
	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/utils"
)

type ProfileStore struct {
	store StoreInterface

	config ProfileConfig
}

type ProfileConfig struct {
	Name            string          `json:"p"`
	Endpoint        string          `json:"e"`
	TlsNoVerify     bool            `json:"t"`
	AuthCredentials AuthCredentials `json:"a"`
}

func NewProfileStore(newStore NewStoreInterface, profileName string, endpoint string, tlsNoVerify bool) (*ProfileStore, error) {
	if err := validateProfileName(profileName); err != nil {
		return nil, err
	}
	u, err := utils.NormalizeEndpoint(endpoint)
	if err != nil {
		return nil, err
	}

	p := &ProfileStore{
		store: newStore(config.AppName, getStoreKey(profileName)),
		config: ProfileConfig{
			Name:        profileName,
			Endpoint:    u.String(),
			TlsNoVerify: tlsNoVerify,
		},
	}
	return p, nil
}

func LoadProfileStore(newStore NewStoreInterface, profileName string) (*ProfileStore, error) {
	if err := validateProfileName(profileName); err != nil {
		return nil, err
	}

	p := &ProfileStore{
		store: newStore(config.AppName, getStoreKey(profileName)),
	}
	return p, p.Get()
}

func (p *ProfileStore) Get() error {
	return p.store.Get(&p.config)
}

func (p *ProfileStore) Save() error {
	return p.store.Set(p.config)
}

func (p *ProfileStore) Delete() error {
	return p.store.Delete()
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
	u, err := utils.NormalizeEndpoint(endpoint)
	if err != nil {
		return err
	}
	p.config.Endpoint = u.String()
	return p.Save()
}

// TLS No Verify
func (p *ProfileStore) GetTLSNoVerify() bool {
	return p.config.TlsNoVerify
}

func (p *ProfileStore) SetTLSNoVerify(tlsNoVerify bool) error {
	p.config.TlsNoVerify = tlsNoVerify
	return p.Save()
}

// utility functions

func getStoreKey(n string) string {
	return STORE_KEY_PROFILE + "-" + n
}
