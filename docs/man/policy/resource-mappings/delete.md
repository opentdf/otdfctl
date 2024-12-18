---
title: Delete a resource mapping
command:
  name: delete
  flags:
    - name: id
      description: The ID of the resource mapping to delete
      default: ''
    - name: force
      description: Force deletion without interactive confirmation (dangerous)
---

For more information about resource mappings, see the `resource-mappings` subcommand.

## Examples

```shell
otdfctl policy resource-mappings delete --id=3ff446fb-8fb1-4c04-8023-47592c90370c
```

```shell
  SUCCESS   Deleted resource-mappings: 3ff446fb-8fb1-4c04-8023-47592c90370c                                                                                                                                                                                                                                                                                                                                                                          
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │3ff446fb-8fb1-4c04-8023-47592c90370c                                                                                                            │
│Attribute Value Id                                                       │                                                                                                                                                │
│Attribute Value                                                          │                                                                                                                                                │
│Terms                                                                    │                                                                                                                                                │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy resource-mappings list --json' to see all properties
```
