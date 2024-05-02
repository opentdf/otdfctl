package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/tui/constants"
)

// const (
// 	mainMenu menuState = iota
// 	namespaceMenu
// 	attributeMenu
// 	entitlementMenu
// 	resourceEncodingMenu
// 	subjectEncodingMenu
// )

// type menuState int

type AttributeSubItem struct {
	title       string
	description string
}

// type AttributeItem struct {
// 	id          menuState
// 	title       string
// 	description string
// }

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
	// list list.Model
	// sdk  handlers.Handler
}

func InitAttributeView(id string, h handlers.Handler) (AttributeView, tea.Cmd) {
	read := Read{title: "Read Attribute"}
	m := AttributeView{read: read}
	attr, err := h.GetAttribute(id)
	if err != nil {
		return m, nil
	}
	model, _ := InitRead("Read Attribute", []list.Item{
		// AppMenuItem{title: "Namespaces", description: "Manage namespaces", id: namespaceMenu},
		AttributeSubItem{title: "ID", description: attr.Id},
		// AppMenuItem{title: "Entitlements", description: "Manage entitlements", id: entitlementMenu},
		// AppMenuItem{title: "Resource Encodings", description: "Manage resource encodings", id: resourceEncodingMenu},
		// AppMenuItem{title: "Subject Encodings", description: "Manage subject encodings", id: subjectEncodingMenu},
	})
	m.read = model.(Read)

	m.read.list = list.New([]list.Item{}, list.NewDefaultDelegate(), constants.WindowSize.Width, constants.WindowSize.Height)
	m.read.list.Title = "Read Attribute"
	m.read.list.SetItems([]list.Item{
		// AppMenuItem{title: "Namespaces", description: "Manage namespaces", id: namespaceMenu},
		AttributeSubItem{title: "ID", description: attr.Id},
		// AppMenuItem{title: "Entitlements", description: "Manage entitlements", id: entitlementMenu},
		// AppMenuItem{title: "Resource Encodings", description: "Manage resource encodings", id: resourceEncodingMenu},
		// AppMenuItem{title: "Subject Encodings", description: "Manage subject encodings", id: subjectEncodingMenu},
	})

	model, msg := m.Update(WindowMsg())
	m = model.(AttributeView)
	return m, msg
	// return m, func() tea.Msg { return nil }
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
