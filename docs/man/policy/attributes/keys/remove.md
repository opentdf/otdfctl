---
title: Remove a Public Key Mapping
command:
  name: remove
  aliases:
    - r
  flags:
    - name: public-key-id
      shorthand: i
      description: ID of the Public Key
      required: true
    - name: definition
      shorthand: d
      description: ID or FQN of the Attribute Definition
      required: true
---

Remove a public key mapping from an attribute definition.

## Example

```shell
otdfctl policy attributes definitions keys remove --public-key-id=62857b55-560c-4b67-96e3-33e4670ecb3b --definition=62857b55-560c-4b67-96e3-33e4670ecb3b
```

```shell
otdfctl policy attributes definitions keys remove --public-key-id=62857b55-560c-4b67-96e3-33e4670ecb3b --definition=https://example.com/attr/attr1
```
