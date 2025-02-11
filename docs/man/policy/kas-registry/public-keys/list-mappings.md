---
title: List Public Key Mappings
command:
  name: list-mappings
  aliases:
    - lm
  flags:
    - name: kas
      shorthand: k
      description: Key Access Server ID, Name or URI.
    - name: public-key-id
      shorthand: p
      description: Public Key ID
    - name: limit
      shorthand: l
      description: Limit retrieved count
    - name: offset
      shorthand: o
      description: Offset (page) quantity from start of the list
---

List public key mappings shows a list of Key Access Servers and associated public keys. Each Key also has a list of namespaces, attribute defnitions and attribute values that are associated with it.

## Example

```shell
# List public key mappings
otdfctl policy kas-registry public-keys list-mappings
```

```shell
# List public key mappings with Key Access Server ID
otdfctl policy kas-registry public-keys list-mappings --kas=62857b55-560c-4b67-96e3-33e4670ecb3b
```

```shell
# List public key mappings with Key Access Server Name
otdfctl policy kas-registry public-keys list-mappings --kas=example-kas
```

```shell
# List public key mappings with Key Access Server URI
otdfctl policy kas-registry public-keys list-mappings --kas=https://example.com/kas
```

```shell
# List public key mappings with Public Key ID
otdfctl policy kas-registry public-keys list-mappings --public-key-id=62857b55-560c-4b67-96e3-33e4670ecb3b
```
