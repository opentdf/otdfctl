package tui

import tea "github.com/charmbracelet/bubbletea"

type LabelUpdate struct {
	update Update
}

func InitLabelUpdate(label LabelItem) LabelUpdate {
	return LabelUpdate{
		update: InitUpdate([]string{"Key", "Value"}, []string{label.title, label.description}),
	}
}

func (m LabelUpdate) Init() tea.Cmd {
	return nil
}

func (m LabelUpdate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.update.Update(msg)
}

func (m LabelUpdate) View() string {
	return m.update.View()
}
