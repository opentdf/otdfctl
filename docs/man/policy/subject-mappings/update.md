---
title: Update a subject mapping 
command:
  name: update
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
      default: ""
    - name: action-custom
      description: The custom action to map to a subject set
      shorthand: c
      required: false
      default: ""
    - name: subject-condition-set-id
      description: Known pre-existing Subject Condition Set Id
      required: true
      default: ""
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ""
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---

'Actions' are updated in place, destructively replacing the current set. If you want to add or remove actions, you must provide the full set of actions on update.
