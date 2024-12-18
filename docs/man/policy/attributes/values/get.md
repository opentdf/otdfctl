---
title: Get an attribute value
command:
  name: get
  aliases:
    - g
  flags:
    - name: id
      shorthand: i
      description: The ID of the attribute value to get
---

Retrieve an attribute value along with its metadata.

For more general information about attribute values, see the `values` subcommand.

## Example

```shell
otdfctl policy attributes values get --id 355743c1-c0ef-4e8d-9790-d49d883dbc7d
```

```shell
  SUCCESS   Found values: 355743c1-c0ef-4e8d-9790-d49d883dbc7d                                                                                                                                                                                                                                                                                                                                                                     
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │355743c1-c0ef-4e8d-9790-d49d883dbc7d                                                                                                            │
│FQN                                                                      │https://opentdf.io/attr/myattribute/value/myvalue1                                                                                              │
│Value                                                                    │myvalue1                                                                                                                                        │
│Created At                                                               │Tue Dec 17 19:06:55 UTC 2024                                                                                                                    │
│Updated At                                                               │Tue Dec 17 19:06:55 UTC 2024                                                                                                                    │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy attributes values get --id=355743c1-c0ef-4e8d-9790-d49d883dbc7d --json' to see all properties
```
