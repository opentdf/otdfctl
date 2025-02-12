---
title: Add a Public Key
command:
  name: add
  aliases:
    - a
  flags:
    - name: public-key-id
      shorthand: i
      description: ID of the Public Key
      required: true
    - name: value
      shorthand: v
      description: ID or FQN of the Attribute Value
      required: true
---

Add a public key mapping to an attribute value.

## Example

```shell
otdfctl policy attributes values keys add --public-key-id=62857b55-560c-4b67-96e3-33e4670ecb3b --value=62857b55-560c-4b67-96e3-33e4670ecb3b
```

```shell
otdfctl policy attributes values keys add --public-key-id=62857b55-560c-4b67-96e3-33e4670ecb3b --value=https://example.com/attr/attr1/value/val1
```
