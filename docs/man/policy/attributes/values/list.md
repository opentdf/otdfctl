---
title: List attribute values
command:
  name: list
  aliases:
    - ls
    - l
  flags:
    - name: attribute-id
      shorthand: a
      description: The ID of the attribute to list values for
    - name: state
      shorthand: s
      description: Filter by state
      enum:
        - active
        - inactive
        - any
      default: active
    - name: limit
      shorthand: l
      description: Limit retrieved count
    - name: offset
      shorthand: o
      description: Offset (page) quantity from start of the list
---

By default, the list will only provide `active` values if unspecified, but the filter can be controlled with the `--state` flag.

For more general information about attribute values, see the `values` subcommand.

## Example

```shell
otdfctl policy attributes values list --attribute-id 3c51a593-cbf8-419d-b7dc-b656d0bedfbb
```

```shell
  SUCCESS   Found values list                                                                                                                                                                                                                                                                                                                                                                                                               
╭───────────────────────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────┬───────────────────────────────────────────┬──────────────┬──────────────┬──────────────╮
│ID                                                                     │Fqn                                                      │Active                                     │Labels        │Created At    │Updated At    │
├───────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────┼───────────────────────────────────────────┼──────────────┼──────────────┼──────────────┤
│355743c1-c0ef-4e8d-9790-d49d883dbc7d                                   │https://opentdf.io/attr/myattribute/value/myvalue1       │true                                       │[]            │Tue Dec 17 19…│Tue Dec 17 19…│
│b20458b0-1855-4608-8869-3f6199bc2878                                   │https://opentdf.io/attr/myattribute/value/myvalue2       │true                                       │[]            │Tue Dec 17 19…│Tue Dec 17 19…│
╰───────────────────────────────────────────────────────────────────────┴─────────────────────────────────────────────────────────┴───────────────────────────────────────────┴──────────────┴──────────────┴──────────────╯
  NOTE   Use 'otdfctl policy attributes values get --id=<id> --json' to see all properties
```
