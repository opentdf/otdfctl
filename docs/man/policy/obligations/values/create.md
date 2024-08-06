---
title: Create an obligation
command:
  name: create
  flags:
    - name: value
      shorthand: v
      description: Value being added to the existing obligation (i.e. 'watermark')
      required: true
    - name: obligation
      shorthand: o
      description: ID of the parent obligation
      required: true
    - name: attr-val
      shorthand: a
      description: ID of assigned attribute value(s) for derived obligations
      required: false
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

For more information about the significance of obligations and how they are utilized for derived obligations on attribute values
or when added directly to a TDF, see the parent command above.
