package cli

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
)

type Table table.Table

var defaultTableWidth int

func NewTable() *table.Table {
	t := table.New()
	return t.Border(lipgloss.NormalBorder()).
		Width(defaultTableWidth).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(25))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == 0:
				return lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(25)).Bold(true)
			case row%2 == 0:
				// bright grey text on default background
				return lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(236))
			default:
				// odd rows: bright grey text on dark grey background
				return lipgloss.NewStyle().Foreground(lipgloss.ANSIColor(250)).Background(lipgloss.Color("#2B2B2B"))
			}
		})
}

func init() {
	// dynamically set the default table width based on terminal size breakpoints
	w, _, err := term.GetSize(0)
	if err != nil {
		w = 80
	}
	if w > 180 {
		defaultTableWidth = 180
	} else if w > 120 {
		defaultTableWidth = 120
	} else {
		defaultTableWidth = 80
	}
}
