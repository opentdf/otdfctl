---
title: Create an attribute value
command:
  name: create
  aliases:
    - new
    - add
    - c
  flags:
    - name: attribute-id
      shorthand: a
      description: The ID of the attribute to create a value for
    - name: value
      shorthand: v
      description: The value to create
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

Add a single new value underneath an existing attribute.

For a hierarchical attribute, a new value is added in lowest hierarchy (last).

For more information on attribute values, see the `values` subcommand.

## Example

```shell
otdfctl policy attributes values create --attribute-id 3c51a593-cbf8-419d-b7dc-b656d0bedfbb --value myvalue1
```

```shell
  SUCCESS   Created values: 355743c1-c0ef-4e8d-9790-d49d883dbc7d                                                                                                                                                                                                                                                                                                                                                                              
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
