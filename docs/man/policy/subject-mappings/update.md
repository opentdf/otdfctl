---
title: Update a subject mapping
command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      description: The ID of the subject mapping to update
      shorthand: i
      required: true
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
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---

Update a Subject Mapping to alter entitlement of an entity to an Attribute Value.

`Actions` are updated in place, destructively replacing the current set. If you want to add or remove actions, you must provide the full set of actions on update.

At this time, creation of a new SCS during update of a subject mapping is not supported.

For more information about subject mappings, see the `subject-mappings` subcommand.

For more information about subject condition sets, see the `subject-condition-sets` subcommand.
