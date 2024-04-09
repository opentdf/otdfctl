---
title: Create an attribute namespace
command:
  name: create
  flags:
    - name: name
      shorthand: n
      description: Name of the attribute namespace
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      type: string-slice
      default: ""
---
