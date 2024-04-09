---
title: Update an attribute
command:
  name: update
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      type: string-slice
      default: ""
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      type: bool
      default: false
---
