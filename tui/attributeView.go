package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/opentdf/otdfctl/pkg/cli"
	"github.com/opentdf/otdfctl/pkg/handlers"
	"github.com/opentdf/otdfctl/tui/constants"
)

const (
	id = iota
	name
	namespace
	rule
	values
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
	cyan     = lipgloss.Color("#00FFFF")
)

const useHighPerformanceRenderer = false

var (
	inputStyle    = lipgloss.NewStyle().Foreground(constants.Magenta)
	continueStyle = lipgloss.NewStyle().Foreground(cyan)
)

var (
	titleStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Right = "├"
		return lipgloss.NewStyle().BorderStyle(b).Padding(0, 1)
	}()

	infoStyle = func() lipgloss.Style {
		b := lipgloss.RoundedBorder()
		b.Left = "┤"
		return titleStyle.Copy().BorderStyle(b)
	}()
)

// type TextWrapper struct{}

// func View(m TextWrapper)  {}
// func Value(m TextWrapper) {}

type AttributeView struct {
	// inputs        []interface{}
	// focused       int
	// err           error
	// keys          []string
	// title         string
	// ready         bool
	// viewport      viewport.Model
	// width, height int
	// list          []list.Item
	// idx           int
	// editMode      bool
	// sdk           handlers.Handler
	list  list.Model
	read  Read
	width int
}

// func SetupViewport(m AttributeView, msg tea.WindowSizeMsg) (AttributeView, []tea.Cmd) {
// 	var cmds []tea.Cmd
// 	headerHeight := lipgloss.Height(m.CreateHeader())
// 	footerHeight := lipgloss.Height(m.CreateFooter())
// 	verticalMarginHeight := headerHeight + footerHeight
// 	m.width = msg.Width

// 	if !m.ready {
// 		// Since this program is using the full size of the viewport we
// 		// need to wait until we've received the window dimensions before
// 		// we can initialize the viewport. The initial dimensions come in
// 		// quickly, though asynchronously, which is why we wait for them
// 		// here.
// 		m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
// 		m.viewport.YPosition = headerHeight
// 		m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
// 		m.ready = true

// 		// This is only necessary for high performance rendering, which in
// 		// most cases you won't need.
// 		//
// 		// Render the viewport one line below the header.
// 		m.viewport.YPosition = headerHeight + 1
// 	} else {
// 		m.viewport.Width = msg.Width
// 		m.viewport.Height = msg.Height - verticalMarginHeight
// 	}

// 	if useHighPerformanceRenderer {
// 		// Render (or re-render) the whole viewport. Necessary both to
// 		// initialize the viewport and when the window is resized.
// 		//
// 		// This is needed for high-performance rendering only.
// 		cmds = append(cmds, viewport.Sync(m.viewport))
// 	}
// 	return m, cmds
// }

func InitAttributeView(id string, sdk handlers.Handler) (tea.Model, tea.Cmd) {
	read, _ := InitRead("Read Attribute")
	m := AttributeView{
		// title:    "Attribute",
		// keys:     []string{"ID", "Name", "Namespace", "Rule", "Values"},
		// inputs:   make([]interface{}, 5),
		// focused:  0,
		// viewport: viewport.Model{},
		// width:    80,
		// height:   20,
		// sdk:      sdk,
		read: read.(Read),
	}
	attr, err := sdk.GetAttribute(id)
	if err != nil {
		return m, nil
	}
	var vals []string
	for _, val := range attr.Values {
		vals = append(vals, val.Value)
	}
	// m.read.title = "Read Attribute"
	m.read.keys = []string{"Id", "Name", "Rule", "Values", "Namespace", "Labels", "Created At", "Updated At"}
	m.read.vals = []string{attr.Id, attr.Name, attr.Rule.String(), cli.CommaSeparated(vals), attr.Namespace.Name, "TODO", attr.Metadata.CreatedAt.String(), attr.Metadata.UpdatedAt.String()}
	return m, nil
}

func (m AttributeView) Init() tea.Cmd {
	return nil
}

func (m AttributeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// case tea.WindowSizeMsg:
	// 	constants.WindowSize = msg
	// 	m.list.SetSize(msg.Width, msg.Height)
	// 	m.width = msg.Width
	// 	return m, nil
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m AttributeView) View() string {
	return m.read.View()
}

// func (m AttributeView) IsNew() bool {
// 	return m.idx >= len(m.list)
// }

// func (m AttributeView) ChangeMode() AttributeView {
// 	m.editMode = m.IsNew() || !m.editMode
// 	return m
// }

// func (m AttributeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
// 	var cmds []tea.Cmd // = make([]tea.Cmd, len(m.inputs))
// 	var editing bool
// 	item := AttributeItem{
// 		id:        m.inputs[id].(textinput.Model).Value(),
// 		name:      m.inputs[name].(textinput.Model).Value(),
// 		namespace: m.inputs[namespace].(textinput.Model).Value(),
// 		rule:      m.inputs[rule].(textinput.Model).Value(),
// 		values:    strings.Split(m.inputs[values].(textinput.Model).Value(), ","),
// 	}
// 	switch msg := msg.(type) {
// 	case tea.KeyMsg:
// 		switch msg.Type {
// 		case tea.KeyShiftLeft:
// 			listIdx := m.idx
// 			if m.IsNew() {
// 				listIdx -= 1
// 			}
// 			return InitAttributeList(item.id, m.sdk)
// 		case tea.KeyShiftRight:
// 			if !m.IsNew() {
// 				// edit
// 				m.list[m.idx] = list.Item(item)
// 			} else {
// 				// create
// 				m.list = append(m.list, list.Item(item))
// 			}

