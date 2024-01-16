package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type FlagHelperStringSliceOptions struct {
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

func (f FlagHelper) GetRequiredStringSlice(flag string, v []string, opts FlagHelperStringSliceOptions) []string {
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
