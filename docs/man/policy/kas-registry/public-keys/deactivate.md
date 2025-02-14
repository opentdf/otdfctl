---
title: Deactivate a Public Key
command:
  name: deactivate
  aliases:
    - d
  flags:
    - name: id
      shorthand: i
      description: ID of the Public Key
      required: true
---

Deactivate a public key.

## Example

```shell
otdfctl policy kas-registry public-keys deactivate --id=62857b55-560c-4b67-96e3-33e4670ecb3b
```
