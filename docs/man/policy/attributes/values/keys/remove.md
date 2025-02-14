---
title: Remove a Public Key
command:
  name: remove
  aliases:
    - r
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

Remove a public key mapping from an attribute value.

## Example

```shell
# Remove a public key mapping with Value ID
otdfctl policy attributes values keys remove --public-key-id=62857b55-560c-4b67-96e3-33e4670ecb3b --value=62857b55-560c-4b67-96e3-33e4670ecb3b
```

```shell
# Remove a public key mapping with Value FQN
otdfctl policy attributes values keys remove --public-key-id=62857b55-560c-4b67-96e3-33e4670ecb3b --value=https://example.com/attr/attr1/value/val1
```

