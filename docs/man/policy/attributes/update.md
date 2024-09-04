---
title: Update an attribute
command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ""
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---

# Update an attribute

Attribute Definition changes can be dangerous, so this command is for updates considered "safe" (currently just mutations to metadata `labels`).

For unsafe updates, see the dedicated `unsafe update` command. For more general information, see the `attributes` subcommand.

For more general information about attributes, see the `attributes` subcommand.
