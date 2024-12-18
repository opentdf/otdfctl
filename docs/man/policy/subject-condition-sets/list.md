---
title: List Subject Condition Set

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

For more information about subject condition sets, see the `subject-condition-sets` subcommand.

## Example

```shell
otdfctl policy subject-condition-set list
```

```shell
  SUCCESS   Found subject-condition-sets list                                                                                                                                                                               
                                                                                                                                                                                                                            
╭──────────────────────────────────────────────────────────────────────────────────────┬─────────────────────────────────────────────────────────────────────┬──────────────────┬──────────────────┬───────────────────────╮
│ID                                                                                    │SubjectSets                                                          │Labels            │Created At        │Updated At             │
├──────────────────────────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────────────────────┼──────────────────┼──────────────────┼───────────────────────┤
│8b80eb7c-cecb-44d4-91a7-f14ada74d4ce                                                  │[{"conditionGroups":[{"conditions":[{"subjectExternalSelectorValue":…│[]                │Mon Dec 16 16:00:…│Mon Dec 16 16:00:33 UT…│
│bfade235-509a-4a6f-886a-812005c01db5                                                  │[{"conditionGroups":[{"conditions":[{"subjectExternalSelectorValue":…│[]                │Wed Dec 18 06:44:…│Wed Dec 18 06:44:39 UT…│
╰──────────────────────────────────────────────────────────────────────────────────────┴─────────────────────────────────────────────────────────────────────┴──────────────────┴──────────────────┴───────────────────────╯
  NOTE   Use 'otdfctl policy subject-condition-sets get --id=<id> --json' to see all properties 
```
