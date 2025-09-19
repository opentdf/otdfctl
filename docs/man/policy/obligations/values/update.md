---
title: Update an obligation value
command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      shorthand: i
      description: ID of the obligation value to update
      required: true
    - name: value
      shorthand: v
      description: Optional updated value of the obligation value (must be unique within the definition)
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

Update the `value` and/or metadata labels for an obligation value.

If PEPs rely on this value, a value update could break access.

Make sure you know what you are doing.

For more information about obligation values, see the manual for the `values` subcommand.

## Example

```shell
otdfctl policy obligations values update --id 3c51a593-cbf8-419d-b7dc-b656d0bedfbb --value new_value
```
