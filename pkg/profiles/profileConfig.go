package profiles

import (
	"fmt"
	"time"

	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/opentdf/otdfctl/pkg/utils"
)

// URN-based namespace template without UUID, using only profile name for uniqueness
var URNNamespaceTemplate = fmt.Sprintf("urn:opentdf:%s:profile:v1", config.AppName) // e.g., urn:opentdf:otdfctl:profile:v1:<profileName>

// ProfileStore manages profile configurations and handles storage
type ProfileStore struct {
	store  StoreInterface
	config ProfileConfig
}

// ProfileConfig defines the structure of a profile with flexible attributes and timestamps
type ProfileConfig struct {
	Name            string                 `json:"profile"`         // Profile name (unique identifier)
	Endpoint        string                 `json:"endpoint"`        // Profile endpoint
	TlsNoVerify     bool                   `json:"tlsNoVerify"`     // TLS verification setting
	AuthCredentials AuthCredentials        `json:"authCredentials"` // Authentication credentials
	Attributes      map[string]interface{} `json:"attributes"`      // Flexible map of additional attributes
	CreatedAt       time.Time              `json:"createdAt"`       // Timestamp for profile creation
	UpdatedAt       time.Time              `json:"updatedAt"`       // Timestamp for last profile update
	Version         string                 `json:"version"`         // profile version
}

// NewProfileStore creates a new profile store with flexible attributes and timestamps
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
			Attributes:  make(map[string]interface{}), // Empty map for flexible attributes
			CreatedAt:   time.Now().UTC(),             // Set creation time
			UpdatedAt:   time.Now().UTC(),             // Set initial update time
			Version:     URNNamespaceTemplate,         // Set profile version to URN-based namespace template
		},
	}
	return p, nil
}

// LoadProfileStore loads an existing profile using its profile name
func LoadProfileStore(newStore NewStoreInterface, profileName string) (*ProfileStore, error) {
	if err := validateProfileName(profileName); err != nil {
		return nil, err
	}

	p := &ProfileStore{
		store: newStore(config.AppName, getStoreKey(profileName)),
	}
	return p, p.Get()
}

// Get loads the profile configuration into p.config
func (p *ProfileStore) Get() error {
	return p.store.Get(&p.config)
}

// Save saves the current profile configuration to storage and updates UpdatedAt timestamp
func (p *ProfileStore) Save() error {
	p.config.UpdatedAt = time.Now().UTC()
	return p.store.Set(p.config)
}

// Delete removes the profile from storage
func (p *ProfileStore) Delete() error {
	return p.store.Delete()
}

// Generate a unique namespace for a profile using only the profile name
func (p *ProfileConfig) GetNamespace() string {
	return URNNamespaceTemplate
}

// GetProfileName retrieves the profile name
func (p *ProfileStore) GetProfileName() string {
	return p.config.Name
}

// SetAttribute allows adding or updating an attribute in the profile's Attributes map
func (p *ProfileStore) SetAttribute(key string, value interface{}) error {
	p.config.Attributes[key] = value
	return p.Save()
}

// GetAttribute retrieves an attribute by key from the profile's Attributes map
func (p *ProfileStore) GetAttribute(key string) (interface{}, bool) {
	value, exists := p.config.Attributes[key]
	return value, exists
}

// GetEndpoint retrieves the Endpoint value from ProfileConfig
func (p *ProfileStore) GetEndpoint() string {
	return p.config.Endpoint
}

// SetEndpoint updates the Endpoint in ProfileConfig after normalizing it, then saves the profile
func (p *ProfileStore) SetEndpoint(endpoint string) error {
	u, err := utils.NormalizeEndpoint(endpoint)
	if err != nil {
		return err
	}
	p.config.Endpoint = u.String()
	return p.Save() // Save the updated profile configuration
}

// GetTLSNoVerify retrieves the TlsNoVerify setting from ProfileConfig
func (p *ProfileStore) GetTLSNoVerify() bool {
	return p.config.TlsNoVerify
}

// utility functions

// getStoreKey generates a unique key for storing the profile using the profile name
func getStoreKey(name string) string {
	return fmt.Sprintf("%s-%s", STORE_KEY_PROFILE, name)
}
