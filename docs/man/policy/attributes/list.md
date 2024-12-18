---
title: List attribute definitions
command:
  name: list
  aliases:
    - l
  flags:
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

By default, the list will only provide `active` attributes if unspecified, but the filter can be controlled with the `--state` flag.

For more general information about attributes, see the `attributes` subcommand.

## Example

```shell
otdfctl policy attributes list
```

```shell
  SUCCESS   Found attributes list                                                                                                                                                                                                                                                                                                                                                                                                                 
╭──────────────────────────────────────────────────┬────────────────────────────────────────┬──────────────────────────────┬────────────────────┬────────────────────┬────────────────────┬──────────┬──────────┬──────────╮
│ID                                                │Namespace                               │Name                          │Rule                │Values              │Active              │Labels    │Created At│Updated At│
├──────────────────────────────────────────────────┼────────────────────────────────────────┼──────────────────────────────┼────────────────────┼────────────────────┼────────────────────┼──────────┼──────────┼──────────┤
│3c51a593-cbf8-419d-b7dc-b656d0bedfbb              │opentdf.io                              │myattribute                   │ANY_OF              │[]                  │true                │[]        │Tue Dec 1…│Tue Dec 1…│
│6a261d68-0899-4e17-bb2f-124abba7c09c              │example.com                             │attr1                         │ANY_OF              │[value1, value2]    │true                │[]        │Fri Nov  …│Fri Nov  …│
│e1536f25-d287-43ed-9ad9-2cf4a7698e5f              │example.com                             │attr2                         │ALL_OF              │[value2, value1]    │true                │[]        │Fri Nov  …│Fri Nov  …│
╰──────────────────────────────────────────────────┴────────────────────────────────────────┴──────────────────────────────┴────────────────────┴────────────────────┴────────────────────┴──────────┴──────────┴──────────╯
  NOTE   Use 'otdfctl policy attributes get --id=<id> --json' to see all properties
```
