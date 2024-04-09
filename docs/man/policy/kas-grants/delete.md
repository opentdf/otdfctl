---
title: Delete a grant

command:
  name: delete
  description: Delete a grant
  flags:
    - name: attribute-id
      shorthand: a
      description: The attribute to delete
      required: true
      type: string
    - name: value-id
      shorthand: v
      description: The value of the attribute to delete
      required: true
      type: string
    - name: kas-id
      shorthand: k
      description: The Key Access Server ID
      required: true
      type: string
---