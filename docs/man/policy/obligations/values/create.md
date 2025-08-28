---
title: Create an obligation value
command:
  name: create
  aliases:
    - c
    - add
    - new
  flags:
    - name: obligation
      shorthand: o
      description: Identifier of the associated obligation (ID or FQN)
      required: true
    - name: value
      shorthand: v
      description: Value of the obligation (i.e. 'value1', must be unique within the definition)
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

Add a value to an obligation in the platform Policy.

For more information, see the `obligations` subcommand.

## Examples

Create an obligation value for the obligation with ID '3c51a593-cbf8-419d-b7dc-b656d0bedfbb', and value 'my_value':

```shell
otdfctl policy obligations values create --obligation 3c51a593-cbf8-419d-b7dc-b656d0bedfbb --value my_value
```
