---
title: List Public Keys
command:
  name: list
  aliases:
    - l
  flags:
    - name: kas
      shorthand: k
      description: Key Access Server ID, Name or URI.
    - name: show-public-key
      description: Show the public key
      default: false
    - name: limit
      shorthand: l
      description: Limit retrieved count
    - name: offset
      shorthand: o
      description: Offset (page) quantity from start of the list
---

List public keys shows a list of public keys.

## Example

```shell
otdfctl policy kas-registry public-key list
```

```shell
# List public keys with Key Access Server ID
otdfctl policy kas-registry public-keys list --kas=62857b55-560c-4b67-96e3-33e4670ecb3b
```

```shell
# List public keys with Key Access Server Name
otdfctl policy kas-registry public-keys list --kas=example-kas
```

```shell
# List public keys with Key Access Server URI
otdfctl policy kas-registry public-keys list --kas=https://example.com/kas
```