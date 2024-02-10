package man

import (
	_ "embed"
	"fmt"
	"strings"

	"github.com/adrg/frontmatter"
)

//go:embed policy-resourceMappings.md
var policyResourceMappingsEn string

//go:embed policy-resourceMappings.fr.md
var policyResourceMappingsFr string

var PolicyResourceMappings = map[string]Doc{}

type Doc struct {
	Command string   `yaml:"command"`
	Aliases []string `yaml:"aliases"`
	Short   string   `yaml:"short"`
	Long    string
}

func (m Doc) ShortWithSubCommands(subcommands []string) string {
	return m.Short + " [" + strings.Join(subcommands, ", ") + "]"
}

func init() {
	PolicyResourceMappings["en"] = ProcessDoc(policyResourceMappingsEn)
	PolicyResourceMappings["fr"] = ProcessDoc(policyResourceMappingsFr)
}

func ProcessDoc(doc string) Doc {
	if len(doc) <= 0 {
		return Doc{}
	}
	var matter Doc
	rest, err := frontmatter.Parse(strings.NewReader(doc), &matter)
	if err != nil {
		fmt.Print(err)
		return Doc{}
	}
	matter.Long = string(rest)
	return matter
}
