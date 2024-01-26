package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/opentdf/tructl/tui/constants"
)

type AttributeList struct {
	list  list.Model
	width int
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

func (m AttributeItem) Title() string {
	return m.name
}

func (m AttributeItem) Description() string {
	return m.description
}

func InitAttributeList() AttributeList {
	// TODO: fetch items from API

	m := AttributeList{}
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), constants.WindowSize.Width, constants.WindowSize.Height)
	m.list.Title = "Attributes"
	m.list.SetItems([]list.Item{
		AttributeItem{
			id:          1,
			namespace:   "demo.com",
			name:        "relto",
			rule:        "hierarchical",
			description: "The relto attribute is used to describe the relationship of the resource to the country of origin.",
			values:      []string{"USA", "GBR"},
		},
	})

	return m
}

func (m AttributeList) Init() tea.Cmd {
	return nil
}

func StyleAttr(attr string) string {
	return lipgloss.NewStyle().
		Foreground(constants.Magenta).
		Render(attr)
}

func CreateFormat(num int) string {
	var format string
	for i := 0; i < num; i++ {
		format += "%s %s\n"
	}
	return format
}

func (m AttributeList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.list.SetSize(msg.Width, msg.Height)
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+[", "backspace":
			am, _ := InitAppMenu()
			return am.Update(tea.WindowSizeMsg{Width: constants.WindowSize.Width, Height: constants.WindowSize.Height})
			// return am, cmd
		case "c":
			// show the add attribute form
			// InitAttributeCreateView()
			return m, nil
		case "enter":
			item := m.list.Items()[0].(AttributeItem)
			attr_keys := []string{"Name", "Namespace", "Rule", "Description", "Values"}
			content := fmt.Sprintf(
				CreateFormat(len(attr_keys)),
				StyleAttr(attr_keys[0]), item.name,
				StyleAttr(attr_keys[1]), item.namespace,
				StyleAttr(attr_keys[2]), item.rule,
				StyleAttr(attr_keys[3]), item.description,
				StyleAttr(attr_keys[4]), item.values,
			)
			wrapped := wordwrap.String(content, m.width)
			am := AttributeView{}
			am.title = "Attribute"
			am.content = wrapped
			return am.Update(tea.WindowSizeMsg{Width: constants.WindowSize.Width, Height: constants.WindowSize.Height})
		}
	}
	return m, nil
}

func (m AttributeList) View() string {
	return ViewList(m.list)
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
