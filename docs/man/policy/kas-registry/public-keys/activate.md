---
title: Activate a Public Key
command:
  name: activate
  aliases:
    - a
  flags:
    - name: id
      shorthand: i
      description: ID of the Public Key
      required: true
---

Activate a public key.

## Example

```shell
otdfctl policy kas-registry public-keys activate --id=62857b55-560c-4b67-96e3-33e4670ecb3b
```
