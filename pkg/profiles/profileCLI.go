package profiles

import (
	"github.com/opentdf/otdfctl/internal/auth"
	"github.com/opentdf/otdfctl/pkg/utils"
)

type ProfileCLI struct {
	Name        string `json:"profile"`
	Endpoint    string `json:"endpoint"`
	TlsNoVerify bool   `json:"tlsNoVerify"`
	// TODO: use pointer?
	AuthCredentials *auth.AuthCredentials `json:"authCredentials"`
}

// Satisfy go-osprofiles.NamedProfile interface
func (p *ProfileCLI) GetName() string {
	return p.Name
}

// Endpoint
func (p *ProfileCLI) GetEndpoint() string {
	return p.Endpoint
}

func (p *ProfileCLI) SetEndpoint(endpoint string) error {
	u, err := utils.NormalizeEndpoint(endpoint)
	if err != nil {
		return err
	}
	p.Endpoint = u.String()
	return nil
}

// TLS No Verify
func (p *ProfileCLI) GetTLSNoVerify() bool {
	return p.TlsNoVerify
}

func (p *ProfileCLI) SetTLSNoVerify(tlsNoVerify bool) {
	p.TlsNoVerify = tlsNoVerify
}

// AuthCredentials
func (p *ProfileCLI) GetAuthCredentials() *auth.AuthCredentials {
	return p.AuthCredentials
}

func (p *ProfileCLI) SetAuthCredentials(c *auth.AuthCredentials) {
	p.AuthCredentials = c
}
