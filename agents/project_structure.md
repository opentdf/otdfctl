# Project Structure

This document provides a detailed breakdown of the otdfctl codebase organization.

## Root Directory

```
otdfctl/
├── cmd/              - Cobra command definitions
├── pkg/              - Shared packages and business logic
├── tui/              - Terminal UI components (WIP - avoid)
├── docs/             - Documentation
├── e2e/              - BATS end-to-end tests
├── agents/           - Agent documentation (this directory)
├── adr/              - Architecture Decision Records
├── .github/          - CI/CD workflows and scripts
├── main.go           - Application entry point
├── Makefile          - Build and test automation
└── go.mod            - Go module definition
```

## `/cmd` - Command Definitions

Command-line interface built with Cobra. Each subdirectory represents a command group.

### Structure
```
cmd/
├── root.go                 - Root command setup, version handling, logging
├── execute.go              - Command execution logic
├── interactive.go          - Interactive mode (TUI)
├── profile.go              - Profile management commands
├── auth/                   - Authentication commands
│   ├── auth.go
│   ├── login.go
│   ├── logout.go
│   ├── clientCredentials.go
│   ├── printAccessToken.go
│   └── clearCachedCredentials.go
├── config/                 - Configuration commands
├── policy/                 - Policy management (attributes, namespaces, etc.)
│   ├── policy.go
│   ├── attributes.go
│   ├── attributeValues.go
│   ├── namespaces.go
│   ├── kasRegistry.go
│   ├── kasKeys.go
│   ├── kasGrants.go
│   ├── subjectMappings.go
│   ├── resourceMappings.go
│   ├── resourceMappingGroups.go
│   ├── registeredResources.go
│   ├── subjectConditionSets.go
│   ├── obligations.go
│   ├── actions.go
│   ├── baseKeys.go
│   ├── keyManagement.go
│   └── keyManagementProvider.go
├── tdf/                    - TDF operations
│   ├── tdf.go
│   ├── encrypt.go
│   ├── decrypt.go
│   └── inspect.go
├── dev/                    - Development/debugging commands
└── common/                 - Common command utilities
```

### Key Files
- **cmd/root.go:24** - `RootCmd` definition
- **cmd/root.go:86** - Command tree setup
- **cmd/root.go:63** - Persistent pre-run for logging setup

## `/pkg` - Shared Packages

Core business logic and utilities.

### `/pkg/auth` - Authentication
Handles OAuth2 flows, token management, and authentication with the platform.

- `auth.go` - Core authentication logic
- `errors.go` - Authentication error types

### `/pkg/cli` - CLI Utilities
Common CLI operations and helpers.

- `cli.go` - Main CLI context and initialization
- `printer.go` - Output formatting (styled vs JSON)
- `table.go` - Table rendering for terminal
- `tabular.go` - Tabular data helpers
- `style.go` - Terminal styling
- `confirm.go` - User confirmation prompts
- `pipe.go` - Piped input handling
- `flagValues.go` - Flag value extraction
- `errors.go` - CLI error handling
- `messages.go` - User-facing messages
- `sdkHelpers.go` - SDK integration helpers
- `utils.go` - General utilities
- `clioptions.go` - CLI option handling

Key pattern: Commands use `cli.New(cmd, args)` to get a CLI context with access to flags and utilities.

### `/pkg/config` - Configuration
Application configuration and build-time constants.

- `config.go` - Config structures, version info, test mode flag

Build-time variables set via `-ldflags`:
- `Version`
- `CommitSha`
- `BuildTime`
- `TestMode`

### `/pkg/handlers` - Business Logic
Core handler functions for CRUD operations. These are called by commands in `/cmd`.

**Policy handlers:**
- `attribute.go` - Attribute CRUD
- `attributeValues.go` - Attribute value CRUD
- `namespaces.go` - Namespace CRUD
- `kas-registry.go` - KAS registry operations
- `kas-keys.go` - KAS key management
- `kas-grants.go` - KAS grant management
- `subjectmappings.go` - Subject mapping CRUD
- `resourceMappings.go` - Resource mapping CRUD
- `resourceMappingGroups.go` - Resource mapping group CRUD
- `registeredResources.go` - Registered resource CRUD
- `subjectConditionSets.go` - Subject condition set CRUD
- `obligations.go` - Obligation management
- `actions.go` - Action management
- `base-keys.go` - Base key operations
- `provider-config.go` - Provider configuration

