package cli

import "github.com/charmbracelet/lipgloss"

func SuccessMessage(msg string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		styleSuccessStatusBar.
			MarginBottom(1).
			Render("SUCCESS"), msg)
}

func FooterMessage(msg string) string {
	if msg == "" {
		return ""
	}
	w := lipgloss.Width
	note := footerLabelStyle.Render("NOTE")
	footer := footerTextStyle.Copy().Width(defaultTableWidth - w(note)).Render(msg)
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		note,
		footer,
	)
}

func ErrorMessage(msg string, err error) string {
	if err != nil {
		msg += ": " + err.Error()
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		styleErrorStatusBar.
			PaddingRight(3).
			Render("ERROR"), msg)
}

func WarningMessage(msg string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		styleWarningStatusBar.
			Padding(0, 2).
			MarginRight(1).
			Render("WARNING"), msg)
}
