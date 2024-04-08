---
title: Create a new subject mapping 
command:
  name: create
  flags:
    - name: attribute-value-id
      description: The ID of the attribute value to map to a subject set
      shorthand: a
      required: true
      default: ""
    - name: action-standard
      description: The standard action to map to a subject set
      shorthand: s
      type: enum
      values: ["DECRYPT", "TRANSMIT"]
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
    - name: subject-condition-set-new
      description: JSON array of Subject Sets to create a new Subject Condition Set associated with the created Subject Mapping
      required: false
      default: ""
---