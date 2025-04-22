---
title: Create a Custom Action
command:
  name: create
  aliases:
    - c
    - add
    - new
  flags:
    - name: name
      shorthand: n
      description: Name of the custom action (must be unique within Policy)
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

Add a custom `action` to the platform Policy.

An action `name` is normalized to lower case and may contain hyphens or dashes between other alphanumeric
characters.

For more information, see the `actions` subcommand.

## Examples

Create a custom action named 'install_package': 

```shell
otdfctl policy actions create --name install_package
```

