---
title: Update a resource mapping
command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      description: The ID of the resource mapping to update.
      default: ''
    - name: attribute-value-id
      description: The ID of the attribute value to map to the resource.
      default: ''
    - name: terms
      description: The synonym terms to match for the resource mapping.
      default: ''
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---

Alter the attribute value associated with a resource mapping's terms, or fully replace the terms in a given resource mapping.

For more information about resource mappings, see the `resource-mappings` subcommand.

## Examples

```shell
otdfctl policy resource-mappings update --id=3ff446fb-8fb1-4c04-8023-47592c90370c --terms newterm1,newterm2
```

```shell
  SUCCESS   Updated resource-mappings: 3ff446fb-8fb1-4c04-8023-47592c90370c                                                                                                                                                                                                                                                                                                                                                                       
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │3ff446fb-8fb1-4c04-8023-47592c90370c                                                                                                            │
│Attribute Value Id                                                       │891cfe85-b381-4f85-9699-5f7dbfe2a9ab                                                                                                            │
│Attribute Value                                                          │myvalue1                                                                                                                                        │
│Terms                                                                    │newterm1, newterm2                                                                                                                              │
│Created At                                                               │Wed Dec 18 05:53:53 UTC 2024                                                                                                                    │
│Updated At                                                               │Wed Dec 18 05:58:11 UTC 2024                                                                                                                    │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy resource-mappings get --id=3ff446fb-8fb1-4c04-8023-47592c90370c --json' to see all properties
```
