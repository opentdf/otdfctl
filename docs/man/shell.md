---
title: Interactive Shell

command:
  name: shell
  aliases:
    - sh
---

# Interactive Shell

Launch an interactive shell that provides a filesystem-like interface for navigating and managing OpenTDF Platform resources.

## Usage

```bash
otdfctl shell
```

## Overview

The shell provides an interactive environment where you can navigate through platform resources using familiar commands like `cd` and `ls`. Resources are organized in a hierarchical structure:

```
/
├── namespaces/
│   └── <namespace-name>/
│       ├── attribute-definitions/
│       │   └── <attribute-name>/
│       │       └── attribute-values/
│       │           └── <value-name>
│       └── ...
└── registered-resources/
    └── <resource-name>/
        └── ...
```

## Available Commands

### Navigation

- **`ls`** - List items in the current context
- **`cd <path>`** - Change to the specified path
  - `cd /` - Go to root
  - `cd ..` - Go up one level
  - `cd namespaces` - Enter namespaces directory
  - `cd namespaces/example.com` - Navigate to specific namespace
- **`pwd`** - Print current working directory

### Information

- **`help`** - Show available commands for current context
- **`clear`** - Clear the output display

### Profile Management

- **`profile`** - Show the current active profile
- **`profile list`** - List all available profiles (shows current with `>` and default with `*`)
- **`profile use <name>`** - Switch to a different profile (resets navigation to root)

### Resource Details

- **`get`** or **`show`** - Display detailed information about the current resource
  - Shows comprehensive details including IDs, names, status, and related data
  - Only available when navigated to a specific resource (namespace, attribute, value)
  - Not available at collection levels (e.g., `/namespaces/` or `/attribute-definitions/`)

### Resource Management

Commands are context-aware based on your current location:

- **`create`** - Create a new resource (prompts for required information) *(coming soon)*
- **`update`** - Update the current resource *(coming soon)*
- **`delete`** - Delete the current resource *(coming soon)*

### Shell Control

- **`exit`** or **`quit`** - Exit the shell
- **Ctrl+C** - Exit the shell

### Command History

The shell maintains a history of executed commands that you can navigate:

- **Up Arrow** - Navigate to previous commands in history
- **Down Arrow** - Navigate to next commands in history (or return to current input)
- History navigation preserves your current input if you haven't executed it yet
- Duplicate consecutive commands are automatically filtered from history
- History is maintained for the duration of your shell session

## Context-Aware Operations

The shell understands your current location and automatically provides context to commands. For example:

```bash
platform:/> cd namespaces/example.com/attribute-definitions
platform:/namespaces/example.com/attribute-definitions/> create
# Will prompt you to create a new attribute in the example.com namespace
# No need to specify the namespace - it's inferred from your location
```

## Examples

### Navigate to a namespace

```bash
platform:/> ls
namespaces/
registered-resources/

platform:/> cd namespaces
platform:/namespaces/> ls
example.com
default
...

platform:/namespaces/> cd example.com
platform:/namespaces/example.com/> ls
attribute-definitions/
```

### View attributes in a namespace

```bash
platform:/namespaces/example.com/> cd attribute-definitions
platform:/namespaces/example.com/attribute-definitions/> ls
classification
clearance
...
```

### Create a new attribute

```bash
platform:/namespaces/example.com/attribute-definitions/> create
# Interactive prompts will guide you through attribute creation
```

### Switch profiles

```bash
platform:/> profile list
Available profiles:

> production (current)
* staging (default)
  development

platform:/> profile use development
✓ Switched to profile: development

development:/> profile
Current profile: development
```

### View resource details

```bash
platform:/namespaces/> cd example.com
platform:/namespaces/example.com/> get

Namespace Details

Name: example.com
ID: abc-123-def-456
FQN: https://example.com
Active: true

platform:/namespaces/example.com/> cd attribute-definitions/classification
platform:/namespaces/example.com/attribute-definitions/classification/> get

Attribute Definition Details

Name: classification
ID: attr-123-456
Namespace: example.com
Rule: RULE_TYPE_ALL_OF
Active: true

Values: 3
  • secret
  • top-secret
  • unclassified
```

## Authentication

The shell uses your current profile for authentication. Make sure you're logged in before launching the shell:

```bash
otdfctl auth login
otdfctl shell
```

## Tips

- Use **Tab** to autocomplete commands and paths
- Use **Up/Down arrows** to navigate through command history (like bash/zsh)
- Use **`help`** at any level to see available commands
- Use **`get`** or **`show`** to inspect detailed information about any resource
- Use **`profile list`** and **`profile use <name>`** to switch between profiles without exiting the shell
- Switching profiles automatically resets navigation to root
- Paths can be absolute (`/namespaces/example.com`) or relative (`../..`)
- The `get` command is only available when you're at a specific resource (not at collection levels)
- Command history avoids storing duplicate consecutive commands
