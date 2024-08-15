---
title: Update a resource mapping
command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      description: The ID of the resource mapping to update.
      default: ""
    - name: attribute-value-id  
      description: The ID of the attribute value to map to the resource.
      default: ""
    - name: terms
      description: The synonym terms to match for the resource mapping.
      default: ""
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ""
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---
