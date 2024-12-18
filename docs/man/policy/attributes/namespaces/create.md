---
title: Create an attribute namespace
command:
  name: create
  aliases:
    - c
    - add
    - new
  flags:
    - name: name
      shorthand: n
      description: Name of the attribute namespace
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

Creation of a `namespace` is required to add attributes or any other policy objects beneath.

For more information, see the `namespaces` subcommand.

## Example

```shell
otdfctl policy attributes namespaces create --name opentdf.io
```

```shell
SUCCESS   Created namespaces: 7650f02a-be00-4faa-a1d1-37cded5e23dc
┌────────────────────────────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│Property                                                                    │Value                                                                                                │
├────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Name                                                                        │opentdf.io                                                                                           │
│Id                                                                          │7650f02a-be00-4faa-a1d1-37cded5e23dc                                                                 │
│Created At                                                                  │Mon Jun 24 11:02:00 UTC 2024                                                                         │
│Updated At                                                                  │Mon Jun 24 11:02:00 UTC 2024                                                                         │
└────────────────────────────────────────────────────────────────────────────┴─────────────────────────────────────────────────────────────────────────────────────────────────────┘
NOTE   Use 'namespaces get --id=7650f02a-be00-4faa-a1d1-37cded5e23dc --json' to see all properties
```
