package cli

import "github.com/charmbracelet/lipgloss"

var statusBarStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#FF5F87")).
	Padding(0, 1).
	MarginRight(1)

var footerLabelStyle = lipgloss.NewStyle().
	Inherit(statusBarStyle).
	Foreground(lipgloss.Color("#FFFDF5")).
	PaddingLeft(1)

var footerTextStyle = lipgloss.
	NewStyle().
	Background(lipgloss.ANSIColor(236)).
	PaddingLeft(1).
	Inherit(statusBarStyle)
