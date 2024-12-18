---
title: Update attribute value

command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      shorthand: i
      description: The ID of the attribute value to update
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---

Attribute Value changes can be dangerous, so this command is for updates considered "safe" (currently just mutations to metadata `labels`).

For unsafe updates, see the dedicated `unsafe update` command. For more general information, see the `values` subcommand.

For more general information about attributes, see the `attributes` subcommand.

## Example

```shell
otdfctl policy attributes values update --id 355743c1-c0ef-4e8d-9790-d49d883dbc7d --label hello=world
```

```shell
  SUCCESS   Updated values: 355743c1-c0ef-4e8d-9790-d49d883dbc7d                                                                                                                                                                                                                                                                                                                                                                              
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │355743c1-c0ef-4e8d-9790-d49d883dbc7d                                                                                                            │
│FQN                                                                      │https://opentdf.io/attr/myattribute/value/myvalue1                                                                                              │
│Value                                                                    │myvalue1                                                                                                                                        │
│Created At                                                               │Tue Dec 17 19:06:55 UTC 2024                                                                                                                    │
│Updated At                                                               │Tue Dec 17 19:11:50 UTC 2024                                                                                                                    │
│Labels                                                                   │[hello: world]                                                                                                                                  │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy attributes values get --id=355743c1-c0ef-4e8d-9790-d49d883dbc7d --json' to see all properties
```
