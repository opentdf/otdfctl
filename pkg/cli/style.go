package cli

import "github.com/charmbracelet/lipgloss"

var statusBarStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#FF5F87")).
	Padding(0, 1).
	MarginRight(1)
