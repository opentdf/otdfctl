---
title: Remove a KAS key from an attribute namespace
command:
  name: remove
  flags:
    - name: namespace
      shorthand: n
      description: Can be URI or ID of namespace
      required: true
    - name: keyId
      shorthand: k
      description: ID of the KAS key to remove
      required: true
---

Removes a KAS key from a policy attribute namespace. After removing the key, the attribute namespace can no longer be used with the specified KAS key for encryption and decryption operations.

## Example

```shell
otdfctl policy attributes namespaces remove --namespace 3d25d33e-2469-4990-a9ed-fdd13ce74436 --keyId 8f7e6d5c-4b3a-2d1e-9f8d-7c6b5a432f1d
```

```shell
otdfctl policy attributes namespaces remove --namespace "https://example.com" --keyId 8f7e6d5c-4b3a-2d1e-9f8d-7c6b5a432f1d
```
