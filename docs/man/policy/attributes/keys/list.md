---
title: List Attribute Keys
command:
  name: list
  aliases:
    - l
  flags:
    - name: definition
      shorthand: d
      description: ID or FQN of the Attribute Definition
      required: true
    - name: show-public-key
      description: Show the public key
      default: false
---

List the public key mappings for an attribute definition.

## Example

```shell
# List public key mappings with Definition ID
otdfctl policy attributes definitions keys list --definition=62857b55-560c-4b67-96e3-33e4670ecb3b
```

```shell
# List public key mappings with Definition FQN
otdfctl policy attributes definitions keys list --definition=https://example.com/attr/attr1
```