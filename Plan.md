# Plan for Interactive Shell

This document outlines the plan for creating an interactive shell for `otdfctl`.

## Step 1: Create the Basic Shell Structure

*   Create a new `shell.go` file in the `tui` package.
*   Implement a basic `bubbletea` model with a single input for commands.
*   The shell will have a prompt that displays the current context (e.g., `platform:/>`).

## Step 2: Implement `ls` and `cd` Commands

*   Implement the logic for navigating the hierarchical structure we discussed:
    ```
    /
    ├── namespaces/
    │   ├── <namespace-name>/
    │   │   ├── attribute-definitions/
    │   │   │   ├── <attribute-definition-name>/
    │   │   │   │   ├── attribute-values/
    │   │   │   │   │   ├── <attribute-value-name>
    │   │   │   │   │   └── ...
    │   │   │   │   └── ...
    │   │   ├── obligations/
    │   │   │   ├── <obligation-name>/
    │   │   │   │   ├── obligation-values/
    │   │   │   │   │   ├── <obligation-value-name>
    │   │   │   │   │   └── ...
    │   │   │   │   └── ...
    │   │   │   └── ...
    │   │   └── ...
    └── registered-resources/
        ├── <registered-resource-name>/
        │   ├── registered-resource-values/
        │   │   ├── <registered-resource-value-name>
        │   │   └── ...
        │   └── ...
        └── ...
    ```
*   The `ls` command will list the available items in the current context.
*   The `cd` command will change the current context.

## Step 3: Implement the `help` Command

*   The `help` command will display a list of available commands and their usage.

## Step 4: Implement Policy Management Commands

*   Map the existing `otdfctl` policy commands to the interactive shell.
*   Start with attribute management and key/KAS assignment, as requested.
*   The commands will be context-aware. For example, when in the context of an attribute definition, the `assign-key` command will know which attribute definition to assign the key to.

## Step 5: Implement Authentication

*   Integrate the interactive shell with the existing `otdfctl` authentication mechanisms.
*   The shell will use the active profile to authenticate with the platform.

## Step 6: Add Advanced Features

*   Implement command history, allowing the user to cycle through previous commands.
*   Implement autocomplete to suggest commands and arguments based on the current context.