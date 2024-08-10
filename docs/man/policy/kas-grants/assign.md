---
title: Assign a grant

command:
  name: update
  aliases:
    - u
    - create
    - add
    - new
    - upsert
  description: Assign a grant of a KAS to an Attribute Definition or Value
  flags:
    - name: attribute-id
      shorthand: a
      description: The ID of the attribute definition being assigned a KAS Grant
      required: true
    - name: value-id
      shorthand: v
      description: The ID of the attribute value being assigned a KAS Grant
      required: true
    - name: kas-id
      shorthand: k
      description: The ID of the Key Access Server being assigned to the grant
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---
