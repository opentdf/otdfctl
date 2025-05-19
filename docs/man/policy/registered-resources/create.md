---
title: Create a Registered Resource
command:
  name: create
  aliases:
    - c
    - add
    - new
  flags:
    - name: name
      shorthand: n
      description: Name of the registered resource (must be unique within Policy)
      required: true
    - name: value
      shorthand: v
      description: Value of the registered resource (i.e. 'value1', must be unique within the Registered Resource)
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

## Examples

Create a registered resource named 'my_resource' with value 'my_value':

```shell
otdfctl policy registered-resources create --name my_resource --v my_value
```