**Other handlers:**
- `sdk.go` - SDK initialization and helpers
- `tdf.go` - TDF encrypt/decrypt/inspect operations
- `selectors.go` - Resource selection utilities

Pattern: Handlers receive parameters, interact with the OpenTDF Platform SDK, and return results.

### `/pkg/profiles` - Profile Management
User profile storage and retrieval (connection endpoints, credentials, output format).

- `profile.go` - Profile CRUD operations
- `profileConfig.go` - Profile configuration structures
- `profileAuthCreds.go` - Credential storage (keyring integration)
- `errors.go` - Profile error types

Uses OS keyring (via `github.com/jrschumacher/go-osprofiles`) to store credentials securely.

### `/pkg/man` - Documentation
CLI documentation system that drives help text.

- `man.go` - Documentation loading and retrieval
- `docflags.go` - Flag documentation
- `style.go` - Documentation styling

Pattern: Commands use `man.Docs.GetDoc("<command>")` to load documentation from `/docs/man/*.md`.

### `/pkg/tdf` - TDF Operations
TDF (Trusted Data Format) encrypt/decrypt/inspect logic.

- `tdf.go` - TDF operations implementation

### `/pkg/utils` - Utilities
General-purpose utility functions.

- `identifier.go` - ID parsing and validation
- `validators.go` - Input validation
- `read.go` - File reading utilities
- `http.go` - HTTP helpers
- `pemvalidate.go` - PEM validation

## `/tui` - Terminal UI (Work in Progress)

**⚠️ CAUTION**: This is work in progress. Avoid touching until framework is defined.

Built with Bubble Tea and Lipgloss for interactive terminal experiences.

```
tui/
├── common.go
├── appMenu.go
├── attributeList.go
├── attributeView.go
├── attributeCreateView.go
├── labelList.go
├── labelUpdate.go
├── read.go
├── update.go
├── shell.go
├── constants/
│   └── consts.go
└── form/
    └── addAttribute.go
```

## `/docs` - Documentation

### `/docs/man` - Man Page Documentation
Markdown files that drive CLI help text via `man.Docs.GetDoc()`.

Pattern: Each command has a corresponding markdown file with frontmatter defining command metadata.

Example structure:
```markdown
---
command: command-name
parent: parent-command
description: Short description
---

# Longer description and examples
```

## `/e2e` - End-to-End Tests

BATS (Bash Automated Testing System) tests for full CLI workflow testing.

```
e2e/
├── *.bats                     - BATS test files
├── resize_terminal.sh         - Terminal size control
└── testrail-integration/      - TestRail upload scripts
```

## `/adr` - Architecture Decision Records

Documents architectural decisions made in the project.

## Entry Point

**main.go** - Application entry point. Calls `cmd.Execute()` to run the CLI.

## Key Integration Points

### Command → Handler → SDK Flow

1. User runs command (e.g., `otdfctl policy attributes create`)
2. Command in `/cmd/policy/attributes.go` parses flags
3. Command calls handler in `/pkg/handlers/attribute.go`
4. Handler uses OpenTDF Platform SDK to make API calls
5. Handler returns result
6. Command uses `/pkg/cli/printer.go` to format and display output

### Profile-Based Configuration

1. User creates profile: `otdfctl profile create <name> <endpoint>`
2. Profile stored via `/pkg/profiles/` (config + keyring)
3. Commands load profile to get endpoint and credentials
4. SDK initialized with profile credentials
5. Output format (styled/json) stored per-profile, overridable with `--json` flag

### Documentation System

1. Documentation written in `/docs/man/*.md`
2. Commands use `man.Docs.GetDoc("<command>")` to load docs
3. Frontmatter defines command structure
4. Content used for help text

## Dependencies

Key external dependencies (see go.mod:1-114):

- **github.com/spf13/cobra** - CLI framework
- **github.com/opentdf/platform/sdk** - OpenTDF Platform SDK
- **github.com/charmbracelet/bubbletea** - TUI framework
- **github.com/charmbracelet/lipgloss** - Terminal styling
- **golang.org/x/oauth2** - OAuth2 authentication
- **google.golang.org/grpc** - gRPC communication
