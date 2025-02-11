---
title: Get a Public Key
command:
  name: get
  aliases:
    - g
  flags:
    - name: id
      shorthand: i
      description: ID of the Public Key
      required: true
---

Get a public key.

## Example

```shell
otdfctl policy kas-registry public-key get --id=62857b55-560c-4b67-96e3-33e4670ecb3b
```
