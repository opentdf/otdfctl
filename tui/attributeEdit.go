package tui

import (
	"fmt"
	// "log"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
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
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

type model struct {
	inputs  []textinput.Model
	focused int
	err     error
	keys    []string
	// width   int
}

func InitAttributeEdit(names []string, item AttributeItem) (tea.Model, tea.Cmd) {
	var inputs []textinput.Model = make([]textinput.Model, 6)
	inputs[id] = textinput.New()
	inputs[id].Placeholder = "4505 **** **** 1234"
	inputs[id].Focus()
	// inputs[ccn].CharLimit = 20
	// inputs[id].Width = 30
	inputs[id].Prompt = ""

	inputs[name] = textinput.New()
	inputs[name].Placeholder = "MM/YY "
	// inputs[exp].CharLimit = 5
	// inputs[name].Width = 5
	inputs[name].Prompt = ""

	inputs[namespace] = textinput.New()
	inputs[namespace].Placeholder = "XXX"
	// inputs[cvv].CharLimit = 3
	// inputs[namespace].Width = 5
	inputs[namespace].Prompt = ""
	m := model{
		keys:    names,
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
	return m.Update(WindowMsg())
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyShiftLeft: //, tea.KeyBackspace:
			return InitAttributeList()
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				return m, tea.Quit
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	case tea.WindowSizeMsg:
		for i := range m.inputs {
			m.inputs[i].Width = msg.Width
		}
		// m.width = msg.Width
		return m, nil
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func CreateEditFormat(num int) string {
	var format string
	for i := 0; i < num; i++ {
		format += "\n%s\n%s\n"
	}
	return format
}

func (m model) View() string {
	content := fmt.Sprintf("Edit Attribute\n"+CreateEditFormat(len(m.inputs))+"\n%s",
		inputStyle.Width(len(m.keys[id])).Render(m.keys[id]),
		m.inputs[id].View(),
		inputStyle.Width(len(m.keys[name])).Render(m.keys[name]),
		wordwrap.String(m.inputs[name].View(), 15),
		inputStyle.Width(len(m.keys[namespace])).Render(m.keys[namespace]),
		m.inputs[namespace].View(),
		inputStyle.Width(len(m.keys[rule])).Render(m.keys[rule]),
		m.inputs[rule].View(),
		inputStyle.Width(len(m.keys[description])).Render(m.keys[description]),
		m.inputs[description].View(),
		inputStyle.Width(len(m.keys[values])).Render(m.keys[values]),
		m.inputs[values].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
	wrapped := wordwrap.String(content, 15)
	return wrapped
}

// nextInput focuses the next input field
func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *model) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
