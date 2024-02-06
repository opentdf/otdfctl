package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opentdf/tructl/tui/constants"
)

type AttributeList struct {
	list  list.Model
	width int
}

type AttributeItem struct {
	id          string
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

func InitAttributeList(items []AttributeItem) (tea.Model, tea.Cmd) {
	// TODO: fetch items from API

	m := AttributeList{}
	m.list = list.New([]list.Item{}, list.NewDefaultDelegate(), constants.WindowSize.Width, constants.WindowSize.Height)
	m.list.Title = "Attributes"
	if len(items) > 0 {
		var newItems []list.Item
		for _, item := range items {
			newItems = append(newItems, item)
		}
		m.list.SetItems(newItems)
	} else {
		// m.list.SetItems([]list.Item{
		// 	AttributeItem{
		// 		id:          "8a6755f2-efa8-4758-b893-af9a488e0bea",
		// 		namespace:   "demo.com",
		// 		name:        "relto",
		// 		rule:        "hierarchical",
		// 		description: "The relto attribute is used to describe the relationship of the resource to the country of origin.",
		// 		values:      []string{"USA", "GBR"},
		// 	},
		// })
	}
	return m.Update(WindowMsg())
}

func (m AttributeList) Init() tea.Cmd {
	return nil
}

func StyleAttr(attr string) string {
	return lipgloss.NewStyle().
		Foreground(constants.Magenta).
		Render(attr)
}

func CreateViewFormat(num int) string {
	var format string
	for i := 0; i < num; i++ {
		format += "%s %s\n"
	}
	return format
}

func (m AttributeList) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	attr_keys := []string{"Id", "Name", "Namespace", "Rule", "Description", "Values"}
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
			am.list.Select(1)
			return am.Update(WindowMsg())
		case "c":
			// show the add attribute form
			// InitAttributeCreateView()
			// return m, nil
			// attr_keys := []string{"Id", "Name", "Namespace", "Rule", "Description", "Values"}
			return InitAttributeView(attr_keys, AttributeItem{}, "[Create Attribute]")
		// case "e":
		// 	item := m.list.Items()[0].(AttributeItem)
		// 	attr_keys := []string{"Id", "Name", "Namespace", "Rule", "Description", "Values"}
		// 	return InitAttributeView(attr_keys, item)
		case "enter", "e":
			item := m.list.Items()[0].(AttributeItem)
			// attr_keys := []string{"Id", "Name", "Namespace", "Rule", "Description", "Values"}
			return InitAttributeView(attr_keys, item, "[Edit Attribute]")
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
