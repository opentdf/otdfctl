---
title: Get a subject mapping
command:
  name: get
  aliases:
    - g
  flags:
    - name: id
      description: The ID of the subject mapping to get
      shorthand: i
      required: true
      default: ''
---

Retrieve the specifics of a Subject Mapping.

For more information about subject mappings, see the `subject-mappings` subcommand.

```shell
otdfctl policy subject-mappings get --id 39866dd2-368b-41f6-b292-b4b68c01888b
```

```shell
  SUCCESS   Found subject-mappings: 39866dd2-368b-41f6-b292-b4b68c01888b                                                                                                                                                                                                                                                                                                                                                                             
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │39866dd2-368b-41f6-b292-b4b68c01888b                                                                                                            │
│Attribute Value: Id                                                      │891cfe85-b381-4f85-9699-5f7dbfe2a9ab                                                                                                            │
│Attribute Value: Value                                                   │myvalue1                                                                                                                                        │
│Actions                                                                  │[{"Value":{"Standard":1}}]                                                                                                                      │
│Subject Condition Set: Id                                                │8dc98f65-5f0a-4444-bfd1-6a818dc7b447                                                                                                            │
│Subject Condition Set                                                    │[{"condition_groups":[{"conditions":[{"subject_external_selector_value":".example.field.one","operator":1,"subject_external_values":["myvalue",…│
│Created At                                                               │Wed Dec 18 15:40:50 UTC 2024                                                                                                                    │
│Updated At                                                               │Wed Dec 18 15:40:50 UTC 2024                                                                                                                    │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy subject-mappings get --id=39866dd2-368b-41f6-b292-b4b68c01888b --json' to see all properties 
```
