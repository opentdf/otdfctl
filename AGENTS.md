# Repository Guidelines

## Project Structure & Module Organization

- `main.go`: CLI entrypoint.
- `cmd/`: Cobra command tree (grouped by domain: `auth/`, `policy/`, `tdf/`, etc.).
- `pkg/`: Core implementation used by commands (notably `pkg/handlers/`, `pkg/profiles/`, `pkg/auth/`, `pkg/utils/`).
- `docs/man/`: User-facing command docs consumed by the CLI (see `pkg/man`).
- `e2e/`: End-to-end tests written in BATS (`*.bats`) plus helper scripts.
- `tui/`: Experimental interactive UI; treat as unstable unless you’re explicitly working on it.
- `adr/`: Architecture decision records.

## Build, Test, and Development Commands

This repo is Go-first (module: `github.com/opentdf/otdfctl`) and pins a toolchain in `go.mod` (`go1.24.11`).

- `make run`: Run locally via `go run .`.
- `make test`: Run unit tests (`go test -v ./...`).
- `go test ./... -short -race -cover`: Matches CI’s unit test flags.
- `make build`: Cross-compile release artifacts into `target/` and `output/` (also runs tests and checksum steps).
- `make build-test`: Build a test-mode binary (`otdfctl_testbuild`) for local workflows.
- `make test-bats`: Run e2e BATS suite (requires `bats` installed and a running OpenTDF platform).

## Coding Style & Naming Conventions

- Go formatting is enforced: run `gofmt` (and prefer `goimports` for imports) before pushing.
- Package names: lowercase; exported identifiers: `PascalCase`; errors: `ErrX` where appropriate.
- Keep command wiring in `cmd/**` and business logic in `pkg/**` (especially `pkg/handlers/**`).
- If you add/modify commands or flags, update the matching docs in `docs/man/`.

## Testing Guidelines

- Unit tests live alongside code as `*_test.go`; keep table tests readable and deterministic.
- E2E tests are `e2e/*.bats`; be mindful of terminal sizing (`e2e/resize_terminal.sh`).

## Commit & Pull Request Guidelines

- Follow Conventional Commits as seen in history (e.g., `feat(core): …`, `fix(ci): …`, `chore(dependabot): …`; use `!` for breaking changes).
- DCO sign-off is required: `git commit -s -m "feat(core): …"`.
- PR titles are linted for semantic format: types `fix|feat|chore|docs` and scopes `main|core|tui|demo|ci|dependabot`.
- PRs should include a short description, linked issue (if any), and note any user-visible CLI output changes.
