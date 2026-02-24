# otdfctl Architecture

This document describes the architecture of `otdfctl`, the CLI for managing the OpenTDF Platform. It is split into two sections: **As-Is** (current state) and **To-Be** (intended direction). The goal is to enable alignment assessment between what exists in the code and what the architecture intends.

---

## As-Is: Current Architecture

### Overview

otdfctl is a Go CLI application built on [Cobra](https://cobra.dev/) that provides CRUD operations for the OpenTDF Platform's policy objects, TDF encrypt/decrypt, and profile-based configuration. It communicates with the platform via gRPC through the `opentdf/platform/sdk`.

**Module**: `github.com/opentdf/otdfctl`
**Go version**: 1.25.0
**Current version**: 0.29.0 (managed by release-please)

### Layered Architecture

The codebase follows a three-layer architecture:

```
┌─────────────────────────────────────────────────────┐
│  cmd/                Commands (Cobra)               │
│  ├── policy/         Policy CRUD subcommands        │
│  ├── tdf/            Encrypt/decrypt/inspect        │
│  ├── auth/           Login/logout/credentials       │
│  ├── dev/            Dev/playground commands         │
│  ├── common/         Shared command helpers          │
│  ├── root.go         Root command, init, wiring     │
│  ├── profile.go      Profile management commands    │
│  └── interactive.go  TUI launcher                   │
├─────────────────────────────────────────────────────┤
│  pkg/                Core Library                    │
│  ├── cli/            CLI framework (Cli, Printer,   │
│  │                   flags, tables, styles, confirm) │
│  ├── handlers/       SDK wrapper (Handler struct,    │
│  │                   CRUD methods per resource)      │
│  ├── man/            Doc-driven command system       │
│  ├── profiles/       Profile management (stores)     │
│  ├── auth/           OIDC auth flows                 │
│  ├── config/         Build-time constants            │
│  ├── tdf/            TDF type constants              │
│  └── utils/          URL, PEM, file, HTTP helpers    │
├─────────────────────────────────────────────────────┤
│  External                                            │
│  ├── opentdf/platform/sdk        gRPC SDK            │
│  ├── opentdf/platform/protocol   Protobuf types      │
│  └── opentdf/platform/lib        Crypto, flattening  │
└─────────────────────────────────────────────────────┘
```

**Data flow for a typical command**:
```
User → Cobra command → cli.New() → common.NewHandler() → handlers.Method() → SDK → Platform gRPC
                                         ↓
                                   Profile store (auth, endpoint)
                                         ↓
                                   common.HandleSuccess() → Printer (styled | JSON) → stdout
```

### Documentation-Driven Command System

The most distinctive architectural pattern. Commands are defined in Markdown files with YAML frontmatter, not in Go code.

**Location**: `docs/man/` (100+ markdown files)
**Mechanism**: `pkg/man/` parses embedded markdown at `init()` time via `go:embed`

Each markdown file defines:
- `command.name`: The cobra `Use` field
- `command.arguments`: Positional args
- `command.aliases`: Command aliases
- `command.flags`: Flag definitions (name, description, shorthand, default, enum)
- `command.hidden`: Whether the command is hidden
- Markdown body: Becomes the `Long` description (rendered via glamour)

**Usage pattern in Go code**:
```go
// Retrieve a doc-defined command and attach a run function
doc := man.Docs.GetCommand("policy/attributes/create", man.WithRun(createAttribute))
// Retrieve flag definitions from the doc
doc.Flags().StringP(doc.GetDocFlag("name").Name, ...)
```

**Supports i18n**: Files can be suffixed with language codes (e.g., `create.fr.md`), with English as the default fallback.

### Command Initialization Pattern

All command groups follow an `InitCommands()` pattern:

1. `cmd/root.go:init()` adds top-level commands to `RootCmd`
2. `cmd/root.go:init()` calls `InitCommands()` on each command group
3. Each `InitCommands()` calls per-resource `init*Commands()` functions
4. `init*Commands()` wires doc-based commands, flags, and run functions

This explicit initialization (rather than Go `init()`) gives control over ordering and avoids circular dependencies.

### Handler Pattern (Service Layer)

`pkg/handlers/Handler` wraps the platform SDK:

```go
type Handler struct {
    sdk              *sdk.SDK
    platformEndpoint string
}
```

Every command that talks to the platform follows this flow:
1. `c := cli.New(cmd, args)` — create CLI context with flag helpers and printer
2. `h := common.NewHandler(c)` — load profile, authenticate, create SDK connection
3. `defer h.Close()` — clean up the SDK connection
4. Call a handler method (e.g., `h.CreateAttribute(...)`)
5. Format output and exit via `common.HandleSuccess()` or `cli.ExitWith*()`

**`common.NewHandler()`** orchestrates auth:
- If `--host` + auth flags are set → creates an in-memory profile (no filesystem)
- Otherwise → loads from the profile store (filesystem)
- Validates credentials before returning

**Handler files** (18 files in `pkg/handlers/`):
Each file corresponds to a platform resource type (attributes, namespaces, kas-registry, etc.) and wraps the SDK's gRPC client methods.

### CLI Framework (`pkg/cli`)

The `Cli` struct is created at the start of every command handler:

```go
type Cli struct {
    cmd     *cobra.Command
    args    []string
    Flags   *flagHelper   // typed flag accessors
    printer *Printer      // dual-mode output
}
```

Key components:
- **`flagHelper`** (`flagValues.go`): Typed flag accessors (`GetRequiredString`, `GetRequiredID`, `GetOptionalBool`, `GetStringSlice`, `GetState`, etc.)
- **`Printer`** (`printer.go`): Dual-mode output controlled by `--json` flag or profile `outputFormat` setting
- **`errors.go`**: Exit functions (`ExitWithError`, `ExitWithNotFoundError`, `ExitWithWarning`, `ExitWithSuccess`) — all call `os.Exit()`
- **`table.go`**: Table rendering via `evertras/bubble-table`
- **`tabular.go`**: Key-value display for single-resource views
- **`style.go`**: Adaptive color palette and lipgloss styles
- **`confirm.go`**: Interactive confirmation prompts via `charmbracelet/huh`
- **`sdkHelpers.go`**: SDK data transformers (`GetSimpleAttribute`, `ConstructMetadata`, key algorithm converters)
- **`pipe.go`**: Stdin/file reading utilities

All command code accesses flags through `c.Flags` (e.g., `c.Flags.GetRequiredString("name")`).

### Dual Output Mode

Every command supports two output formats (per ADR-0001):

| Mode | Trigger | Behavior |
|------|---------|----------|
| **Styled** | Default, or profile `outputFormat: styled` | lipgloss-colored terminal output with bubble-table rendering |
| **JSON** | `--json` flag, or profile `outputFormat: json` | Structured JSON to stdout/stderr |

The `common.HandleSuccess()` function checks both the `--json` flag and the profile output format to decide which mode to use.

### Profile System

Replaces traditional config files with a profile-based approach:

**Package**: `pkg/profiles/` (wraps `jrschumacher/go-osprofiles`)

**Profile drivers**:
| Driver | Usage |
|--------|-------|
| Filesystem | Default. Stores in OS app config directories |
| Keyring | OS keyring (deprecated, auto-migrates to filesystem) |
| In-memory | For flag-based ephemeral usage (`--host` + auth flags) |

**Profile contents** (`ProfileConfig`):
- `Name`, `Endpoint`, `TLSNoVerify`, `OutputFormat`
- `AuthCredentials`: type (client-credentials or access-token), client ID/secret, scopes, access token

**Migration**: `common.InitProfile()` auto-migrates from keyring to filesystem on every invocation.

### Authentication

**Package**: `pkg/auth/`

Three auth mechanisms:
1. **Client credentials** — stored in profile or via `--with-client-creds` / `--with-client-creds-file` flags
2. **Access token** — stored in profile or via `--with-access-token` flag
3. **PKCE login** — interactive browser-based flow via `auth login`

Auth is resolved in `common.NewHandler()` and validated via `auth.ValidateProfileAuthCredentials()` before SDK initialization.

### Error Handling

Commands use `cli.ExitWith*()` functions which:
1. Print a styled or JSON error message (respecting output format)
2. Call `os.Exit()` with an appropriate exit code

gRPC `NotFound` status codes get special treatment via `ExitWithNotFoundError`. Sentinel errors are defined per package in `errors.go` files.

Destructive operations use `cli.ConfirmAction()` / `cli.ConfirmTextInput()` unless `--force` is passed.

### Command Structure

```
otdfctl
├── auth
│   ├── login
│   ├── logout
│   ├── client-credentials
│   ├── clear-client-credentials
│   └── print-access-token
├── profile (aliases: profiles, prof)
│   ├── create, list, get, delete, delete-all
│   ├── set-default, set-endpoint, set-output-format
│   ├── migrate, cleanup
├── encrypt, decrypt, inspect  (TDF group)
├── policy
│   ├── actions [get, list, create, update, delete]
│   ├── attributes [get, list, create, update, deactivate]
│   │   ├── values [get, list, create, update, deactivate]
│   │   │   ├── key [assign, remove]
│   │   │   └── unsafe [delete, reactivate, update]
│   │   ├── key [assign, remove]
│   │   ├── unsafe [delete, reactivate, update]
│   │   └── namespaces [get, list, create, update, deactivate]
│   │       ├── key [assign, remove]
│   │       └── unsafe [delete, reactivate, update]
│   ├── kas-registry [get, list, create, update, delete]
│   │   └── key [get, list, create, update, rotate, import, list-mappings]
│   │       ├── base [get, set]
│   │       └── unsafe [delete]
│   ├── kas-grants [assign, unassign, list]
│   ├── key-management
│   │   └── provider [get, list, create, update, delete]
│   ├── subject-condition-sets [get, list, create, update, delete]
│   ├── subject-mappings [get, list, create, update, delete]
│   ├── obligations [get, list, create, update, delete]
│   ├── resource-mappings [get, list, create, update, delete]
│   ├── resource-mapping-groups [get, list, create, update, delete]
│   └── registered-resources [get, list, create, update, delete]
├── dev
│   ├── design-system
│   └── selectors [generate, test]
└── interactive (launches Bubble Tea TUI)
```

### TUI (`tui/`)

A Bubble Tea-based interactive mode. Launched via `otdfctl interactive`.

**Status**: Work in progress (README warns against modifying)

Contains views for: app menu, attribute list/detail/create, label list/update, read, update. Uses `tui/constants/` for shared state and `tui/form/` for form components.

### Build System

**Makefile targets**:
| Target | Description |
|--------|-------------|
| `run` | `go run .` |
| `test` | `go test -v ./...` |
| `build` | Cross-compile for 8 targets (darwin/linux/windows x amd64/arm/arm64), zip, checksum |
| `build-test` | Build with `TestMode=true` for BATS e2e testing |
| `test-bats` | Build test binary, resize terminal, run BATS |
| `clean` | Remove `target/` |

Build-time injection via `-ldflags`: `Version`, `CommitSha`, `BuildTime`, and optionally `TestMode`.

**CI/CD** (GitHub Actions):
- `ci.yaml`: govulncheck, golangci-lint, unit tests (short, race, cover), e2e (BATS against containerized platform), platform cross-tests
- `release.yaml`: release-please versioning, cross-compile, artifact upload
- Supporting: codeql, pr-lint, backport, dependabot, dependency-review

### Testing

**Unit tests**: Minimal — 3 test files:
- `cmd/execute_test.go` — Execute/MountRoot functions
- `pkg/utils/identifier_test.go` — `NormalizeEndpoint()`
- `pkg/utils/pemvalidate_test.go` — PEM validation

**End-to-end tests (BATS)**: Primary testing strategy — 19 test suites in `e2e/`:
- Cover all major command groups (attributes, namespaces, auth, encrypt-decrypt, kas, etc.)
- Run against a real platform instance (containerized in CI)
- Test both styled and `--json` output modes
- Test binary built with `TestMode=true` (in-memory profile store)
- Terminal size controlled for consistent output testing
- Optional TestRail integration for result uploading

### Key Dependencies

| Category | Library | Purpose |
|----------|---------|---------|
| CLI | `spf13/cobra` | Command framework |
| SDK | `opentdf/platform/sdk` | Platform gRPC communication |
| Types | `opentdf/platform/protocol/go` | Protobuf-generated types |
| TUI | `charmbracelet/bubbletea` | Terminal UI framework |
| Style | `charmbracelet/lipgloss` | Terminal styling |
| Tables | `evertras/bubble-table` | Table rendering |
| Docs | `charmbracelet/glamour` | Markdown rendering |
| Prompts | `charmbracelet/huh` | Interactive forms |
| Auth | `zitadel/oidc/v3` | OIDC client |
| JWT | `go-jose/go-jose/v3`, `golang-jwt/jwt/v5` | JWT parsing |
| Profiles | `jrschumacher/go-osprofiles` | OS profile storage |
| MIME | `gabriel-vasile/mimetype` | File type detection |
| Frontmatter | `adrg/frontmatter` | YAML frontmatter parsing |

### Conventions and Patterns Summary

1. **Doc-driven commands**: Command metadata lives in `docs/man/*.md`, not Go code
2. **Three-layer separation**: `cmd/` → `pkg/handlers/` → `platform/sdk`
3. **Functional options**: Used for `Handler` and `Cli` construction
4. **Dual output**: Every command supports styled and JSON output
5. **Profile-based config**: No config files; profiles store everything
6. **Explicit initialization**: `InitCommands()` pattern over Go `init()`
7. **Confirmation for destructive ops**: `ConfirmAction`/`ConfirmTextInput` with `--force` bypass
8. **Mountable design**: `MountRoot()` allows embedding as a subcommand of another CLI
9. **Package-level variable state**: Some flag values stored as package-level vars (e.g., `attributeValues`, `metadataLabels`)
10. **Exit-on-error**: Commands call `os.Exit()` on error rather than returning errors up the stack

---

## To-Be: Intended Architecture

### Decided Direction

**`FlagHelper` → `Flags` migration**: Complete. The deprecated `FlagHelper` alias was removed; all code now uses `c.Flags`.

**Deprecated `config` command**: Removed. The `config` command was deprecated in PR #719 when output format storage moved to the profile system, and has since been deleted along with its docs.

**Testing strategy**: E2E (BATS) is the primary testing approach and provides good coverage for the CRUD-heavy command surface. Unit tests are appropriate for non-CRUD helpers and utility functions (e.g., `pkg/utils/`, `pkg/cli/` helpers, `pkg/man/` parsing). Direction: add unit tests for non-trivial helper logic, not for CRUD command wiring.

**Package-level variable state**: The `StringSliceVarP` pattern (package-level vars bound to cobra flags) is a cobra limitation for slice flags, not an architectural concern. The pattern of declaring the var, binding via `VarP`, then reading through `c.Flags.GetStringSlice()` with validation is intentional and should be maintained where cobra requires it.

**`NewHandler` as preRun hook**: There is a TODO and [issue #383](https://github.com/opentdf/otdfctl/issues/383) to move handler initialization to a Cobra `PersistentPreRunE` hook. This would reduce boilerplate in every command handler. Direction: pursue when ready, but not a blocker.

### Open Questions

- What is the intended future of the TUI (`tui/` directory)? (Undecided)
- How should the profile system evolve? (Current state is functional)

### Known TODOs in Code

From `README.md`:
- [ ] Add support for JSON input as piped input
- [ ] Add help level handler for each command
- [ ] Add support for `--verbose` persistent flag
- [x] Helper functions to support common tasks (done via `pkg/cli`)

From `cmd/common/common.go`:
- [ ] Make `NewHandler` a preRun hook ([#383](https://github.com/opentdf/otdfctl/issues/383))

### Alignment Concerns

| Area | Status | Notes |
|------|--------|-------|
| `FlagHelper` → `Flags` | Complete | All code uses `c.Flags`; alias removed |
| Deprecated `config` cmd | Removed | Deleted along with docs |
| Unit test coverage | Low | Only 3 test files; non-CRUD helpers lack coverage |
| Pattern consistency | Mostly consistent | Policy commands follow the same pattern well |
| Output formatting | Consistent | Dual-mode (styled/JSON) is well-implemented |
| Flag boilerplate | Verbose but consistent | `GetDocFlag` pattern is repetitive but uniform |
| Error handling | Consistent | Exit-on-error pattern used throughout |
