package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/opentdf/platform/protocol/go/common"
	"github.com/spf13/cobra"
)

type FlagsStringSliceOptions struct {
	Min int
	Max int
}

type flagHelper struct {
	cmd *cobra.Command
}

func newFlagHelper(cmd *cobra.Command) *flagHelper {
	return &flagHelper{cmd: cmd}
}

func (f flagHelper) GetRequiredString(flag string) string {
	v := f.cmd.Flag(flag).Value.String()
	if v == "" {
		fmt.Println(ErrorMessage("Flag "+flag+" is required", nil))
		os.Exit(1)
	}
	return v
}

func (f flagHelper) GetOptionalString(flag string) string {
	p := f.cmd.Flag(flag)
	if p == nil {
		return ""
	}
	return p.Value.String()
}

func (f flagHelper) GetStringSlice(flag string, v []string, opts FlagsStringSliceOptions) []string {
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

func (f flagHelper) GetRequiredInt32(flag string) int32 {
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

func (f flagHelper) GetOptionalBool(flag string) bool {
	v, _ := f.cmd.Flags().GetBool(flag)
	return v
}

func (f flagHelper) GetRequiredBool(flag string) bool {
	v, e := f.cmd.Flags().GetBool(flag)
	if e != nil {
		fmt.Println(ErrorMessage("Flag "+flag+" is required", nil))
		os.Exit(1)
	}
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

// func (f flagHelper) GetStructSlice(flag string, v []StructFlag[T], opts flagHelperStringSliceOptions) ([]StructFlag[T], err) {
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
