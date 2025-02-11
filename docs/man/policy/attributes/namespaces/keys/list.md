---
title: List Namespace Keys
command:
  name: list
  aliases:
    - l
  flags:
    - name: namespace
      shorthand: n
      description: ID or FQN of the Attribute Namespace
      required: true
---

List the public key mappings for an attribute namespace.

## Example

```shell
# List public key mappings with Namespace ID
otdfctl policy attributes namespaces keys list --namespace=62857b55-560c-4b67-96e3-33e4670ecb3b
```

```shell
# List public key mappings with Namespace FQN
otdfctl policy attributes namespaces keys list --namespace=https://example.namespace
```