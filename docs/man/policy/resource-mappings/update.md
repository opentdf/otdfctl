---
title: Update a resource mapping
command:
  name: update
  flags:
    - name: id
      description: The ID of the resource mapping to update.
      default: ""
    - name: attribute-value-id  
      description: The ID of the attribute value to map to the resource.
      default: ""
    - name: terms
      description: The synonym terms to match for the resource mapping.
      type: string-slice
      default: ""
---
