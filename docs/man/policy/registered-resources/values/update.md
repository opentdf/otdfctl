---
title: Update a Registered Resource Value
command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      shorthand: i
      description: ID of the registered resource value to update
    - name: value
      shorthand: v
      description: Optional updated value of the registered resource value (must be unique within the Registered Resource)
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

Update the `value` and/or metadata labels for a Registered Resource Value.

If PEPs rely on this value, a value update could break access.

Make sure you know what you are doing.

For more information about Registered Resource Values, see the manual for the `values` subcommand.

## Example

```shell
otdfctl policy registered-resources values update --id 3c51a593-cbf8-419d-b7dc-b656d0bedfbb --value new_value
```
