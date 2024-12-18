---
title: Update an attribute definition
command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ""
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---

Attribute Definition changes can be dangerous, so this command is for updates considered "safe" (currently just mutations to metadata `labels`).

For unsafe updates, see the dedicated `unsafe update` command. For more general information, see the `attributes` subcommand.

For more general information about attributes, see the `attributes` subcommand.

## Example

```shell
otdfctl policy attributes update --id=3c51a593-cbf8-419d-b7dc-b656d0bedfbb --label hello=world
```

```shell
  SUCCESS   Updated attributes: 3c51a593-cbf8-419d-b7dc-b656d0bedfbb                                                                                                                                                                                                                                                                                                                                                                           
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │3c51a593-cbf8-419d-b7dc-b656d0bedfbb                                                                                                            │
│Name                                                                     │myattribute                                                                                                                                     │
│Created At                                                               │Tue Dec 17 18:33:06 UTC 2024                                                                                                                    │
│Updated At                                                               │Tue Dec 17 18:39:26 UTC 2024                                                                                                                    │
│Labels                                                                   │[hello: world]                                                                                                                                  │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy attributes get --id=3c51a593-cbf8-419d-b7dc-b656d0bedfbb --json' to see all properties 
```
