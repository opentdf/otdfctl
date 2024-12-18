---
title: Delete a Subject Condition Set

command:
  name: delete
  flags:
    - name: id
      description: The ID of the subject condition set to delete
      shorthand: i
      required: true
    - name: force
      description: Force deletion without interactive confirmation (dangerous)
---

For more information about subject condition sets, see the `subject-condition-sets` subcommand.

## Example

```shell
otdfctl policy subject-condition-sets delete --id=bfade235-509a-4a6f-886a-812005c01db5
```

```shell
  SUCCESS   Deleted subject-condition-sets: bfade235-509a-4a6f-886a-812005c01db5                                                                                                                                                                                                                                                                                                                                                                     
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │bfade235-509a-4a6f-886a-812005c01db5                                                                                                            │
│SubjectSets                                                              │[{"conditionGroups":[{"conditions":[{"subjectExternalSelectorValue":".example.field.one","operator":"SUBJECT_MAPPING_OPERATOR_ENUM_IN","subject…│
│Created At                                                               │Wed Dec 18 06:44:39 UTC 2024                                                                                                                    │
│Updated At                                                               │Wed Dec 18 06:54:28 UTC 2024                                                                                                                    │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy subject-condition-sets list --json' to see all properties  
```
