package man

import (
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/glamour/ansi"
	"golang.org/x/term"
)

func styleDoc(doc string) string {
	w, _, err := term.GetSize(0)
	if err != nil {
		w = 80
	}
	if w > 120 {
		w = 120
	}
	// Set up a new glamour instance
	// with some options
	ds := glamour.DarkStyleConfig
	// ls := glamour.DefaultStyles["light"]

	ds.Document.Margin = uintPtr(0)
	ds.Paragraph.Margin = uintPtr(2)
	// Capitalize headers
	ds.H1.StylePrimitive = ansi.StylePrimitive{
		Upper:  boolPtr(true),
		Color:  stringPtr("#F1F1F1"),
		Format: "# {{.text}}",
	}
	r, _ := glamour.NewTermRenderer(
		// glamour.WithAutoStyle(),
		glamour.WithStyles(ds),
		glamour.WithWordWrap(w),
	)

	// Render the content
	out, _ := r.Render(doc)

	return out
}

func boolPtr(b bool) *bool       { return &b }
func stringPtr(s string) *string { return &s }
func uintPtr(u uint) *uint       { return &u }
