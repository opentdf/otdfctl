---
title: Get an attribute namespace
command:
  name: get
  aliases:
    - g
  flags:
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
SUCCESS   Found namespaces: 7650f02a-be00-4faa-a1d1-37cded5e23dc
┌────────────────────────────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│Property                                                                    │Value                                                                                                │
├────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                          │7650f02a-be00-4faa-a1d1-37cded5e23dc                                                                 │
│Name                                                                        │opentdf.io                                                                                           │
│Created At                                                                  │Mon Jun 24 11:02:00 UTC 2024                                                                         │
│Updated At                                                                  │Mon Jun 24 11:02:00 UTC 2024                                                                         │
└────────────────────────────────────────────────────────────────────────────┴─────────────────────────────────────────────────────────────────────────────────────────────────────┘
NOTE   Use 'namespaces get --id=7650f02a-be00-4faa-a1d1-37cded5e23dc --json' to see all properties
```