---
title: Create a resource mapping
command:
  name: create
  flags:
    - name: attribute-value-id
      description: The ID of the attribute value to map to the resource.
      default: ""
    - name: terms
      description: The synonym terms to match for the resource mapping.
      type: string-slice
      default: ""
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      type: string-slice
      default: ""
---
