---
title: List attributes
command:
  name: list
  aliases:
    - l
  flags:
    - name: state
      shorthand: s
      description: Filter by state
      enum:
        - active
        - inactive
        - any
      default: active
---
