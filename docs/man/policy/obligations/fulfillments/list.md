---
title: List obligation values
command:
  name: list
  flags:
    - name: obligation-id
      shorthand: o
      description: The ID of the obligation to list values for
    - name: state
      shorthand: s
      description: Filter by state
      enum:
        - active
        - inactive
        - any
      default: active
---

Retrieves all obligation valuess stored in platform policy.
