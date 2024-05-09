package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/tui/constants"
	"github.com/opentdf/platform/protocol/go/policy"
)

// func (m LabelList)

type LabelList struct {
	attr *policy.Attribute
	sdk  handlers.Handler
	read Read
}

type LabelItem struct {
	title       string
	description string
}

func (m LabelItem) FilterValue() string {
	return m.title
}

func (m LabelItem) Title() string {
	return m.title
}

func (m LabelItem) Description() string {
	return m.description
}

func InitLabelList(attr *policy.Attribute, sdk handlers.Handler) (tea.Model, tea.Cmd) {
	labels := attr.Metadata.Labels
	var items []list.Item
	for k, v := range labels {
		item := LabelItem{
			title:       k,
			description: v,
		}
		items = append(items, item)
	}
	model, _ := InitRead("Read Labels", items)
	return LabelList{attr: attr, sdk: sdk, read: model.(Read)}, nil
}

func (m LabelList) Init() tea.Cmd {
	return nil
}

func (m LabelList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.read.list.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "backspace":
			return InitAttributeView(m.attr.Id, m.sdk)
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter", "e":
			// edit?
			return m, nil
		case "c":
			// create new label
			return m, nil
		case "d":
			// delete label
			return m, nil
		}
	}
	return m, nil
}

func (m LabelList) View() string {
	return ViewList(m.read.list)
}
