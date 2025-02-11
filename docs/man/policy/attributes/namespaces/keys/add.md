---
title: Add a Public Key Mapping
command:
  name: add
  aliases:
    - a
  flags:
    - name: public-key-id
      shorthand: i
      description: ID of the Public Key
      required: true
    - name: namespace
      shorthand: n
      description: ID or FQN of the Attribute Namespace
      required: true
---

Add a public key mapping to an attribute namespace.

## Example

```shell
# Add a public key mapping with Namespace ID
otdfctl policy attributes namespaces keys add --public-key-id=62857b55-560c-4b67-96e3-33e4670ecb3b --namespace=62857b55-560c-4b67-96e3-33e4670ecb3b
```

```shell
# Add a public key mapping with Namespace FQN
otdfctl policy attributes namespaces keys add --public-key-id=62857b55-560c-4b67-96e3-33e4670ecb3b --namespace=https://example.namespace
```
