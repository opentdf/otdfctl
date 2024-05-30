---
title: Create an attribute
command:
  name: create
  flags:
    - name: name
      shorthand: n
      description: Name of the attribute
      required: true
    - name: rule
      shorthand: r
      description: Rule of the attribute
      enum:
        - ANY_OF
        - ALL_OF
        - HIERARCHY
      required: true
    - name: value
      shorthand: v
      description: Value of the attribute
      required: true
    - name: namespace
      shorthand: s
      description: Namespace of the attribute
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ""
---
