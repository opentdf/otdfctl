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

// type AttributeItem struct {
// 	id          menuState
// 	title       string
// 	description string
// }

// func (m AttributeItem) FilterValue() string {
// 	return m.title
// }

// func (m AttributeItem) Title() string {
// 	return m.title
// }

// func (m AttributeItem) Description() string {
// 	return m.description
// }

type AttributeView struct {
	list list.Model
	view tea.Model
	sdk  handlers.Handler
}

func InitAttributeView(id string, h handlers.Handler) (AttributeView, tea.Cmd) {
	m := AttributeView{
		view: nil,
		sdk:  h,
	}
	attr, err := h.GetAttribute(id)
	if err != nil {
		return m, nil
	}
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 8, 8)
	m.list.Title = "OpenTDF"
	m.list.SetItems([]list.Item{
		// AppMenuItem{title: "Namespaces", description: "Manage namespaces", id: namespaceMenu},
		AttributeItem{title: "ID", description: "ID >" + attr.Id, _id: attr.Id, id: 0, name: "OK"},
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
		m.list.SetSize(msg.Width, msg.Height)
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
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m AttributeView) View() string {
	return ViewList(m.list)
}
