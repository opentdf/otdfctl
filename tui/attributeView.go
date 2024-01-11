package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/opentdf/tructl/tui/constants"
)

type AttributeModel struct {
	list list.Model
}

type AttributeItem struct {
	id          int
	namespace   string
	name        string
	description string
	rule        string
	values      []string
}

func (m AttributeItem) FilterValue() string {
	return m.name
}

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
			rule:        "heirarchical",
			description: "The relto attribute is used to describe the relationship of the resource to the country of origin.",
			values:      []string{"USA", "GBR"},
		},
	})
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
		}
	}
	return m, nil
}

func (m AttributeModel) View() string {
	return m.list.View()
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
