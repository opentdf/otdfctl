package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/platform/protocol/go/policy"
)

type LabelUpdate struct {
	update Update
	attr   *policy.Attribute
	sdk    handlers.Handler
}

func InitLabelUpdate(label LabelItem, attr *policy.Attribute, sdk handlers.Handler) LabelUpdate {
	return LabelUpdate{
		update: InitUpdate([]string{"Key", "Value"}, []string{label.title, label.description}),
		attr:   attr,
		sdk:    sdk,
	}
}

func (m LabelUpdate) Init() tea.Cmd {
	return nil
}

func (m LabelUpdate) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			// return InitLabelList(m.attr, m.sdk)
			if m.update.focusIndex == len(m.update.inputs) {
				return InitLabelList(m.attr, m.sdk)
			}
		}
	}
	update, cmd := m.update.Update(msg)
	m.update = update.(Update)
	return m, cmd
}

func (m LabelUpdate) View() string {
	return m.update.View()
}
