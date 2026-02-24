# otdfctl Architecture

## Overview

`otdfctl` is a modular, documentation-driven CLI tool. Its architecture ensures that command definitions, arguments, and flags are defined in Markdown documentation files (`docs/man/`). These docs are parsed at runtime to generate CLI help and to drive the registration of command-line arguments and flags, ensuring that code and documentation remain in sync.

## Project Structure

- **Commands:**
  - Source files for commands are in `cmd/`.
  - Each command is registered using `man.Docs.GetCommand`, which loads its definition from the corresponding Markdown doc in `docs/man/`.
  - The command string passed to `GetCommand` uses forward slashes (e.g., `dev/hello`, `policy/subject-mappings/update`) to represent the command path.
  - Arguments and flags are registered in code using the doc-driven helpers, e.g.:
    ```go
    cmd.Flags().StringP(
        cmd.GetDocFlag("flagname").Name,
        cmd.GetDocFlag("flagname").Shorthand,
        cmd.GetDocFlag("flagname").Default,
        cmd.GetDocFlag("flagname").Description,
    )
    ```
    This ensures the flag's name, shorthand, default, and description are always sourced from the documentation.
- **Handlers:**
  - Business logic is implemented in `pkg/handlers/`.
  - Command handlers in `cmd/` delegate to these packages for core operations.
- **Configuration:**
  - Configuration management is handled in `pkg/config/` and `otdfctl.yaml`.
- **Authentication/Profiles:**
  - Authentication flows and user profiles are managed in `pkg/auth/` and `pkg/profiles/`.
- **Documentation:**
  - All command and subcommand documentation lives in `docs/man/`.
  - These Markdown files define command names (with forward slashes), arguments, flags, descriptions, and usage examples.
  - The CLI help system is generated from these docs at runtime.
- **TUI:**
  - Experimental text-based UI in `tui/`.
- **Testing:**
  - End-to-end BATS tests are in `e2e/`.

## Command and Flag Registration Pattern

- The canonical source for command structure, arguments, and flags is the Markdown documentation in `docs/man/`.
- In Go code, commands are registered using `man.Docs.GetCommand("path/to/command", ...)`.
- Flags and arguments are registered using the `GetDocFlag` helper, which pulls all metadata from the doc:
  ```go
  cmd.Flags().StringP(
      cmd.GetDocFlag("flagname").Name,
      cmd.GetDocFlag("flagname").Shorthand,
      cmd.GetDocFlag("flagname").Default,
      cmd.GetDocFlag("flagname").Description,
  )
  ```
- This pattern is used for all commands and subcommands, ensuring that the CLI and its documentation are always in sync.

## Adding or Modifying Commands

1. **Write or update the Markdown documentation** for your command in `docs/man/<command>/`. This defines the command's arguments, flags, and help text.
2. **Implement the command handler** in `cmd/`, using `man.Docs.GetCommand` to load the command and inject flags/args from the doc using `GetDocFlag`.
3. **Add or update business logic** in `pkg/handlers/` as needed.
4. **Add or update tests** in `e2e/`.
5. **Run and verify** the command, ensuring CLI help matches the documentation.

See `docs/example-add-subcommand.md` for a step-by-step example.

---

Update this document as the architecture evolves.
