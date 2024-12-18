---
title: Delete an attribute definition
command:
  name: delete
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute definition
      required: true
---

# Unsafe Delete Warning

Deleting an Attribute Definition cascades deletion of any Attribute Values and any associated mappings underneath.

Any existing TDFs containing the deleted attribute of this name will be rendered inaccessible until it has been recreated.

Make sure you know what you are doing.

For more general information about attributes, see the `attributes` subcommand.

## Example

```shell
otdfctl policy attributes unsafe delete --id 3c51a593-cbf8-419d-b7dc-b656d0bedfbb
```

```shell
  SUCCESS   Deleted unsafe: 3c51a593-cbf8-419d-b7dc-b656d0bedfbb                                                                                                                                                                                                                                                                                                                                                        
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Deleted                                                                  │true                                                                                                                                            │
│Id                                                                       │3c51a593-cbf8-419d-b7dc-b656d0bedfbb                                                                                                            │
│Name                                                                     │myattribute                                                                                                                                     │
│Created At                                                               │Tue Dec 17 18:33:06 UTC 2024                                                                                                                    │
│Updated At                                                               │Tue Dec 17 19:46:00 UTC 2024                                                                                                                    │
│Labels                                                                   │[hello: world]                                                                                                                                  │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy attributes list --json' to see all properties
```
