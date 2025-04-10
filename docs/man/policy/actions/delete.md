---
title: Delete a Custom Action
command:
  name: delete
  flags:
    - name: id
      shorthand: i
      description: ID of the custom action
      required: true
---

Removes a Custom Action from platform Policy. Standard Actions cannot be deleted.

Action deletion cascades to entitlement Subject Mappings, Obligations, and Non-Data
Resource entitlement requirements.

Make sure you know what you are doing.

For more information about Actions, see the manual for the `actions` subcommand.

## Example 

```shell
otdfctl policy actions delete --id 217b300a-47f9-4bee-be8c-d38c880053f7
```
