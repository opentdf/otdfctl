---
title: List an obligation definition
command:
  name: list
  aliases:
    - l
  flags:
    - name: limit
      shorthand: l
      description: Limit retrieved count
    - name: offset
      shorthand: o
      description: Offset (page) quantity from start of the list
    - name: namespace
      shorthand: n
      description: Namespace ID or FQN by which to filter results
---

For more information about Obligations, see the `obligations` subcommand.

## Example

```shell
otdfctl policy obligations list
```
