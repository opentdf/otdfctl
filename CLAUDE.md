# otdfctl - OpenTDF Platform CLI

## What is this project?

**otdfctl** is a command-line interface for managing the OpenTDF Platform. It's built with Go using the Cobra CLI framework and provides CRUD operations for platform resources.

### Tech Stack
- **Language**: Go 1.24
- **CLI Framework**: Cobra
- **TUI Framework**: Bubble Tea, Lipgloss (work in progress)
- **Testing**: Go testing + BATS (Bash Automated Testing System)
- **Platform SDK**: OpenTDF Platform SDK (github.com/opentdf/platform/sdk)

### Project Structure
```
cmd/          - Cobra command definitions (auth, policy, tdf, config, dev)
pkg/          - Shared packages (auth, cli, config, handlers, profiles, man, tdf, utils)
tui/          - Terminal UI components (work in progress - avoid touching)
docs/man/     - Documentation that drives CLI help text
e2e/          - End-to-end BATS tests
agents/       - Agent documentation (progressive disclosure)
```

## Why does this exist?

This CLI simplifies setup, facilitates migration, and aids in configuration management for OpenTDF Platform instances. It uses profile-based configuration to manage connections to different platform instances.

## How to work on this project

### Essential Commands
- **Run**: `go run .` or `make run`
- **Test**: `go test -v ./...` or `make test`
- **Build**: `make build` (cross-platform) or `go build`
- **BATS Tests**: `make test-bats` (requires platform running)

### Before You Start

**Read the relevant agent docs** (in `agents/` directory) based on your task:
- `building_and_testing.md` - Build process, testing, BATS setup
- `project_structure.md` - Detailed codebase layout and component descriptions
- `development_workflow.md` - Adding commands, handling errors, documentation
- `tui_development.md` - TUI framework (⚠️ work in progress)

Only read the docs that are relevant to your current task.

### Key Patterns
- Output format (styled/json) is stored per-profile and can be overridden with `--json` flag
- All commands use the `cli.New()` helper and handlers in `pkg/handlers/`
- Documentation in `/docs/man` drives CLI help text via `man.Docs.GetDoc()`
- Test mode can be enabled by building with `make build-test`

### Critical Notes
- The TUI is work in progress - avoid touching it unless specifically instructed
- Never use interactive git commands (`-i` flag) as they're not supported
- Use linters/formatters instead of manual style corrections
