package man

import (
	"fmt"
	"io/fs"
	"log/slog"
	"strings"

	"github.com/adrg/frontmatter"
	docsEmbed "github.com/opentdf/tructl/docs"
	"github.com/spf13/cobra"
)

var manLang string

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
	slog.Debug("Loading docs from embed")
	Docs = Manual{
		Docs: make(map[string]*Doc),
		En:   make(map[string]*Doc),
		Fr:   make(map[string]*Doc),
	}

	err := fs.WalkDir(docsEmbed.ManFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// extract language from filename
		p := strings.Split(d.Name(), ".")
		cmd := p[0]
		lang := "en"

		// check if file is a markdown file
		if p[len(p)-1] != "md" {
			return nil
		} else if len(p) < 2 || len(p) > 3 { // check if file complies with naming convention
			return nil
		} else if len(p) == 3 {
			lang = p[1]
		}

		// remove extension and extract command from path
		p = strings.Split(path, "/")
		// remove leading and trailing slashes
		p = p[1 : len(p)-1]
		// if the last element is not _index, it is a subcommand
		if cmd != "_index" {
			p = append(p, cmd)
		}
		cmd = strings.Join(p, "/")

		if cmd == "" {
			cmd = "<root>"
		}

		slog.Debug("Found doc", slog.String("cmd", cmd), slog.String("lang", lang))
		c, err := docsEmbed.ManFiles.ReadFile(path)
		if err != nil {
			return fmt.Errorf("Could not read file, %s: %s ", path, err.Error())
		}

		doc, err := processDoc(string(c))
		if err != nil {
			return fmt.Errorf("Could not process doc, %s: %s", path, err.Error())
		}

		slog.Debug("Adding doc: ", cmd, " ", lang, "\n")
		switch lang {
		case "fr":
			Docs.Fr[cmd] = doc
		case "en":
			Docs.En[cmd] = doc
		default:
			return fmt.Errorf("Unknown language, " + lang)
		}
		return nil
	})
	if err != nil {
		panic("Could not read embedded files: " + err.Error())
	}
}

func processDoc(doc string) (*Doc, error) {
	if len(doc) <= 0 {
		return nil, fmt.Errorf("Empty document")
	}
	var matter struct {
		Title   string `yaml:"title"`
		Command struct {
			Name    string   `yaml:"name"`
			Aliases []string `yaml:"aliases"`
			Flags   []struct {
				Name  string `yaml:"name"`
				Short string `yaml:"short"`
			} `yaml:"flags"`
		} `yaml:"command"`
	}
	rest, err := frontmatter.Parse(strings.NewReader(doc), &matter)
	if err != nil {
		return nil, err
	}

	c := matter.Command

	if c.Name == "" {
		return nil, fmt.Errorf("required 'command' property")
	}

	long := "# " + matter.Title + "\n\n"
	long += strings.TrimSpace(string(rest))

	d := Doc{
		cobra.Command{
			Use:     c.Name,
			Aliases: c.Aliases,
			Short:   matter.Title,
		},
	}

	d.Long = long

	return &d, nil
}
