---
title: Deactivate an attribute definition
command:
  name: deactivate
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute
      required: true
    - name: force
      description: Force deactivation without interactive confirmation (dangerous)
---

Deactivation preserves uniqueness of the attribute and values underneath within policy and all existing relations,
essentially reserving them.

However, a deactivation of an attribute means its associated values cannot be entitled in an access decision.

For information about reactivation, see the `unsafe reactivate` subcommand.

For more general information about attributes, see the `attributes` subcommand.

## Example

```shell
otdfctl policy attributes deactivate --id 3c51a593-cbf8-419d-b7dc-b656d0bedfbb
```

```shell
  SUCCESS   Deactivated attributes: 3c51a593-cbf8-419d-b7dc-b656d0bedfbb                                                                                                                                                                                                                                                                                                                                    
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Name                                                                     │myattribute                                                                                                                                     │
│Rule                                                                     │ANY_OF                                                                                                                                          │
│Values                                                                   │[myvalue1]                                                                                                                                      │
│Namespace                                                                │opentdf.io                                                                                                                                      │
│Created At                                                               │Tue Dec 17 18:33:06 UTC 2024                                                                                                                    │
│Updated At                                                               │Tue Dec 17 19:41:47 UTC 2024                                                                                                                    │
│Labels                                                                   │[hello: world]                                                                                                                                  │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy attributes list --json' to see all properties 
```
