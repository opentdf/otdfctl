package cli

import "github.com/charmbracelet/lipgloss/table"

func NewTabular() *table.Table {
	t := NewTable()
	t.Headers("Property", "Value")
	return t
}
