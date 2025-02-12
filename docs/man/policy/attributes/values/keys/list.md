---
title: List Key Mappings
command:
  name: list
  aliases:
    - l
  flags:
    - name: value
      shorthand: v
      description: ID or FQN of the Attribute Value
      required: true
    - name: show-public-key
      description: Show the public key
      default: false
---

List the public key mappings for an attribute value.

## Example

```shell
# List public key mappings with Value ID
otdfctl policy attributes values keys list --value=62857b55-560c-4b67-96e3-33e4670ecb3b
```

```shell
# List public key mappings with Value FQN
otdfctl policy attributes values keys list --value=https://example.com/attr/attr1/value/val1
```