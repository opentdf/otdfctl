---
title: Get a Provider Config
command:
  name: get
  aliases:
    - g
  flags:
    - name: id
      shorthand: i
      description: ID of the provider config to retrieve
    - name: name
      shorthand: n
      description: Name of the provider config to retrieve
---

Retrieves a provider config by its ID or name.

## Examples

```shell
otdfctl key provider-config get --id <provider-config-id>
```