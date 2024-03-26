package man

import (
	"embed"
	"fmt"
	"strings"

	"github.com/adrg/frontmatter"
	"github.com/spf13/cobra"
)

var manLang string

//go:embed *.md
var manFiles embed.FS
var Docs Manual

type Doc struct {
	cobra.Command
}

func (d Doc) GetShort(subCmds []string) string {
	return fmt.Sprintf("%s [%s]", d.Short, strings.Join(subCmds, ", "))
}

type Manual struct {
	lang string
	Docs map[string]*Doc
	En   map[string]*Doc
	Fr   map[string]*Doc
}

func (m Manual) GetDoc(cmd string) *Doc {
	if m.lang != "en" {
		switch m.lang {
		case "fr":
			if _, ok := m.Fr[cmd]; ok {
				return m.Fr[cmd]
			}
		}
	}

	if _, ok := m.En[cmd]; !ok {
		panic(fmt.Sprintf("No doc found for cmd, %s", cmd))
	}

	return m.En[cmd]
}

func init() {
	Docs = Manual{
		Docs: make(map[string]*Doc),
		En:   make(map[string]*Doc),
		Fr:   make(map[string]*Doc),
	}

	dir, err := manFiles.ReadDir(".")
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
		if len(p) < 2 || len(p) > 3 || p[len(p)-1] != "md" {
			continue
		} else if len(p) == 3 {
			lang = p[1]
		}

		c, err := manFiles.ReadFile(f.Name())
		if err != nil {
			panic("Could not read file: " + f.Name())
		}

		d, err := ProcessDoc(string(c))
		if err != nil {
			panic(fmt.Sprintf("Could not process doc, %s: %s", f.Name(), err.Error()))
		}

		switch lang {
		case "fr":
			Docs.Fr[cmd] = d
		case "en":
			Docs.En[cmd] = d
		default:
			panic("Unknown language: " + lang)
		}
	}
}

func ProcessDoc(doc string) (*Doc, error) {
	if len(doc) <= 0 {
		return nil, fmt.Errorf("Empty document")
	}
	var matter struct {
		Use     string   `yaml:"command"`
		Aliases []string `yaml:"aliases"`
		Short   string   `yaml:"short"`
	}
	rest, err := frontmatter.Parse(strings.NewReader(doc), &matter)
	if err != nil {
		return nil, err
	}

	if matter.Use == "" {
		return nil, fmt.Errorf("required 'command' property")
	}

	long := strings.TrimSpace(string(rest))

	d := Doc{
		cobra.Command{
			Use:     matter.Use,
			Aliases: matter.Aliases,
			Short:   matter.Short,
		},
	}

	d.Long = long

	return &d, nil
}
