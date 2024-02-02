package cli

import "github.com/charmbracelet/lipgloss"

func SuccessMessage(msg string) string {
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		statusBarStyle.Background(lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}).
			MarginBottom(1).
			Render("SUCCESS"), msg)
}

func FooterMessage(msg string) string {
	w := lipgloss.Width
	note := footerLabelStyle.Render("NOTE ")
	footer := footerTextStyle.Copy().Width(defaultTableWidth - w(note)).Render(msg)
	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		note,
		footer,
	)
}

func ErrorMessage(msg string, err error) string {
	if err != nil {
		msg = ": " + err.Error()
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Left,
		statusBarStyle.Background(lipgloss.AdaptiveColor{Light: "#FF5F87", Dark: "#FF5F87"}).
			PaddingRight(3).
			Render("ERROR"), msg)
}
