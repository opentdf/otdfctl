package tui

import (
	"fmt"
	"strings"

	// "log"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/opentdf/tructl/tui/constants"
)

// type (
// 	errMsg error
// )

const (
	id = iota
	name
	namespace
	rule
	description
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

type AttributeView struct {
	inputs        []interface{}
	focused       int
	err           error
	keys          []string
	title         string
	ready         bool
	viewport      viewport.Model
	width, height int
}

func SetupViewport(m AttributeView, msg tea.WindowSizeMsg) (AttributeView, []tea.Cmd) {
	var (
		cmds []tea.Cmd
	)
	headerHeight := lipgloss.Height(m.CreateHeader())
	footerHeight := lipgloss.Height(m.CreateFooter())
	verticalMarginHeight := headerHeight + footerHeight
	m.width = msg.Width
	if !m.ready {
		// Since this program is using the full size of the viewport we
		// need to wait until we've received the window dimensions before
		// we can initialize the viewport. The initial dimensions come in
		// quickly, though asynchronously, which is why we wait for them
		// here.
		m.viewport = viewport.New(msg.Width, msg.Height-verticalMarginHeight)
		m.viewport.YPosition = headerHeight
		m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
		m.ready = true

		// This is only necessary for high performance rendering, which in
		// most cases you won't need.
		//
		// Render the viewport one line below the header.
		m.viewport.YPosition = headerHeight + 1
	} else {
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - verticalMarginHeight
	}
	// var cmd tea.Cmd
	// m.viewport, cmd = m.viewport.Update(tea.KeyPgDown)
	// cmds = append(cmds, cmd)
	// m.viewport.GotoBottom()
	// cmds = append(cmds, viewport.Sync(m.viewport))
	if useHighPerformanceRenderer {
		// Render (or re-render) the whole viewport. Necessary both to
		// initialize the viewport and when the window is resized.
		//
		// This is needed for high-performance rendering only.
		cmds = append(cmds, viewport.Sync(m.viewport))
	}
	return m, cmds
}

func InitAttributeView(names []string, item AttributeItem) (tea.Model, tea.Cmd) {
	var inputs []interface{}

	ti0 := textinput.New()
	ti0.Focus()
	ti0.SetValue(item.id)
	inputs = append(inputs, ti0)

	ti1 := textinput.New()
	ti1.SetValue(item.name)
	inputs = append(inputs, ti1)

	ti2 := textinput.New()
	ti2.SetValue(item.namespace)
	inputs = append(inputs, ti2)

	ti3 := textinput.New()
	ti3.SetValue(item.rule)
	inputs = append(inputs, ti3)

	ti4 := textarea.New()
	ti4.SetValue(item.description)
	inputs = append(inputs, ti4)

	ti5 := textinput.New()
	ti5.SetValue(strings.Join(item.values, ","))
	inputs = append(inputs, ti5)

	m := AttributeView{
		keys:    names,
		inputs:  inputs,
		focused: 0,
		err:     nil,
		title:   "[Edit Attribute]",
	}
	return m.Update(WindowMsg())
}

func (m AttributeView) Init() tea.Cmd {
	return textinput.Blink
}

func (m AttributeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	item := AttributeItem{
		id:          m.inputs[id].(textinput.Model).Value(),
		name:        m.inputs[name].(textinput.Model).Value(),
		namespace:   m.inputs[namespace].(textinput.Model).Value(),
		rule:        m.inputs[rule].(textinput.Model).Value(),
		description: m.inputs[description].(textarea.Model).Value(),
		values:      strings.Split(m.inputs[values].(textinput.Model).Value(), ","),
	}
	saveModel, saveCmd := InitAttributeList([]AttributeItem{item})
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyShiftLeft: //, tea.KeyBackspace:
			return InitAttributeList([]AttributeItem{})
		case tea.KeyShiftRight:
			return saveModel, saveCmd
		case tea.KeyEnter:
			// if m.focused == len(m.inputs)-1 {
			// 	return saveModel, saveCmd
			// }
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP, tea.KeyUp:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN, tea.KeyDown:
			m.nextInput()
		}
		for i := range m.inputs {
			if i == description {
				tempInput := m.inputs[i].(textarea.Model)
				tempInput.Blur()
				m.inputs[i] = tempInput
			} else {
				tempArea := m.inputs[i].(textinput.Model)
				tempArea.Blur()
				m.inputs[i] = tempArea
			}
		}
		if m.focused == description {
			tempArea := m.inputs[m.focused].(textarea.Model)
			tempArea.Focus()
			m.inputs[m.focused] = tempArea
		} else {
			tempInput := m.inputs[m.focused].(textinput.Model)
			tempInput.Focus()
			m.inputs[m.focused] = tempInput
		}

	case tea.WindowSizeMsg:
		m, cmds = SetupViewport(m, msg)
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	var cmd tea.Cmd

	for i := range m.inputs {
		if i == description {
			m.inputs[i], cmd = m.inputs[i].(textarea.Model).Update(msg)
		} else {
			m.inputs[i], cmd = m.inputs[i].(textinput.Model).Update(msg)
		}
		cmds = append(cmds, cmd)
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func CreateEditFormat(num int) string {
	var format string
	for i := 0; i < num; i++ {
		format += "\n%s %s\n"
	}
	return format
}

func (m AttributeView) View() string {
	content := fmt.Sprintf(CreateEditFormat(len(m.inputs)),
		inputStyle.Width(len(m.keys[id])).Render(m.keys[id]),
		m.inputs[id].(textinput.Model).View(),
		inputStyle.Width(len(m.keys[name])).Render(m.keys[name]),
		m.inputs[name].(textinput.Model).View(),
		inputStyle.Width(len(m.keys[namespace])).Render(m.keys[namespace]),
		m.inputs[namespace].(textinput.Model).View(),
		inputStyle.Width(len(m.keys[rule])).Render(m.keys[rule]),
		m.inputs[rule].(textinput.Model).View(),
		inputStyle.Width(len(m.keys[description])).Render(m.keys[description]),
		m.inputs[description].(textarea.Model).View(),
		inputStyle.Width(len(m.keys[values])).Render(m.keys[values]),
		m.inputs[values].(textinput.Model).View(),
	)

	if !m.ready {
		return "\n  Initializing..."
	}
	wrapped := wordwrap.String(content, m.width)
	m.viewport.SetContent(wrapped)
	return fmt.Sprintf("%s\n%s\n%s", m.CreateHeader(), m.viewport.View(), m.CreateFooter())
}

// nextInput focuses the next input field
func (m *AttributeView) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *AttributeView) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func CreateLine(width int, text string) string {
	return strings.Repeat("─", max(0, width-lipgloss.Width(text)))
}

func (m AttributeView) CreateHeader() string {
	title := titleStyle.Render(m.title)
	line := CreateLine(m.viewport.Width, title)
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
}

func (m AttributeView) CreateFooter() string {
	info := infoStyle.Render(fmt.Sprintf("discard: shift + left arrow | save: shift + right arrow | scroll: %3.f%%", m.viewport.ScrollPercent()*100))
	line := CreateLine(m.viewport.Width, info)
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
}