// 			return InitAttributeList(item.id, m.sdk)
// 		case tea.KeyEnter:
// 			m.nextInput()
// 		case tea.KeyCtrlC, tea.KeyEsc:
// 			if m.editMode {
// 				m.editMode = false
// 			} else {
// 				return m, tea.Quit
// 			}
// 		case tea.KeyShiftTab, tea.KeyCtrlP, tea.KeyUp:
// 			m.prevInput()
// 		case tea.KeyTab, tea.KeyCtrlN, tea.KeyDown:
// 			m.nextInput()
// 		}
// 		if msg.String() == "i" && !m.editMode {
// 			editing = true
// 			m = m.ChangeMode()
// 			var cmd tea.Cmd
// 			tempInput := m.inputs[m.focused].(textinput.Model)
// 			cmd = tempInput.Cursor.SetMode(0)
// 			m.inputs[m.focused] = tempInput
// 			return m, cmd
// 		}
// 		for i := range m.inputs {

// 			tempArea := m.inputs[i].(textinput.Model)
// 			tempArea.Blur()
// 			m.inputs[i] = tempArea

// 		}

// 		tempInput := m.inputs[m.focused].(textinput.Model)
// 		tempInput.Focus()
// 		m.inputs[m.focused] = tempInput

// 	case tea.WindowSizeMsg:
// 		m, cmds = SetupViewport(m, msg)
// 	// We handle errors just like any other message
// 	case errMsg:
// 		m.err = msg
// 		return m, nil
// 	}

// 	var cmd tea.Cmd
// 	if m.editMode || m.IsNew() && !editing {
// 		for i := range m.inputs {

// 			m.inputs[i], cmd = m.inputs[i].(textinput.Model).Update(msg)

// 			cmds = append(cmds, cmd)
// 		}
// 	}
// 	m.viewport, cmd = m.viewport.Update(msg)
// 	cmds = append(cmds, cmd)
// 	return m, tea.Batch(cmds...)
// }

// func CreateEditFormat(num int) string {
// 	var format string
// 	prefix := "\n%s"
// 	postfix := "%s\n"
// 	var middle string
// 	for i := 0; i < num; i++ {
// 		format += prefix + middle + postfix
// 	}
// 	return format
// }

// func (m AttributeView) View() string {
// 	content := fmt.Sprintf(CreateEditFormat(len(m.inputs)),
// 		inputStyle.Width(len(m.keys[id])).Render(m.keys[id]),
// 		m.inputs[id].(textinput.Model).View(),
// 		inputStyle.Width(len(m.keys[name])).Render(m.keys[name]),
// 		m.inputs[name].(textinput.Model).View(),
// 		inputStyle.Width(len(m.keys[namespace])).Render(m.keys[namespace]),
// 		m.inputs[namespace].(textinput.Model).View(),
// 		inputStyle.Width(len(m.keys[rule])).Render(m.keys[rule]),
// 		m.inputs[rule].(textinput.Model).View(),
// 		inputStyle.Width(len(m.keys[values])).Render(m.keys[values]),
// 		m.inputs[values].(textinput.Model).View(),
// 	)

// 	if !m.ready {
// 		return "\n  Initializing..."
// 	}
// 	wrapped := wordwrap.String(content, m.width)
// 	m.viewport.SetContent(wrapped)
// 	return fmt.Sprintf("%s\n%s\n%s", m.CreateHeader(), m.viewport.View(), m.CreateFooter())
// }

// // nextInput focuses the next input field
// func (m *AttributeView) nextInput() {
// 	m.focused = (m.focused + 1) % len(m.inputs)
// }

// // prevInput focuses the previous input field
// func (m *AttributeView) prevInput() {
// 	m.focused--
// 	// Wrap around
// 	if m.focused < 0 {
// 		m.focused = len(m.inputs) - 1
// 	}
// }

// func CreateLine(width int, text string) string {
// 	return strings.Repeat("─", max(0, width-lipgloss.Width(text)))
// }

// func (m AttributeView) CreateHeader() string {
// 	title := titleStyle.Render(m.title)
// 	line := CreateLine(m.viewport.Width, title)
// 	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
// }

// func (m AttributeView) CreateFooter() string {
// 	var prefix string
// 	if m.editMode || m.IsNew() {
// 		prefix = "discard: shift + left arrow | save: shift + right arrow"
// 	} else {
// 		prefix = "enter edit mode: i | go back: shift + left arrow"
// 	}
// 	info := infoStyle.Render(fmt.Sprintf(prefix+" | scroll: %3.f%%", m.viewport.ScrollPercent()*100))
// 	line := CreateLine(m.viewport.Width, info)
// 	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
// }
