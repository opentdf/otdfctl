---
title: List attribute namespaces
command:
  name: list
  aliases:
    - ls
    - l
  flags:
    - name: state
      shorthand: s
      description: Filter by state [active, inactive, any]
    - name: limit
      shorthand: l
      description: Limit retrieved count
    - name: offset
      shorthand: o
      description: Offset (page) quantity from start of the list
---

For more general information, see the `namespaces` subcommand.

## Example

```shell
otdfctl policy attributes namespaces list
```

```shell
SUCCESS   Found namespaces list                                                                                                                                                                                                                                                                                                                                                                                                                
╭───────────────────────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────┬───────────────────────────────────────────┬──────────────┬──────────────┬──────────────╮
│ID                                                                     │Name                                                     │Active                                     │Labels        │Created At    │Updated At    │
├───────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────┼───────────────────────────────────────────┼──────────────┼──────────────┼──────────────┤
│87ba60e1-da12-4889-95fd-267968bf0896                                   │scenario.com                                             │true                                       │[]            │Fri Nov  1 14…│Fri Nov  1 14…│
│8f1d8839-2851-4bf4-8bf4-5243dbfe517d                                   │example.com                                              │true                                       │[]            │Fri Nov  1 14…│Fri Nov  1 14…│
│d69cf14d-744b-48cf-aab4-43756e97a8e5                                   │example.net                                              │true                                       │[]            │Fri Nov  1 14…│Fri Nov  1 14…│
│0d94e00a-7bd3-4482-afe3-f1e4b03c1353                                   │example.org                                              │true                                       │[]            │Fri Nov  1 14…│Fri Nov  1 14…│
│e3802200-7d16-45c4-be55-3f1a2e90adb1                                   │opentdf.io                                               │true                                       │[]            │Tue Dec 17 16…│Tue Dec 17 16…│
╰───────────────────────────────────────────────────────────────────────┴─────────────────────────────────────────────────────────┴───────────────────────────────────────────┴──────────────┴──────────────┴──────────────╯
  NOTE   Use 'otdfctl policy attributes namespaces get --id=<id> --json' to see all properties 
```
