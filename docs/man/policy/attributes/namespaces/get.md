---
title: Get an attribute namespace
command:
  name: get
  aliases:
    - g
  flags:
    - name: fqn
      shorthand: f
      description: FQN of the attribute namespace
    - name: id
      shorthand: i
      description: ID of the attribute namespace
---

For more information, see the `namespaces` subcommand.

## Example

```shell
otdfctl policy attributes namespaces get --id=7650f02a-be00-4faa-a1d1-37cded5e23dc
```
```shell
otdfctl policy attributes namespaces get --fqn=https://opentdf.io # OpenTDF currently requires the protocol be included with the FQN
```
