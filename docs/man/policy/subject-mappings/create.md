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

Create a Subject Mapping to entitle an entity (via existing or new Subject Condition Set) to an Attribute Value.

For more information about subject mappings, see the `subject-mappings` subcommand.

For more information about subject condition sets, see the `subject-condition-sets` subcommand.
