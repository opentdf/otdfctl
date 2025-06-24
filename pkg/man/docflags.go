package man

import (
	"fmt"

	"github.com/opentdf/otdfctl/pkg/cli"
)

type DocFlag struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Shorthand   string   `yaml:"shorthand"`
	Default     string   `yaml:"default"`
	Enum        []string `yaml:"enum"`
}

type DocArgument struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Required    bool   `yaml:"required"`
	Type        string `yaml:"type"`
}

// FlexibleArgs can handle both []string and []DocArgument formats
type FlexibleArgs struct {
	StringArgs []string
	ObjectArgs []DocArgument
}

func (f *FlexibleArgs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try to unmarshal as []string first
	var stringArgs []string
	if err := unmarshal(&stringArgs); err == nil {
		f.StringArgs = stringArgs
		return nil
	}

	// If that fails, try to unmarshal as []DocArgument
	var objectArgs []DocArgument
	if err := unmarshal(&objectArgs); err == nil {
		f.ObjectArgs = objectArgs
		return nil
	}

	return fmt.Errorf("arguments must be either a list of strings or a list of argument objects")
}

// ToStringSlice converts FlexibleArgs to a simple string slice for compatibility
func (f *FlexibleArgs) ToStringSlice() []string {
	if len(f.StringArgs) > 0 {
		return f.StringArgs
	}
	
	result := make([]string, len(f.ObjectArgs))
	for i, arg := range f.ObjectArgs {
		result[i] = arg.Name
	}
	return result
}

func (d *Doc) GetDocFlag(name string) DocFlag {
	for _, f := range d.DocFlags {
		if f.Name == name {
			if len(f.Enum) > 0 {
				f.Description = fmt.Sprintf("%s %s", f.Description, cli.CommaSeparated(f.Enum))
			}
			return f
		}
	}
	panic(fmt.Sprintf("No doc flag found for name, %s for command %s", name, d.Use))
}

func (f DocFlag) DefaultAsBool() bool {
	return f.Default == "true"
}
