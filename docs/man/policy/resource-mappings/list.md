---
title: List resource mappings
command:
  name: list
  aliases:
    - l
  flags:
    - name: limit
      shorthand: l
      description: Limit retrieved count
    - name: offset
      shorthand: o
      description: Offset (page) quantity from start of the list
---

For more information about resource mappings, see the `resource-mappings` subcommand.

## Examples

```shell
otdfctl policy resource-mappings get --id=3ff446fb-8fb1-4c04-8023-47592c90370c
```

```shell
  SUCCESS   Found resource-mappings list                                                                                                                                                                                    
                                                                                                                                                                                                                            
╭────────────────────────────────────────────────────────┬─────────────────────────────────────────────┬─────────────────────────────────────────────┬─────────────────────────────────┬───────────┬───────────┬───────────╮
│ID                                                      │Attribute Value Id                           │Attribute Value                              │Terms                            │Labels     │Created At │Updated At │
├────────────────────────────────────────────────────────┼─────────────────────────────────────────────┼─────────────────────────────────────────────┼─────────────────────────────────┼───────────┼───────────┼───────────┤
│3ff446fb-8fb1-4c04-8023-47592c90370c                    │891cfe85-b381-4f85-9699-5f7dbfe2a9ab         │myvalue1                                     │term1, term2                     │[]         │Wed Dec 18…│Wed Dec 18…│
│02092d67-fffa-4030-9775-b5cd5d581e1f                    │74babca6-016f-4f3e-a99b-4e46ea8d0fd8         │myvalue2                                     │term2, term4                     │[]         │Fri Nov  1…│Fri Nov  1…│
╰────────────────────────────────────────────────────────┴─────────────────────────────────────────────┴─────────────────────────────────────────────┴─────────────────────────────────┴───────────┴───────────┴───────────╯
  NOTE   Use 'otdfctl policy resource-mappings get --id=<id> --json' to see all properties
```
