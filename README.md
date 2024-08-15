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

The [tests](./tests) directory contains e2e Bash Automated Test System (bats) tests for all of the cli functionality.

To install bats on MacOS:
```
$ brew install bats-core
```
Or with NPM on any OS:
```
# To install globally:
$ npm install -g bats

# To install into your project and save it as one of the "devDependencies" in
# your package.json:
$ npm install --save-dev bats
```

These tests require the platform to be running and provisioned with basic keycloak clients/users. Before running, clone https://github.com/opentdf/platform and follow [the quickstart](https://github.com/opentdf/platform?tab=readme-ov-file#quick-start) to spin it up.

Build the cli:
```
$ go build .
```

Run the bats with:
```
$ bats tests/*.bats
```
