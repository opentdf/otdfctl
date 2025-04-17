---
title: Update a Provider Config
command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      shorthand: i
      description: ID of the provider config to update
      required: true
    - name: name
      shorthand: n
      description: New name for the provider config
    - name: config
      shorthand: c
      description: New JSON configuration for the provider
    - name: label
      shorthand: l
      description: Metadata labels for the provider config
---

Updates an existing provider config with the specified parameters.

## Examples

```shell
otdfctl key provider-config update --id <id> --name <new-name> --config <new-json-config>
```