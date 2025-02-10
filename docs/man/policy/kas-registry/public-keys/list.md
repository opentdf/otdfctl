---
title: List Public Keys
command:
  name: list
  aliases:
    - l
  flags:
    - name: kas
      shorthand: k
      description: Key Access Server ID, Name or URI.
    - name: show-public-key
      description: Show the public key
      default: false
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
otdfctl policy kas-registry public-key list
```
