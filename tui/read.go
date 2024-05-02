package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/opentdf/otdfctl/tui/constants"
)

type Read struct {
	title string
	keys  []string
	vals  []string
	list  list.Model
	width int
	// list list.Model
}

// type item string
type item struct {
	title string
}

func (i item) FilterValue() string { return i.title }

func (i item) Title() string {
	return i.title
}

func InitRead(title string, items []list.Item) (tea.Model, tea.Cmd) {
	m := Read{title: title}
	m.list = list.New(items, list.NewDefaultDelegate(), constants.WindowSize.Width, constants.WindowSize.Height)
	return m.Update(WindowMsg())
}

func (m Read) Init() tea.Cmd {
	return nil
}

func (m Read) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		m.list.SetSize(msg.Width, msg.Height)
		m.width = msg.Width
		return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Read) View() string {
	// key and value
	// items := []list.Item{}
	// for i, key := range m.keys {
	// 	items = append(items, item{title: key + " > " + m.vals[i]})
	// }
	// l := list.New(items, list.NewDefaultDelegate(), constants.WindowSize.Width, constants.WindowSize.Height)
	// l.Title = m.title
	// l.SetItems(items)
	// l := list.Model{
	// 	Title: m.title,
	// 	Items: []list.Item{},
	// }
	// return m.vals[0]
	// m.list = l
	return ViewList(m.list)
	// a := ""
	// for _, i := range items {
	// 	a += i.(item).title + "\n"
	// }
	// return a
}
