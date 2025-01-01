---
title: Get an attribute definition
command:
  name: get
  aliases:
    - g
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute
---

Retrieve an attribute along with its metadata, rule, and values.

For more general information about attributes, see the `attributes` subcommand.

## Example

```shell
otdfctl policy attributes get --id=3c51a593-cbf8-419d-b7dc-b656d0bedfbb
```

```shell
  SUCCESS   Found attributes: 3c51a593-cbf8-419d-b7dc-b656d0bedfbb                                                                                                                                                          
                                                                                                                                                                                                                            
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │3c51a593-cbf8-419d-b7dc-b656d0bedfbb                                                                                                            │
│Name                                                                     │myattribute                                                                                                                                     │
│Rule                                                                     │ANY_OF                                                                                                                                          │
│Values                                                                   │[]                                                                                                                                              │
│Namespace                                                                │opentdf.io                                                                                                                                      │
│Created At                                                               │Tue Dec 17 18:33:06 UTC 2024                                                                                                                    │
│Updated At                                                               │Tue Dec 17 18:33:06 UTC 2024                                                                                                                    │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy attributes get --id=3c51a593-cbf8-419d-b7dc-b656d0bedfbb --json' to see all properties 
```
