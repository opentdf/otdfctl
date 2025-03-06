# otdfctl: cli to manage OpenTDF Platform

This command line interface is used to manage OpenTDF Platform.

The main goals are to:

- simplify setup
- facilitate migration
- aid in configuration management

## TODO list

- [ ] Add support for json input as piped input
- [ ] Add help level handler for each command
- [ ] Add support for `--verbose` persistent flag
- [ ] Helper functions to support common tasks like pretty printing and json output

## Usage

The CLI is configured via the `otdfctl.yaml`. There is an example provided in `otdfctl-example.yaml`.

Run `cp otdfctl-example.yaml otdfctl.yaml` to copy the example config when running the CLI.

Load up the platform (see its [README](https://github.com/opentdf/platform?tab=readme-ov-file#run) for instructions).

## Development

### CLI

The CLI is built using [cobra](https://cobra.dev/).

The primary function is to support CRUD operations using commands as arguments and flags as the values.

The output format (currently `styled` or `json`) is configurable in the `otdfctl.yaml` or via CLI flag.

#### To add a command

1. Capture the flag value and validate the values
   1. Alt support JSON input as piped input
2. Run the handler which is located in `pkg/handlers` and pass the values as arguments
3. Handle any errors and return the result in a lite TUI format

### TUI

> [!CAUTION]
> This is a work in progress please avoid touching until framework is defined

The TUI will be used to create an interactive experience for the user.

## Documentation

Documentation drives the CLI in this project. This can be found in `/docs/man` and is used in the
CLI via the `man.Docs.GetDoc()` function.

## Testing

The CLI is equipped with a test mode that can be enabled by building the CLI with `config.TestMode = true`.
For convenience, the CLI can be built with `make build-test`.

**Test Mode features**:

- Use the in-memory keyring provider for user profiles
- Enable provisioning profiles for testing via `OTDFCTL_TEST_PROFILE` environment variable

### BATS

> [!NOTE]
> Bat Automated Test System (bats) is a TAP-compliant testing framework for Bash. It provides a simple way to verify that the UNIX programs you write behave as expected.

BATS is used to test the CLI from an end-to-end perspective. To run the tests you will need to ensure the following
prerequisites are met:

- bats is installed on your system
  - MacOS: `brew install bats-core bats-support bats-assert`
- The platform is running and provisioned with basic keycloak clients/users
  - See the [platform README](https://github.com/opentdf/platform) for instructions

To run the tests you can either run `make test-bats` or execute specific test suites with `bats e2e/<test>.bats`.

#### Terminal Size

Some tests for output rendered in the terminal will vary in behavior depending on terminal size.

Terminal size when testing:

1. set to standard defaults if running `make test-bats`
2. can be set manually by mouse in terminal where tests are triggered
3. can be set by argument `./e2e/resize_terminal.sh < rows height > < columns width >`
4. can be set by environment variable, i.e. `export TEST_TERMINAL_WIDTH="200"` (200 is columns width)

## Status

In deevelopment.


I think your an nice friend.
