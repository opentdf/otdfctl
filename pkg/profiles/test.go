package profiles

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/opentdf/otdfctl/pkg/config"
	"github.com/zalando/go-keyring"
)

const testModeMsg = `
********************
RUNNING IN TEST MODE

test config: %s
********************

`

var (
	testProfile *Profile
	testCfg     = os.Getenv("OTDFCTL_TEST_PROFILE")
)

type testConfig struct {
	// global config is used to get the store in a bad state
	GlobalConfig config.Config `json:"globalConfig,omitempty"`

	// set the default profile
	DefaultProfile string `json:"defaultProfile,omitempty"`

	// profiles to add
	Profiles []ProfileConfig `json:"profiles,omitempty"`
}

func init() {
	// If running in test mode, use the mock keyring
	//nolint:nestif,forbidigo // test mode mocking so nested blocks and format directive make sense
	if config.TestMode == "true" {
		fmt.Printf(testModeMsg, testCfg)

		keyring.MockInit()

		// configure the keyring based on the test config
		// unmarsal the test config
		if testCfg != "" {
			var err error
			var cfg testConfig
			//nolint:musttag // test config is annotated and this is a linter issue?
			if err := json.NewDecoder(bytes.NewReader([]byte(testCfg))).Decode(&cfg); err != nil {
				panic(err)
			}

			testProfile, err = New()
			if err != nil {
				panic(err)
			}

			for _, p := range cfg.Profiles {
				err := testProfile.AddProfile(p.Name, p.Endpoint, p.TlsNoVerify, cfg.DefaultProfile == p.Name)
				if err != nil {
					panic(err)
				}
			}

			// set default
			if cfg.DefaultProfile != "" {
				if err := testProfile.SetDefaultProfile(cfg.DefaultProfile); err != nil {
					panic(err)
				}
			}
		}
	}
}
