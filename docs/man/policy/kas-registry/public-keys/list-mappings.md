---
title: List Public Key Mappings
command:
  name: list-mappings
  aliases:
    - lm
  flags:
    - name: kas
      shorthand: k
      description: Key Access Server ID, Name or URI.
    - name: public-key-id
      shorthand: p
      description: Public Key ID
    - name: limit
      shorthand: l
      description: Limit retrieved count
    - name: offset
      shorthand: o
      description: Offset (page) quantity from start of the list
---

For more information about registration of Key Access Servers, see the manual for `kas-registry`.

## Example

```shell
otdfctl policy kas-registry list
```
