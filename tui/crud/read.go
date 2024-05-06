package crud

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/opentdf/otdfctl/tui/common"
	"github.com/opentdf/otdfctl/tui/constants"
)

type Read struct {
	list  list.Model
	width int
}

func InitRead(title string, items []list.Item) (tea.Model, tea.Cmd) {
	m := Read{}
	m.list = list.New(items, list.NewDefaultDelegate(), constants.WindowSize.Width, constants.WindowSize.Height)
	m.list.Title = title
	return m.Update(common.WindowMsg())
}

func (m Read) Init() tea.Cmd {
	return nil
}

func (m Read) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	constants.WindowSize = msg
	// 	m.list.SetSize(msg.Width, msg.Height)
	// 	m.width = msg.Width
	// 	return m, nil
	// case tea.KeyMsg:
	// 	switch msg.Type {
	// 	case tea.KeyCtrlC, tea.KeyEsc:
	// 		return m, tea.Quit
	// 	}
	// }
	return m, nil
}

func (m Read) View() string {
	return common.ViewList(m.list)
}
