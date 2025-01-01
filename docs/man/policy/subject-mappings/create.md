---
title: Create a new subject mapping
command:
  name: create
  aliases:
    - new
    - add
    - c
  flags:
    - name: attribute-value-id
      description: The ID of the attribute value to map to a subject set
      shorthand: a
      required: true
      default: ''
    - name: action-standard
      description: The standard action to map to a subject set
      enum:
        - DECRYPT
        - TRANSMIT
      shorthand: s
      required: true
      default: ''
    - name: action-custom
      description: The custom action to map to a subject set
      shorthand: c
      required: false
      default: ''
    - name: subject-condition-set-id
      description: Known preexisting Subject Condition Set Id
      required: true
      default: ''
    - name: subject-condition-set-new
      description: JSON array of Subject Sets to create a new Subject Condition Set associated with the created Subject Mapping
      required: false
      default: ''
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

The possible values for standard actions are DECRYPT and TRANSMIT.

Create a Subject Mapping to entitle an entity (via existing or new Subject Condition Set) to an Attribute Value.

For more information about subject mappings, see the `subject-mappings` subcommand.

For more information about subject condition sets, see the `subject-condition-sets` subcommand.

## Examples

Create a subject mapping linking to an existing subject condition set:
```shell
otdfctl policy subject-mapping create --attribute-value-id 891cfe85-b381-4f85-9699-5f7dbfe2a9ab --action-standard DECRYPT --subject-condition-set-id 8dc98f65-5f0a-4444-bfd1-6a818dc7b447
```

```shell
  SUCCESS   Created subject-mappings: 39866dd2-368b-41f6-b292-b4b68c01888b                                                                                                                                                                                                                                                                                                                                                                          
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │39866dd2-368b-41f6-b292-b4b68c01888b                                                                                                            │
│Attribute Value Id                                                       │891cfe85-b381-4f85-9699-5f7dbfe2a9ab                                                                                                            │
│Actions                                                                  │[{"Value":{"Standard":1}}]                                                                                                                      │
│Subject Condition Set: Id                                                │8dc98f65-5f0a-4444-bfd1-6a818dc7b447                                                                                                            │
│Subject Condition Set                                                    │[{"condition_groups":[{"conditions":[{"subject_external_selector_value":".example.field.one","operator":1,"subject_external_values":["myvalue",…│
│Created At                                                               │Wed Dec 18 15:40:50 UTC 2024                                                                                                                    │
│Updated At                                                               │Wed Dec 18 15:40:50 UTC 2024                                                                                                                    │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy subject-mappings get --id=39866dd2-368b-41f6-b292-b4b68c01888b --json' to see all properties  
```

Or you can create a mapping that linked to a new subject condition set:
```shell
otdfctl policy subject-mapping create --attribute-value-id 891cfe85-b381-4f85-9699-5f7dbfe2a9ab --action-standard DECRYPT --subject-condition-set-new '[                                           
  {
    "condition_groups": [
      {
        "conditions": [
          {
            "operator": 1,
            "subject_external_values": ["myvalue", "myothervalue"],
            "subject_external_selector_value": ".example.field.one"
          },
          {
            "operator": 2,
            "subject_external_values": ["notpresentvalue"],
            "subject_external_selector_value": ".example.field.two"
          }
        ],
        "boolean_operator": 2
      }
    ]
  }
]'
```

```shell
  SUCCESS   Created subject-mappings: d71c4028-ce64-453b-8aa7-6edb45fbb848                                                                                                                                                                                                                                                                                                                                                                           
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │d71c4028-ce64-453b-8aa7-6edb45fbb848                                                                                                            │
│Attribute Value Id                                                       │891cfe85-b381-4f85-9699-5f7dbfe2a9ab                                                                                                            │
│Actions                                                                  │[{"Value":{"Standard":1}}]                                                                                                                      │
│Subject Condition Set: Id                                                │738736ee-880d-40da-acae-672d1deff00f                                                                                                            │
│Subject Condition Set                                                    │[{"condition_groups":[{"conditions":[{"subject_external_selector_value":".example.field.one","operator":1,"subject_external_values":["myvalue",…│
│Created At                                                               │Wed Dec 18 15:41:55 UTC 2024                                                                                                                    │
│Updated At                                                               │Wed Dec 18 15:41:55 UTC 2024                                                                                                                    │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy subject-mappings get --id=d71c4028-ce64-453b-8aa7-6edb45fbb848 --json' to see all properties
```
