package tui

import (
	"fmt"
	"strings"

	// "log"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
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

var (
	inputStyle    = lipgloss.NewStyle().Foreground(constants.Magenta)
	continueStyle = lipgloss.NewStyle().Foreground(cyan)
)

type model struct {
	inputs  []interface{}
	focused int
	err     error
	keys    []string
}

func InitAttributeEdit(names []string, item AttributeItem) (tea.Model, tea.Cmd) {
	// inputs := make([]interface{}, len(names)) //= [make([]tea.Model, len(names))
	var inputs []interface{}
	// var inputs []textinput.Model = make([]textinput.Model, len(names))
	// inputs = append(inputs, textinput.New())
	ti0 := textinput.New()
	ti0.Focus()
	ti0.SetValue(item.id)
	inputs = append(inputs, ti0)

	// inputs[id] = textinput.New()
	// inputs[id].Placeholder = "4505 **** **** 1234"
	// inputs[id].Focus()
	// inputs[id].SetValue(item.id)
	// inputs[ccn].CharLimit = 20
	// inputs[id].Width = 30
	// inputs[id].Prompt = ""
	ti1 := textinput.New()
	ti1.SetValue(item.name)
	inputs = append(inputs, ti1)
	// inputs[name] = textinput.New()
	// inputs[name].Placeholder = "MM/YY "
	// inputs[exp].CharLimit = 5
	// inputs[name].Width = 5
	// inputs[name].Prompt = ""
	// inputs[name].SetValue(item.name)

	ti2 := textinput.New()
	ti2.SetValue(item.namespace)
	inputs = append(inputs, ti2)
	// inputs[namespace] = textinput.New()
	// inputs[namespace].Placeholder = "XXX"
	// inputs[cvv].CharLimit = 3
	// inputs[namespace].Width = 5
	// inputs[namespace].Prompt = ""

	// inputs[namespace].SetValue(item.namespace)
	ti3 := textinput.New()
	ti3.SetValue(item.rule)
	inputs = append(inputs, ti3)
	// inputs[rule] = textinput.New()
	// inputs[rule].SetValue(item.rule)
	ti4 := textarea.New()
	ti4.SetValue(item.description)
	inputs = append(inputs, ti4)
	// inputs[description] = textarea.New()
	// inputs[description].SetValue(item.description)
	ti5 := textinput.New()
	ti5.SetValue(strings.Join(item.values, ","))
	inputs = append(inputs, ti5)
	// inputs[values] = textinput.New()
	// inputs[values].SetValue(strings.Join(item.values, ","))
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
			if m.focused == len(m.inputs)-1 {
				return saveModel, saveCmd
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP, tea.KeyUp:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN, tea.KeyDown:
			m.nextInput()
		}
		for i := range m.inputs {
			// var (
			// 	tempInput textinput.Model
			// 	tempArea  textarea.Model
			// )
			if i == description {
				tempInput := m.inputs[i].(textarea.Model)
				tempInput.Blur()
				m.inputs[i] = tempInput
			} else {
				tempArea := m.inputs[i].(textinput.Model)
				tempArea.Blur()
				m.inputs[i] = tempArea
			}
			// tempInput := m.inputs[i].(textinput.Model)
			// tempInput.Blur()
			// m.inputs[i] = tempInput
			// .Blur()
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
		// m.inputs[m.focused].Focus()
	// case tea.WindowSizeMsg:
	// 	for i := range m.inputs {
	// 		m.inputs[i].Width = msg.Width
	// 	}
	// 	return m, nil
	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		if i == description {
			m.inputs[i], cmds[i] = m.inputs[i].(textarea.Model).Update(msg)
		} else {
			m.inputs[i], cmds[i] = m.inputs[i].(textinput.Model).Update(msg)
		}
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
	return fmt.Sprintf("\n\n%s\n"+CreateEditFormat(len(m.inputs))+"\n%s",
		"[Edit Attribute]",
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
		// fmt.Sprintf("[%s]", m.inputs[values].View()),
		m.inputs[values].(textinput.Model).View(),
		continueStyle.Render("<<Save>>"),
	) + "\n"
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
