# Building and Testing otdfctl

This document covers how to build, test, and verify changes to the otdfctl CLI.

## Building the Project

### Standard Build
```bash
go build
```

This creates an `otdfctl` binary in the current directory.

### Build with Version Info
The Makefile injects version information at build time:

```bash
make build
```

This sets:
- `Version` - from `SEM_VER` env var or git tags
- `CommitSha` - from `COMMIT_SHA` env var or git
- `BuildTime` - current UTC timestamp

These are injected into `pkg/config` via `-ldflags`.

### Test Mode Build
For testing with in-memory profiles and test provisioning:

```bash
make build-test
```

This creates `otdfctl_testbuild` with:
- `config.TestMode = true`
- In-memory keyring provider
- Support for `OTDFCTL_TEST_PROFILE` environment variable

### Cross-Platform Builds
```bash
make build
```

Builds for all platforms:
- darwin-amd64, darwin-arm64
- linux-amd64, linux-arm, linux-arm64
- windows-amd64, windows-arm, windows-arm64

Output goes to `target/` directory and is zipped to `output/` with checksums.

## Testing

### Unit Tests
```bash
go test -v ./...
# or
make test
```

Runs all Go unit tests in the project.

### BATS (End-to-End) Tests

**Prerequisites:**
1. Install BATS - Follow [bats-core installation](https://github.com/bats-core/homebrew-bats-core)
2. Platform must be running and provisioned
   - See [platform README](https://github.com/opentdf/platform) for setup

**Run BATS tests:**
```bash
make test-bats
```

This:
1. Builds test binary (`make build-test`)
2. Sets terminal to standard defaults
3. Runs all BATS tests in `e2e/`

**Run specific test suite:**
```bash
bats e2e/<test-name>.bats
```

#### Terminal Size for Tests

Terminal size affects output rendering tests. Control it via:

1. **Standard defaults**: `make test-bats` (automatic)
2. **Manual**: Resize terminal window with mouse
3. **Script**: `./e2e/resize_terminal.sh <rows> <columns>`
4. **Environment**: `export TEST_TERMINAL_WIDTH="200"`

### TestRail Integration (Optional)

Upload BATS results to TestRail for tracking.

**Setup:**
1. Copy config: `cp testrail.config.example.json testrail.config.json`
2. Edit `testrail.config.json` with your TestRail URL, project ID, TAP file path
3. Copy mapping: `cp testname-to-testrail-id.example.json testname-to-testrail-id.json`
4. Map test names to TestRail case IDs (flat or nested JSON)
5. Set credentials:
   ```bash
   export TESTRAIL_USER=you@example.com
   export TESTRAIL_PASS=your_api_key
   ```

**Upload results:**
```bash
bats --tap bats-tests/ > e2e/bats-results.tap
TESTRAIL_CLI_RUN_NAME=optional-name ./testrail-integration/upload-bats-test-results-to-testrail.sh
```

## Makefile Targets Reference

| Target | Description |
|--------|-------------|
| `make run` | Run project with `go run .` |
| `make test` | Run Go unit tests |
| `make build-test` | Build test binary with TestMode |
| `make test-bats` | Build test binary and run BATS tests |
| `make build` | Full cross-platform build with tests |
| `make clean` | Remove `target/` directory |

## Build Configuration

Build configuration is defined in `Makefile:1-96`:

- Binary name: `otdfctl`
- Go module: `github.com/opentdf/otdfctl`
- Config package: `github.com/opentdf/otdfctl/pkg/config`
- Target directory: `target/`
- Output directory: `output/`

The version info is accessible at runtime:
```go
import "github.com/opentdf/otdfctl/pkg/config"

config.Version    // e.g., "0.28.0"
config.CommitSha  // e.g., "1b5ba79..."
config.BuildTime  // e.g., "2025-12-18T10:20:30Z"
```
