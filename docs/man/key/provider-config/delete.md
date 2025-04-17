---
title: Delete a Provider Config
command:
  name: delete
  aliases:
    - d
    - remove
  flags:
    - name: id
      shorthand: i
      description: ID of the provider config to delete
      required: true
---

Deletes a provider config by its unique ID.

## Examples

```shell
otdfctl key provider-config delete --id <provider-config-id>
```