package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/reflow/wordwrap"
	"github.com/opentdf/tructl/tui/constants"
)

// You generally won't need this unless you're processing stuff with
// complicated ANSI escape sequences. Turn it on if you notice flickering.
//
// Also keep in mind that high performance rendering only works for programs
// that use the full size of the terminal. We're enabling that below with
// tea.EnterAltScreen().
// Setting this to true is causing issues and preventing items from rendering
const useHighPerformanceRenderer = false

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
	width, height int
	content       string
	title         string
	ready         bool
	viewport      viewport.Model
}

func SetupViewport(m AttributeView, msg tea.WindowSizeMsg) (AttributeView, []tea.Cmd) {
	var (
		cmds []tea.Cmd
	)
	headerHeight := lipgloss.Height(m.headerView())
	footerHeight := lipgloss.Height(m.footerView())
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
		// wrapped := wordwrap.String(m.content, msg.Width)
		// m.viewport.SetContent(wrapped)
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

	if useHighPerformanceRenderer {
		// Render (or re-render) the whole viewport. Necessary both to
		// initialize the viewport and when the window is resized.
		//
		// This is needed for high-performance rendering only.
		cmds = append(cmds, viewport.Sync(m.viewport))
	}
	return m, cmds
}

func (m AttributeView) Init() tea.Cmd {
	return nil
}

func (m AttributeView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "backspace":
			// m.viewport.HighPerformanceRendering = false
			// _, cmd1 := m.viewport.Update(msg)
			// cmd2 := viewport.Sync(m.viewport)
			attributeList := InitAttributeList()
			am, cmd := attributeList.Update(tea.WindowSizeMsg{Width: constants.WindowSize.Width, Height: constants.WindowSize.Height})
			return am, tea.Sequence(
				// cmd1, cmd2,
				tea.ClearScreen, cmd)
		}

	case tea.WindowSizeMsg:
		m, cmds = SetupViewport(m, msg)
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m AttributeView) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	wrapped := wordwrap.String(m.content, m.width)
	m.viewport.SetContent(wrapped)
	return fmt.Sprintf("%s\n%s\n%s", m.headerView(), m.viewport.View(), m.footerView())
}

func (m AttributeView) headerView() string {
	title := titleStyle.
		// Foreground(lipgloss.Color("#00FFFF")).
		Render(m.title)
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(title)))
	// line = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FFFF")).Render(line)
	return lipgloss.JoinHorizontal(lipgloss.Center, title, line)
	// return title
}

func (m AttributeView) footerView() string {
	info := infoStyle.Render(fmt.Sprintf("%3.f%%", m.viewport.ScrollPercent()*100))
	line := strings.Repeat("─", max(0, m.viewport.Width-lipgloss.Width(info)))
	return lipgloss.JoinHorizontal(lipgloss.Center, line, info)
	// return ""
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
