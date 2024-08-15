---
title: List attribute values
command:
  name: list
  aliases:
    - ls
    - l
  flags:
    - name: attribute-id
      shorthand: a
      description: The ID of the attribute to list values for
    - name: state
      shorthand: s
      description: Filter by state
      enum:
        - active
        - inactive
        - any
      default: active
---

# List attribute values
