---
title: List Registered Resources
command:
  name: get
  aliases:
    - g
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute
---

Retrieve a registered resource along with its metadata and values.

For more general information about registered resources, see the `registered-resources` subcommand.

## Example

```shell
otdfctl policy registered-resources get --id=3c51a593-cbf8-419d-b7dc-b656d0bedfbb
```
