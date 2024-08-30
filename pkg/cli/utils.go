package cli

import (
	"os"
	"strconv"
	"strings"

	"github.com/opentdf/otdfctl/pkg/config"
	"golang.org/x/term"
)

func CommaSeparated(values []string) string {
	return "[" + strings.Join(values, ", ") + "]"
}

// Returns the terminal width (overridden by the TEST_TERMINAL_WIDTH env var for testing)
func TermWidth() int {
	var (
		w   int
		err error
	)
	testSize := os.Getenv(config.TEST_TERMINAL_WIDTH)
	if testSize == "" {
		w, _, err = term.GetSize(0)
		if err != nil {
			return 80
		}
		return w
	}
	if w, err = strconv.Atoi(testSize); err != nil {
		return 80
	}
	return w
}
