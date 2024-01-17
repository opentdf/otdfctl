package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

const (
	splitter = "::"
)

type FlagHelperListOptions struct {
	Min int
	Max int
}

type FlagHelper struct {
	cmd *cobra.Command
}

func NewFlagHelper(cmd *cobra.Command) *FlagHelper {
	return &FlagHelper{cmd: cmd}
}

func (f FlagHelper) GetRequiredString(flag string) string {
	v := f.cmd.Flag(flag).Value.String()
	if v == "" {
		fmt.Println(ErrorMessage("Flag "+flag+" is required", nil))
		os.Exit(1)
	}
	return v
}

func (f FlagHelper) GetOptionalString(flag string) string {
	return f.cmd.Flag(flag).Value.String()
}

func (f FlagHelper) GetStringSlice(flag string, v []string, opts FlagHelperListOptions) []string {
	if len(v) < opts.Min {
		fmt.Println(ErrorMessage(fmt.Sprintf("Flag %s must have at least %d non-empty values", flag, opts.Min), nil))
		os.Exit(1)
	}
	if opts.Max > 0 && len(v) > opts.Max {
		fmt.Println(ErrorMessage(fmt.Sprintf("Flag %s must have at most %d non-empty values", flag, opts.Max), nil))
		os.Exit(1)
	}
	return v
}

func (f FlagHelper) GetRequiredInt32(flag string) int32 {
	v, e := f.cmd.Flags().GetInt32(flag)
	if e != nil {
		fmt.Println(ErrorMessage("Flag "+flag+" is required", nil))
		os.Exit(1)
	}
	// if v == 0 {
	// 	fmt.Println(ErrorMessage("Flag "+flag+" must be greater than 0", nil))
	// 	os.Exit(1)
	// }
	return v
}

func (f FlagHelper) GetKeyValuesMap(flagName string, values []string, opts FlagHelperListOptions) map[string]string {
	if len(values) < opts.Min {
		fmt.Println(ErrorMessage(fmt.Sprintf("Flag %s must have at least %d non-empty values", flagName, opts.Min), nil))
		os.Exit(1)
	}
	if opts.Max > 0 && len(values) > opts.Max {
		fmt.Println(ErrorMessage(fmt.Sprintf("Flag %s must have at most %d non-empty values", flagName, opts.Max), nil))
		os.Exit(1)
	}

	valuesMap := make(map[string]string)
	for _, value := range values {
		k, v := splitKeyValue(value)
		valuesMap[k] = v
	}
	return valuesMap
}

func splitKeyValue(s string) (string, string) {
	parts := strings.SplitN(s, splitter, 2)
	if len(parts) == 2 {
		return parts[0], parts[1]
	}
	return parts[0], ""
}
