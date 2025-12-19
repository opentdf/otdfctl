# Development Workflow

This document covers how to add new features, commands, and maintain the otdfctl codebase.

## Adding a New Command

Follow these steps when adding a new CLI command:

### 1. Create Command Documentation

Create a markdown file in `/docs/man/` with frontmatter defining the command structure.

**Example**: `/docs/man/policy_attributes_create.md`

```markdown
---
command: create
parent: policy-attributes
description: Create a new attribute
flags:
  - name: name
    type: string
    required: true
    description: Attribute name
  - name: namespace
    type: string
    required: true
    description: Namespace ID
---

# Create Attribute

Creates a new attribute in the specified namespace.

## Usage

otdfctl policy attributes create --name <name> --namespace <namespace-id>

## Examples

Create an attribute named "classification":
otdfctl policy attributes create --name classification --namespace <ns-id>
```

The frontmatter drives the command structure and help text via `man.Docs.GetDoc()`.

### 2. Create Handler Function

Add business logic in `/pkg/handlers/`.

**Example**: `/pkg/handlers/attribute.go`

```go
func CreateAttribute(cmd *cobra.Command, name, namespaceID string) error {
    c := cli.New(cmd, nil)

    // Initialize SDK
    h := NewHandler(c)
    defer h.Close()

    // Call SDK
    attr, err := h.sdk.Attributes.CreateAttribute(c.Context(), &policy.CreateAttributeRequest{
        Name:        name,
        NamespaceId: namespaceID,
    })
    if err != nil {
        return err
    }

    // Print result
    c.PrintSuccess("Attribute created", attr)
    return nil
}
```

Key patterns:
- Use `cli.New(cmd, args)` to get CLI context
- Initialize SDK with `NewHandler(c)` and defer `Close()`
- Use `c.Context()` for SDK calls
- Use `c.PrintSuccess()` or `c.ExitWith()` for output

### 3. Create Command Definition

Add command in `/cmd/` following the command hierarchy.

**Example**: `/cmd/policy/attributes.go`

```go
var createCmd = man.Docs.GetCommand("policy-attributes-create",
    man.WithRun(func(cmd *cobra.Command, args []string) {
        c := cli.New(cmd, args)

        // Get and validate flags
        name := c.Flags.GetRequiredString("name")
        namespace := c.Flags.GetRequiredString("namespace")

        // Call handler
        if err := handlers.CreateAttribute(cmd, name, namespace); err != nil {
            c.ExitWithError("Failed to create attribute", err)
        }
    }),
)

func init() {
    attributesCmd.AddCommand(&createCmd.Command)
}
```

Key patterns:
- Use `man.Docs.GetCommand()` to load command from docs
- Use `man.WithRun()` to set the command function
- Use `cli.New()` to get CLI context
- Use `c.Flags.GetRequiredString()` for required flags
- Use `c.Flags.GetOptionalString()` for optional flags
- Use `c.ExitWithError()` for error handling

### 4. Register Command

Add the command to the parent command in its `init()` function.

See `cmd/policy/policy.go` for example of command registration.

## Output Handling

The CLI supports two output formats: `styled` (pretty terminal output) and `json`.

### Using the Printer

```go
c := cli.New(cmd, args)

// Success with data
c.PrintSuccess("Operation completed", result)

// Exit with error
c.ExitWithError("Operation failed", err)

// Custom exit with specific code
c.ExitWith(message, data, exitCode, writer)
```

### Output Format Selection

- **Per-profile**: Set via `otdfctl profile create --output-format <styled|json>`
- **Per-command**: Override with `--json` flag
- **Storage**: Format stored in profile config (see `pkg/profiles/`)

The printer automatically handles formatting based on the selected format.

## Error Handling

### CLI Errors

Use the `cli` package error helpers:

```go
c := cli.New(cmd, args)

// Exit with error message
c.ExitWithError("Failed to connect", err)

// Exit with just a message
c.ExitWithMessage("Invalid configuration", cli.ExitCodeError)
```

### Handler Errors

Return errors from handlers; commands handle the exit:

```go
func MyHandler(cmd *cobra.Command) error {
    // ... do work ...
    if err != nil {
        return fmt.Errorf("failed to process: %w", err)
    }
    return nil
}
```

### Custom Error Types

See `pkg/cli/errors.go` for CLI error types and `pkg/auth/errors.go` for auth-specific errors.

## Profile Management

Commands that need to connect to the platform use profiles.

### Loading a Profile

