package man

import (
	"embed"
	"fmt"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/spf13/cobra"
)

//go:embed *.md
var e embed.FS
var Docs Man

type Man struct {
	Lang string
	Docs struct {
		En map[string]*cobra.Command
		Fr map[string]*cobra.Command
	}
}

func (m Man) GetDoc(key string) *cobra.Command {
	l := "en"
	if m.Lang != "" {
		l = m.Lang
	}

	if l != "en" {
		switch l {
		case "fr":
			if _, ok := m.Docs.Fr[key]; ok {
				return m.Docs.Fr[key]
			}
		}
	}

	if _, ok := m.Docs.En[key]; !ok {
		panic(fmt.Sprintf("No doc found for key %s", key))
	}

	return m.Docs.En[key]
}

var PolicyResourceMappings = map[string]Doc{}

var PolicyAttributes = map[string]Doc{}
var PolicyAttributeValuesCreate = Doc{}

type Doc struct {
	cobra.Command
}

func (m Doc) ShortWithSubCommands(subcommands []string) string {
	return m.Short + " [" + strings.Join(subcommands, ", ") + "]"
}

func init() {
	Docs = Man{}
	Docs.Docs.En = make(map[string]*cobra.Command)
	Docs.Docs.Fr = make(map[string]*cobra.Command)

	dir, err := e.ReadDir(".")
	if err != nil {
		panic("Could not read embedded files")
	}

	for _, f := range dir {
		if f.IsDir() {
			continue
		}

		// extract language from filename
		p := strings.Split(f.Name(), ".")
		cmd := p[0]
		lang := "en"

		// ignore files that are not markdown or that have more than one extension
		if len(p) < 1 || len(p) > 3 || p[len(p)-1] != "md" {
			continue
		} else if len(p) == 3 {
			lang = p[1]
		}

		c, err := e.ReadFile(f.Name())
		if err != nil {
			panic("Could not read file: " + f.Name())
		}

		d := ProcessDoc(string(c))

		switch lang {
		case "fr":
			Docs.Docs.Fr[cmd] = d
		case "en":
		default:
			panic("Unknown language: " + lang)
		}
	}
}

func ProcessDoc(doc string) *cobra.Command {
	if len(doc) <= 0 {
		return &cobra.Command{}
	}
	var matter struct {
		Use     string   `yaml:"command"`
		Aliases []string `yaml:"aliases"`
		Short   string   `yaml:"short"`
	}
	rest, err := frontmatter.Parse(strings.NewReader(doc), &matter)
	if err != nil {
		fmt.Print(err)
		return &cobra.Command{}
	}

	return &cobra.Command{
		Use:     matter.Use,
		Aliases: matter.Aliases,
		Short:   matter.Short,
		Long:    string(rest),
	}
}
