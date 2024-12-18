---
title: Deactivate an attribute value
command:
  name: deactivate
  flags:
    - name: id
      shorthand: i
      description: The ID of the attribute value to deactivate
---

Deactivation preserves uniqueness of the attribute value within policy and all existing relations, essentially reserving it.

However, a deactivation of an attribute value means it cannot be entitled in an access decision.

For information about reactivation, see the `unsafe reactivate` subcommand.

For more information on attribute values, see the `values` subcommand.

## Example

```shell
otdfctl policy attributes values deactivate --id 355743c1-c0ef-4e8d-9790-d49d883dbc7d
```

```shell
  SUCCESS   Deactivated values: 355743c1-c0ef-4e8d-9790-d49d883dbc7d                                                                                                                                                                                                                                                                                                                                                                       
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │355743c1-c0ef-4e8d-9790-d49d883dbc7d                                                                                                            │
│FQN                                                                      │https://opentdf.io/attr/myattribute/value/myvalue1                                                                                              │
│Value                                                                    │myvalue1                                                                                                                                        │
│Created At                                                               │Tue Dec 17 19:06:55 UTC 2024                                                                                                                    │
│Updated At                                                               │Tue Dec 17 19:13:38 UTC 2024                                                                                                                    │
│Labels                                                                   │[hello: world]                                                                                                                                  │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy attributes values list --json' to see all properties
```
