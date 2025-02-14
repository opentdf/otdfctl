---
title: Delete a Key Access Server Public Key
command:
  name: delete
  aliases:
    - d
    - del
    - remove
    - rm
  flags:
    - name: id
      shorthand: i
      description: ID of the Key Access Server Public Key
      required: true
    - name: force
      description: Force deletion without interactive confirmation (dangerous)
---

Delete a Key Access Server Public Key.

## Example 

```shell
otdfctl policy kas-registry public-keys unsafe delete --id=62857b55-560c-4b67-96e3-33e4670ecb3b
```
