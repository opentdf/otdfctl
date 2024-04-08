package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/opentdf/otdfctl/tui/constants"
)

const (
	mainMenu menuState = iota
	namespaceMenu
	attributeMenu
	entitlementMenu
	resourceEncodingMenu
	subjectEncodingMenu
)

type menuState int

type AppMenuItem struct {
	id          menuState
	title       string
	description string
}

func (m AppMenuItem) FilterValue() string {
	return m.title
}

func (m AppMenuItem) Title() string {
	return m.title
}

func (m AppMenuItem) Description() string {
	return m.description
}

type AppMenu struct {
	list list.Model
	view tea.Model
}

func InitAppMenu() (AppMenu, tea.Cmd) {
	m := AppMenu{
		view: nil,
	}
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), 8, 8)
	m.list.Title = "OpenTDF"
	m.list.SetItems([]list.Item{
		AppMenuItem{title: "Namespaces", description: "Manage namespaces", id: namespaceMenu},
		AppMenuItem{title: "Attributes", description: "Manage attributes", id: attributeMenu},
		AppMenuItem{title: "Entitlements", description: "Manage entitlements", id: entitlementMenu},
		AppMenuItem{title: "Resource Encodings", description: "Manage resource encodings", id: resourceEncodingMenu},
		AppMenuItem{title: "Subject Encodings", description: "Manage subject encodings", id: subjectEncodingMenu},
	})
	return m, func() tea.Msg { return nil }
}

func (m AppMenu) Init() tea.Cmd {
	return nil
}

func (m AppMenu) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
		case "enter":
			switch m.list.SelectedItem().(AppMenuItem).id {
			case attributeMenu:
				item := AttributeItem{
					id:          "8a6755f2-efa8-4758-b893-af9a488e0bea",
					namespace:   "demo.com",
					name:        "relto",
					rule:        "hierarchical",
					description: "The relto attribute is used to describe the relationship of the resource to the country of origin.",
					values:      []string{"USA", "GBR"},
				}
				al, cmd := InitAttributeList([]list.Item{item}, 0)
				return al, cmd
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m AppMenu) View() string {
	return ViewList(m.list)
}
