package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/opentdf/platform/protocol/go/common"
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

func (f FlagHelper) GetOptionalString(flag string) string {
	return f.cmd.Flag(flag).Value.String()
}

func (f FlagHelper) GetStringSlice(flag string, v []string, opts FlagHelperStringSliceOptions) []string {
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

// Transforms into enum value and defaults to active state
func GetState(cmd *cobra.Command) common.ActiveStateEnum {
	state := common.ActiveStateEnum_ACTIVE_STATE_ENUM_ACTIVE
	stateFlag := strings.ToUpper(cmd.Flag("state").Value.String())
	if stateFlag != "" {
		if stateFlag == "INACTIVE" {
			state = common.ActiveStateEnum_ACTIVE_STATE_ENUM_INACTIVE
		} else if stateFlag == "ANY" {
			state = common.ActiveStateEnum_ACTIVE_STATE_ENUM_ANY
		}
	}
	return state
}

// func (f FlagHelper) GetStructSlice(flag string, v []StructFlag[T], opts FlagHelperStringSliceOptions) ([]StructFlag[T], err) {
// 	if len(v) < opts.Min {
// 		fmt.Println(ErrorMessage(fmt.Sprintf("Flag %s must have at least %d non-empty values", flag, opts.Min), nil))
// 		os.Exit(1)
// 	}
// 	if opts.Max > 0 && len(v) > opts.Max {
// 		fmt.Println(ErrorMessage(fmt.Sprintf("Flag %s must have at most %d non-empty values", flag, opts.Max), nil))
// 		os.Exit(1)
// 	}
// 	return v
// }

// type StructFlag[T any] struct {
// 	Val T
// }

// func (this StructFlag[T]) String() string {
// 	b, _ := json.Marshal(this)
// 	return string(b)
// }

// func (this StructFlag[T]) Set(s string) error {
// 	return json.Unmarshal([]byte(s), this)
// }
