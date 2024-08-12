package cli

import (
	"strings"

	"golang.org/x/term"
)

func CommaSeparated(values []string) string {
	return "[" + strings.Join(values, ", ") + "]"
}

func TermWidth() int {
	w, _, err := term.GetSize(0)
	if err != nil {
		return 80
	}
	return w
}
