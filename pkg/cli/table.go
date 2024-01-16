package cli

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

type Table table.Table

func NewTable() *table.Table {
	t := table.New()
	return t.Border(lipgloss.NormalBorder()).
		Width(80).
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
