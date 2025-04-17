---
title: Create a Provider Config
command:
  name: create
  aliases:
    - c
  flags:
    - name: name
      shorthand: n
      description: Name of the provider config to create
      required: true
    - name: config
      shorthand: c
      description: JSON configuration for the provider
      required: true
    - name: label
      shorthand: l
      description: Metadata labels for the provider config
---

Creates a new provider config with the specified name and configuration.

## Examples

```shell
otdfctl key provider-config create --name <name> --config <json-config>
```