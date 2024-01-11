package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opentdf/tructl/tui/constants"
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

func InitAppMenu() (tea.Model, tea.Cmd) {
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
				attributeView := InitAttributeView()
				am, cmd := attributeView.Update(constants.WindowSize)
				m.view = am
				return m, cmd
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m AppMenu) View() string {
	// return m.list.View()
	// create a new view with a list view as the main view
	lipgloss.NewStyle().Padding(1, 2, 1, 2)
	return lipgloss.JoinVertical(lipgloss.Top, m.list.View())
}
