# Contributor Guide: otdfctl

## Project Structure
- **Commands:** Add new commands in `cmd/` using Cobra. Each command should have a corresponding handler in `pkg/handlers/`.
- **Handlers:** Place business logic in `pkg/handlers/`.
- **Configuration:** Use `pkg/config/` and `otdfctl.yaml` for config management.
- **Authentication/Profiles:** Use `pkg/auth/` and `pkg/profiles/` for auth flows and user profiles.
- **Documentation:** Update or add Markdown docs in `docs/man/` for each command. These docs are parsed at runtime for CLI help.
- **TUI:** Experimental; avoid major changes unless contributing to TUI development.
- **Testing:** Add/extend BATS tests in `e2e/` for new features. Use test mode for development.

## How to Add a Command

1. **Create or Update the Command File:**
   - For a new top-level command, create a new file in `cmd/` (e.g., `cmd/foo.go`).
   - For a subcommand, add it to the appropriate parent command file (e.g., add a `hello` subcommand to `cmd/dev.go`).
   - Use the CLI helpers (`cli.New`, `c.Args.GetOptionalString`, `c.Flags.GetOptionalBool`, etc.) for argument and flag parsing, following the style in [`docs/example-add-subcommand.md`](docs/example-add-subcommand.md).

2. **Implement Business Logic:**
   - Place complex logic in `pkg/handlers/` and call it from your command handler.

3. **Add Documentation:**
   - Create or update the Markdown documentation for your command or subcommand in `docs/man/<command>/`.
   - Follow the format shown in the example, including argument and flag descriptions, and usage examples.

4. **Add or Update Tests:**
   - Add or extend BATS tests in `e2e/` to cover your new command and its options.

5. **Register the Command:**
   - Ensure your command or subcommand is registered with its parent in the `init()` function.

6. **Verify and Sync:**
   - Run your command to verify it works as expected.
   - Check that CLI help output matches your documentation.

See [`docs/example-add-subcommand.md`](docs/example-add-subcommand.md) for a step-by-step example.

## Documentation-Driven CLI
- The CLI help system is powered by Markdown docs in `docs/man/`.
- Always update docs when adding or changing commands/flags.

## Testing
- Run all BATS tests: `make test-bats`
- Run a specific test: `bats e2e/<test>.bats`
- Use test mode for development: `make build-test`

## TUI Guidelines
- The TUI in `tui/` is experimental. Avoid changes unless coordinated with maintainers.

## Syncing Docs and Code
- When adding or changing commands, always update the corresponding Markdown doc in `docs/man/`.
- Review CLI help output to ensure it matches the documentation.

## Known Limitations
- Some authentication/profile features may not work on all platforms.
- TUI is not production-ready.

---
Update this guide as the project evolves.
