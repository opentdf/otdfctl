# TUI Development

> **⚠️ CAUTION**
>
> The TUI is a work in progress. **Avoid touching TUI code until the framework is defined.**
>
> Only work on TUI components if explicitly instructed by the user.

## Overview

The Terminal User Interface (TUI) provides an interactive experience for otdfctl users, built with:

- **Bubble Tea** - TUI framework (github.com/charmbracelet/bubbletea)
- **Lipgloss** - Terminal styling (github.com/charmbracelet/lipgloss)
- **Bubbles** - TUI components (github.com/charmbracelet/bubbles)

## Current State

The TUI is partially implemented with components for:

- App menu navigation
- Attribute list viewing
- Attribute creation/editing
- Label management

However, the overall framework and patterns are still being defined.

## Directory Structure

```
tui/
├── common.go                - Entry point and shared utilities
├── appMenu.go               - Main menu
├── attributeList.go         - Attribute list view
├── attributeView.go         - Single attribute view
├── attributeCreateView.go   - Attribute creation form
├── labelList.go             - Label list view
├── labelUpdate.go           - Label update form
├── read.go                  - Read operations
├── update.go                - Update operations
├── shell.go                 - Shell/REPL mode
├── constants/
│   └── consts.go            - Constants (window size, etc.)
└── form/
    └── addAttribute.go      - Attribute form
```

## Architecture Patterns

### Bubble Tea Model

Bubble Tea uses the Elm Architecture:

1. **Model** - Application state
2. **Update** - State updates based on messages
3. **View** - Rendering the UI

Example from `tui/common.go:16-35`:

```go
func StartTea(h handlers.Handler) error {
    // Initialize debugging
    f, _ := tea.LogToFile("debug.log", "help")
    defer f.Close()

    // Initialize model
    m, _ := InitAppMenu(h)

    // Start program
    constants.P = tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
    if _, err := constants.P.Run(); err != nil {
        return err
    }
    return nil
}
```

### Message Passing

Components communicate via messages:

```go
type MyMsg struct {
    Data string
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case MyMsg:
        // Handle message
        m.data = msg.Data
    case tea.KeyMsg:
        // Handle keyboard input
    }
    return m, nil
}
```

### View Rendering

Views return strings with ANSI escape codes for styling:

```go
func (m Model) View() string {
    style := lipgloss.NewStyle().Padding(1, 2)
    return style.Render("Hello, World!")
}
```

## Components

### Bubbles Components

The TUI uses pre-built components from `github.com/charmbracelet/bubbles`:

- **list** - Interactive lists
- **textinput** - Text input fields
- **table** - Tables
- **viewport** - Scrollable content

### Custom Components

Custom components are being developed in `/tui`:

- App menu for navigation
- Attribute management forms
- Label editors

## Integration with CLI

The TUI shares handlers with the CLI via `pkg/handlers/`.

Flow:
1. TUI component needs data
2. Calls handler function
3. Handler uses SDK to fetch data
4. TUI displays result

This ensures business logic is shared between CLI and TUI.

## Debugging

TUI logs to `debug.log`:

```go
tea.LogToFile("debug.log", "help")
```

View logs while developing:
```bash
tail -f debug.log
```

## Interactive Mode

The CLI has an interactive command that launches the TUI:

```bash
otdfctl interactive
```

See `cmd/interactive.go` for the entry point.

## Shell Mode

A REPL-style shell is being developed in `tui/shell.go` and `cmd/shell.go`.

This would allow:
```bash
otdfctl shell
> policy attributes list
> policy namespaces create --name test
```

Current status: Work in progress.

## Current Limitations

- Framework patterns not finalized
- Some components incomplete
- Navigation flow being refined
- Error handling needs improvement
- Integration between components in progress

## Future Plans

The TUI aims to provide:

1. **Interactive exploration** - Browse policy resources
2. **Guided workflows** - Step-by-step configuration
3. **Visual feedback** - Better visibility into operations
4. **Form-based input** - Alternative to flags

## If You Must Work on TUI

If explicitly instructed to work on TUI components:

1. **Read existing code** - Understand current patterns
2. **Follow Bubble Tea docs** - https://github.com/charmbracelet/bubbletea
3. **Test thoroughly** - TUI bugs are hard to debug
4. **Use the logger** - Log liberally to `debug.log`
5. **Keep it simple** - Don't over-engineer

### Example: Adding a Simple View

```go
package tui

import (
    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)

type MyViewModel struct {
    content string
}

func (m MyViewModel) Init() tea.Cmd {
    return nil
}

func (m MyViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.KeyMsg:
        if msg.String() == "q" {
            return m, tea.Quit
        }
    }
    return m, nil
}

func (m MyViewModel) View() string {
    style := lipgloss.NewStyle().
        Padding(1, 2).
        Border(lipgloss.RoundedBorder())
    return style.Render(m.content)
}
```

## References

- **Bubble Tea Tutorial** - https://github.com/charmbracelet/bubbletea/tree/master/tutorials
- **Lipgloss Examples** - https://github.com/charmbracelet/lipgloss/tree/master/examples
- **Bubbles Components** - https://github.com/charmbracelet/bubbles

## Remember

**When in doubt, don't touch the TUI.** Work on CLI commands instead.
