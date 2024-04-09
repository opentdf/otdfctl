---
title: Create an attribute value
command:
  name: create
  flags:
    - name: attribute-id
      shorthand: a
      description: The ID of the attribute to create a value for
    - name: value
      shorthand: v
      description: The value to create
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ""
---

This command allows you to manage the values of an attribute.
