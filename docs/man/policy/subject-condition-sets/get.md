---
title: Get a Subject Condition Set

command:
  name: get
  aliases:
    - g
  flags:
    - name: id
      description: The ID of the subject condition set to get
      shorthand: i
      required: true
---

For more information about subject condition sets, see the `subject-condition-sets` subcommand.

## Example

```shell
otdfctl policy subject-condition-sets get --id=bfade235-509a-4a6f-886a-812005c01db5
```

```shell
  SUCCESS   Found subject-condition-sets: bfade235-509a-4a6f-886a-812005c01db5                                                                                                                                                                                                                                                                                                                                                                      
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │bfade235-509a-4a6f-886a-812005c01db5                                                                                                            │
│SubjectSets                                                              │[{"conditionGroups":[{"conditions":[{"subjectExternalSelectorValue":".example.field.one","operator":"SUBJECT_MAPPING_OPERATOR_ENUM_IN","subject…│
│Created At                                                               │Wed Dec 18 06:44:39 UTC 2024                                                                                                                    │
│Updated At                                                               │Wed Dec 18 06:44:39 UTC 2024                                                                                                                    │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy subject-condition-sets get --id=bfade235-509a-4a6f-886a-812005c01db5 --json' to see all properties
```
