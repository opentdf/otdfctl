package profiles

import (
	"testing"

	osprofiles "github.com/jrschumacher/go-osprofiles"
)

type mockNamedProfileExt struct {
	Name      string `json:"name"`
	TestValue string `json:"test_value"`
}

func (m mockNamedProfileExt) GetName() string {
	return m.Name
}

func Test_isProfileCLI(t *testing.T) {
	// Good case
	var profileI osprofiles.NamedProfile
	profileI = &ProfileCLI{}
	prof, err := isProfileCLI(profileI)
	if err != nil {
		t.Errorf("isProfileCLI() error = %v, wantErr %v", err, false)
	}
	if prof == nil {
		t.Errorf("isProfileCLI() profile should not be nil")
	}

	// Bad case
	profileI = &mockNamedProfileExt{}
	p, err := isProfileCLI(profileI)
	if err == nil {
		t.Errorf("isProfileCLI() error = %v, wantErr %v", err, true)
	}
	if p != nil {
		t.Errorf("isProfileCLI() got = %v, want %v", p, nil)
	}
}
