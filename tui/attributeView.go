package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/opentdf/tructl/tui/constants"
	"github.com/charmbracelet/lipgloss"
)

type AttributeModel struct {
	list list.Model
}

type AttributeItem struct {
	id          int
	title string
	namespace   string
	name        string
	description string
	rule        string
	values      []string
}

func (m AttributeItem) FilterValue() string {
	return m.name
}
// func (m AttributeItem) FilterValue() string {
// 	return m.title
// }

func (m AttributeItem) Title() string {
	return m.name
}

func (m AttributeItem) Description() string {
	return m.description
}
// type AppMenuItem struct {
// 	id          menuState
// 	title       string
// 	description string
// }

// func (m AppMenuItem) FilterValue() string {
// 	return m.title
// }

// func (m AppMenuItem) Title() string {
// 	return m.title
// }

// func (m AppMenuItem) Description() string {
// 	return m.description
// }

func InitAttributeView() AttributeModel {
	// TODO: fetch items from API

	m := AttributeModel{}
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), constants.WindowSize.Width, constants.WindowSize.Height)
	m.list.Title = "Attributes"
	m.list.SetItems([]list.Item{
		AttributeItem{
			id:          1,
			namespace:   "demo.com",
			name:        "relto",
			// title: 	 "Relationship To",
			rule:        "heirarchical",
			description: "The relto attribute is used to describe the relationship of the resource to the country of origin.",
			values:      []string{"USA", "GBR"},
		},
		// AttributeItem{
		// 	id:          1,
		// 	namespace:   "demo.com",
		// 	name:        "relto",
		// 	title: 	 "Relationship To",
		// 	rule:        "heirarchical",
		// 	description: "The relto attribute is used to describe the relationship of the resource to the country of origin.",
		// 	values:      []string{"USA", "GBR"},
		// },
	})
	// m.list.SetItems([]list.Item{
	// 	AppMenuItem{title: "relto", description: "The relto attribute is used to describe the relationship of the resource to the country of origin.", id: 1},
	// })
	return m
}

func (m AttributeModel) Init() tea.Cmd {
	return nil
}

func (m AttributeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+[":
			return InitAppMenu()
		case "c":
			// show the add attribute form
			// InitAttributeCreateView()
			return m, nil
		// case "enter":
		// 	return m, View()
		}
	}
	return m, nil
}

func (m AttributeModel) View() string {
	// return m.list.View()
	lipgloss.NewStyle().Padding(1, 2, 1, 2)
	return lipgloss.JoinVertical(lipgloss.Top, m.list.View())
}

// func AddAttribute() {
// 	var namespace string

// 	form := huh.NewForm(
// 		huh.NewGroup(
// 			huh.NewSelect[string]().
// 				Title("Namespace").
// 				Options(
// 					huh.NewOption("demo.com", "demo.com"),
// 					huh.NewOption("demo.net", "demo.net"),
// 				).
// 				Validate(func(str string) error {
// 					// Check if namespace exists
// 					fmt.Println(str)
// 					return nil
// 				}).
// 				Value(&namespace),
// 		),
// 	)

// 	if err := form.Run(); err != nil {
// 		return
// 	}
// }
