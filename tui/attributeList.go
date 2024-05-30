package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/tui/constants"
	"github.com/opentdf/platform/protocol/go/common"
)

type AttributeList struct {
	list list.Model
	sdk  handlers.Handler
}

type AttributeItem struct {
	id   string
	name string
}

func (m AttributeItem) FilterValue() string {
	return m.name
}

func (m AttributeItem) Title() string {
	return m.name
}

func (m AttributeItem) Description() string {
	return m.id
}

func InitAttributeList(id string, sdk handlers.Handler) (tea.Model, tea.Cmd) {
	l := list.New([]list.Item{}, list.NewDefaultDelegate(), constants.WindowSize.Width, constants.WindowSize.Height)
	res, err := sdk.ListAttributes(common.ActiveStateEnum_ACTIVE_STATE_ENUM_ANY)
	if err != nil {
		// return error view
	}
	var attrs []list.Item
	selectIdx := 0
	for i, attr := range res {
		var vals []string
		for _, val := range attr.Values {
			vals = append(vals, val.Value)
		}
		if attr.Id == id {
			selectIdx = i
		}
		item := AttributeItem{
			id:   attr.Id,
			name: attr.Name,
		}
		attrs = append(attrs, item)
	}
	l.Title = "Attributes"
	l.SetItems(attrs)
	l.Select(selectIdx)
	m := AttributeList{sdk: sdk, list: l}
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
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "ctrl+[", "backspace":
			am, _ := InitAppMenu(m.sdk)
			// make enum for Attributes idx in AppMenu
			am.list.Select(0)
			return am.Update(WindowMsg())
		// case "c":
		// create new attribute
		// return InitAttributeView(m.list.Items(), len(m.list.Items()))
		case "enter", "e":
			return InitAttributeView(m.list.Items()[m.list.Index()].(AttributeItem).id, m.sdk)
			// case "ctrl+d":
			// 	m.list.RemoveItem(m.list.Index())
			// 	newIndex := m.list.Index() - 1
			// 	if newIndex < 0 {
			// 		newIndex = 0
			// 	}
			// 	m.list.Select(newIndex)
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
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
