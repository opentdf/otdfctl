package man

import "fmt"

type DocFlag struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Shorthand   string `yaml:"shorthand"`
	Default     string `yaml:"default"`
}

func (d *Doc) GetDocFlag(name string) DocFlag {
	for _, f := range d.DocFlags {
		if f.Name == name {
			return f
		}
	}
	panic(fmt.Sprintf("No doc flag found for name, %s for command %s", name, d.Use))
}

func (f DocFlag) DefaultAsBool() bool {
	return f.Default == "true"
}
