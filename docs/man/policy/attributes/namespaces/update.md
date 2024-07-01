---
title: Update a attribute namespace
command:
  name: update
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute namespace
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ""
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---

# Update an Attribute Namespace

Attribute Namespace changes can be dangerous, so this command is for updates considered "safe."

For unsafe updates, see the dedicated `update` command.