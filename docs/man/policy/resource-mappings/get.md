---
title: Get a resource mapping
command:
  name: get
  aliases:
    - g
  flags:
    - name: id
      description: The ID of the resource mapping to get.
      default: ''
---

For more information about resource mappings, see the `resource-mappings` subcommand.

## Examples

```shell
otdfctl policy resource-mappings get --id=3ff446fb-8fb1-4c04-8023-47592c90370c
```

```shell
  SUCCESS   Found resource-mappings: 3ff446fb-8fb1-4c04-8023-47592c90370c                                                                                                                                                                                                                                                                                                                                                      
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │3ff446fb-8fb1-4c04-8023-47592c90370c                                                                                                            │
│Attribute Value Id                                                       │891cfe85-b381-4f85-9699-5f7dbfe2a9ab                                                                                                            │
│Attribute Value                                                          │myvalue1                                                                                                                                        │
│Terms                                                                    │term1, term2                                                                                                                                    │
│Created At                                                               │Wed Dec 18 05:53:53 UTC 2024                                                                                                                    │
│Updated At                                                               │Wed Dec 18 05:53:53 UTC 2024                                                                                                                    │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy resource-mappings get --id=3ff446fb-8fb1-4c04-8023-47592c90370c --json' to see all properties  
```
