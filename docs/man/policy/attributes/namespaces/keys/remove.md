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
    - name: namespace
      shorthand: d
      description: ID or FQN of the Attribute Namespace
      required: true
---

Remove a public key mapping from an attribute namespace.

## Example

```shell
# Remove a public key mapping with Namespace ID
otdfctl policy attributes namespaces keys remove --public-key-id=62857b55-560c-4b67-96e3-33e4670ecb3b --namespace=62857b55-560c-4b67-96e3-33e4670ecb3b
```

```shell
# Remove a public key mapping with Namespace FQN
otdfctl policy attributes namespaces keys remove --public-key-id=62857b55-560c-4b67-96e3-33e4670ecb3b --namespace=https://example.namespace
```
