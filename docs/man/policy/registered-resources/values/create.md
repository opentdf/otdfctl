---
title: Create Registered Resource Value
command:
  name: create
  aliases:
    - c
    - add
    - new
  flags:
    - name: resource-id
      shorthand: i
      description: ID of the associated registered resource
      required: true
    - name: value
      shorthand: v
      description: Value of the registered resource (i.e. 'value1', must be unique within the Registered Resource)
    - name: action-attribute-value
      shorthand: a
      description: "Optional action attribute values in the format: \"<action_id>|<action_name>;<attribute_value_id|attribute_value_fqn>\""
      default: ''
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

Add a value to a registered resource in the platform Policy.

A registered resource value `value` is normalized to lower case and may contain hyphens or dashes between other alphanumeric characters.

For more information, see the `registered-resources` subcommand.

## Examples

Create a registered resource value for the registered resource with ID '3c51a593-cbf8-419d-b7dc-b656d0bedfbb' with value 'my_value':

```shell
otdfctl policy registered-resources values create --resource-id 3c51a593-cbf8-419d-b7dc-b656d0bedfbb --value my_value
```
