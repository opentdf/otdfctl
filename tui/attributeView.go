package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/tui/constants"
)

type AttributeSubItem struct {
	title       string
	description string
}

func (m AttributeSubItem) FilterValue() string {
	return m.title
}

func (m AttributeSubItem) Title() string {
	return m.title
}

func (m AttributeSubItem) Description() string {
	return m.description
}

type AttributeView struct {
	read Read
}

func InitAttributeView(id string, h handlers.Handler) (AttributeView, tea.Cmd) {
	m := AttributeView{}
	attr, err := h.GetAttribute(id)
	if err != nil {
		return m, nil
	}
	var vals []string
	for _, val := range attr.Values {
		vals = append(vals, val.Value)
	}
	items := []list.Item{
		AttributeSubItem{title: "ID", description: attr.Id},
		AttributeSubItem{title: "Name", description: attr.Name},
		AttributeSubItem{title: "Rule", description: attr.Rule.String()},
		AttributeSubItem{title: "Values", description: cli.CommaSeparated(vals)},
		AttributeSubItem{title: "Namespace", description: attr.Namespace.Name},
		AttributeSubItem{title: "Created At", description: attr.Metadata.CreatedAt.String()},
		AttributeSubItem{title: "Updated At", description: attr.Metadata.UpdatedAt.String()},
	}

	model, _ := InitRead("Read Attribute", items)
	m.read = model.(Read)
	model, msg := m.Update(WindowMsg())
	m = model.(AttributeView)
	return m, msg
}

func (m AttributeView) Init() tea.Cmd {
	return nil
}

func (m AttributeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.read.list.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+d":
			return m, nil
			// case "enter":
			// 	switch m.list.SelectedItem().(AttributeItem).id {
			// 	// case namespaceMenu:
			// 	// 	// get namespaces
			// 	// 	nl, cmd := InitNamespaceList([]list.Item{}, 0)
			// 	// 	return nl, cmd
			// 	case attributeMenu:
			// 		// list attributes
			// 		al, cmd := InitAttributeList("", m.sdk)
			// 		return al, cmd
			// 	}
		}
	}

	var cmd tea.Cmd
	m.read.list, cmd = m.read.list.Update(msg)
	return m, cmd
}

func (m AttributeView) View() string {
	return m.read.View()
}
