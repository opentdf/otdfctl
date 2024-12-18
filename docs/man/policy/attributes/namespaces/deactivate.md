---
title: Deactivate an attribute namespace
command:
  name: deactivate
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute namespace
      required: true
    - name: force
      description: Force deactivation without interactive confirmation (dangerous)
---

Deactivating an Attribute Namespace will make the namespace name inactive as well as any attribute definitions and values beneath.

Deactivation of a Namespace renders any existing TDFs of those attributes inaccessible.

Deactivation will permanently reserve the Namespace name within a platform. Reactivation and deletion are both considered "unsafe"
behaviors.

For information about reactivation, see the `unsafe reactivate` subcommand.

For reactivation, see the `unsafe` command.

## Example 

```shell
otdfctl policy attributes namespaces deactivate --id 7650f02a-be00-4faa-a1d1-37cded5e23dc
```

```shell
SUCCESS   Deactivated namespaces: 7650f02a-be00-4faa-a1d1-37cded5e23dc
┌────────────────────────────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────────────────────────────────────────────────┐
│Property                                                                    │Value                                                                                                │
├────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Name                                                                        │opentdf.io                                                                                           │
│Id                                                                          │7650f02a-be00-4faa-a1d1-37cded5e23dc                                                                 │
│Created At                                                                  │Mon Jun 24 11:02:00 UTC 2024                                                                         │
│Updated At                                                                  │Mon Jun 24 11:02:00 UTC 2024                                                                         │
└────────────────────────────────────────────────────────────────────────────┴─────────────────────────────────────────────────────────────────────────────────────────────────────┘
NOTE   Use 'otdfctl policy attributes namespaces list --json' to see all properties
```
