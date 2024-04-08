package cli

import "strings"

func CommaSeparated(values []string) string {
	return "[" + strings.Join(values, ", ") + "]"
}
