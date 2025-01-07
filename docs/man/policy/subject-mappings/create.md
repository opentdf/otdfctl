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
