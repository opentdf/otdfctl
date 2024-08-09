---
title: Update a grant

command:
  name: update
  aliases:
    - u
    - create
    - add
    - new
    - upsert
  description: Update a grant
  flags:
    - name: attribute-id
      shorthand: a
      description: The attribute to delete
      required: true
    - name: value-id
      shorthand: v
      description: The value of the attribute to delete
      required: true
    - name: kas-id
      shorthand: k
      description: The Key Access Server ID
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ""
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---