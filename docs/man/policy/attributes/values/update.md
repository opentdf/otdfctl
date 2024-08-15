---
title: Update attribute value

command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      shorthand: i
      description: The ID of the attribute value to update
    - name: value
      shorthand: v
      description: The new value
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ""
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---

This command allows you to manage the values of an attribute.
