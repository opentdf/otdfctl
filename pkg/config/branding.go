package config

import (
	"os"
	"path/filepath"
)

// DisplayName is a user-facing CLI name that can be overridden at build time.
//
// Keep AppName stable for config/profile paths; use CLIName() for UX strings.
// Example override: go build -ldflags "-X github.com/opentdf/otdfctl/pkg/config.DisplayName=tructl"
var DisplayName = ""

func CLIName() string {
	if DisplayName != "" {
		return DisplayName
	}

	if len(os.Args) > 0 {
		if base := filepath.Base(os.Args[0]); base != "" && base != "." && base != string(filepath.Separator) {
			return base
		}
	}

	return AppName
}
