---
title: Create an obligation
command:
  name: create
  flags:
    - name: name
      shorthand: n
      description: Name of the obligation (i.e. 'drm' for Digital Rights Management)
      required: true
    - name: value
      shorthand: v
      description: Values of the obligation (i.e. 'watermark')
      required: false
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

An obligation, like an attribute definition, is a parent that can contain one or more values.

For more information about the significance of obligations and how they are utilized for derived obligations on attribute values
or when added directly to a TDF, view the parent command documentation with `--help`.