```go
c := cli.New(cmd, args)

// Load profile
profile, err := profiles.GetCurrent()
if err != nil {
    c.ExitWithError("Failed to load profile", err)
}

// Use profile
endpoint := profile.Endpoint
```

### Profile Structure

Profiles store:
- Endpoint URL
- Authentication credentials (in OS keyring)
- Output format preference
- Other connection settings

See `pkg/profiles/profileConfig.go` for structure.

## SDK Integration

The OpenTDF Platform SDK is used for all platform interactions.

### Initializing the SDK

```go
h := handlers.NewHandler(c)
defer h.Close()

// Use SDK
result, err := h.sdk.SomeService.SomeMethod(c.Context(), request)
```

### SDK Structure

The SDK provides services for:
- `Attributes` - Attribute management
- `Namespaces` - Namespace management
- `KAS` - Key Access Service operations
- `TDF` - TDF operations
- And more...

See `pkg/handlers/sdk.go` for SDK initialization.

## Validation

### Flag Validation

Use the `utils` package validators:

```go
import "github.com/opentdf/otdfctl/pkg/utils"

// Validate UUID
if err := utils.ValidateUUID(id); err != nil {
    return fmt.Errorf("invalid UUID: %w", err)
}

// Validate URL
if err := utils.ValidateURL(endpoint); err != nil {
    return fmt.Errorf("invalid URL: %w", err)
}
```

See `pkg/utils/validators.go` for available validators.

### Input Validation

Validate inputs early in the command function before calling handlers:

```go
name := c.Flags.GetRequiredString("name")
if name == "" {
    c.ExitWithMessage("Name is required", cli.ExitCodeError)
}

if len(name) > 255 {
    c.ExitWithMessage("Name too long (max 255 chars)", cli.ExitCodeError)
}
```

## Piped Input Support

Commands can accept JSON input via stdin.

**TODO**: This feature is on the roadmap but not yet fully implemented. See README.md:11.

Planned pattern:
```bash
echo '{"name": "test"}' | otdfctl policy attributes create
```

See `pkg/cli/pipe.go` for pipe handling utilities.

## Testing Your Changes

### 1. Unit Tests

Add tests alongside your code:

```go
// pkg/handlers/attribute_test.go
func TestCreateAttribute(t *testing.T) {
    // ... test implementation ...
}
```

Run tests:
```bash
go test -v ./pkg/handlers/
```

### 2. Manual Testing

Build and test locally:

```bash
go build
./otdfctl <your-command> <flags>
```

### 3. Integration Tests

Add BATS tests in `/e2e/`:

```bash
# e2e/test_attributes.bats
@test "create attribute" {
    run otdfctl policy attributes create --name test --namespace $NS_ID
    assert_success
    assert_output --partial "Attribute created"
}
```

Run BATS tests:
```bash
make test-bats
```

## Code Conventions

### Follow Existing Patterns

The codebase follows consistent patterns. Before implementing, look at similar existing commands:

- Attributes: `cmd/policy/attributes.go`, `pkg/handlers/attribute.go`
- Auth commands: `cmd/auth/*.go`
- TDF operations: `cmd/tdf/*.go`, `pkg/handlers/tdf.go`

### In-Context Learning

The LLM learns from existing code patterns. Following existing conventions ensures consistency without explicit style guidelines.

### Linting and Formatting

Run linters before committing:

```bash
golangci-lint run
```

The project uses `.golangci.yaml` for linter configuration.

### Git Workflow

Standard git workflow:

1. Create feature branch
2. Make changes
3. Test locally
4. Commit with descriptive message
5. Push and create PR

**Important**: Never use interactive git commands (`git rebase -i`, `git add -i`, etc.) - they're not supported in the agent environment.

## Common Tasks

### Adding a Flag to Existing Command

1. Update documentation in `/docs/man/*.md` frontmatter
2. Access flag in command: `c.Flags.GetOptionalString("newflag")`
3. Pass to handler if needed

### Adding a New Policy Resource Type

1. Add handler in `/pkg/handlers/<resource>.go`
2. Add command in `/cmd/policy/<resource>.go`
3. Add documentation in `/docs/man/policy_<resource>_*.md`
4. Register command in `/cmd/policy/policy.go`

### Modifying Output Format

1. Update handler to return correct structure
2. Ensure structure serializes to JSON properly
3. Update printer usage if needed (`c.PrintSuccess()`)

## References

- **README.md:34-39** - Command addition steps
- **README.md:48-51** - Documentation system
- **Cobra docs** - https://cobra.dev/
- **OpenTDF Platform** - https://github.com/opentdf/platform
