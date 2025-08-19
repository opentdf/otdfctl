---
title: Delete an obligation definition
command:
  name: delete
  flags:
    - name: id
      shorthand: i
      description: ID of the obligation
      required: false
    - name: fqn
      shorthand: f
      description: FQN of the obligation
      required: false
    - name: force
      description: Force deletion without interactive confirmation
---

Removes an obligation definition from platform Policy.

Obligation deletion cascades to the associated Obligation Values and Action Attribute Values.

For more information about Obligations, see the manual for the `obligations` subcommand.

## Example 

```shell
otdfctl policy obligations delete --id 217b300a-47f9-4bee-be8c-d38c880053f7
```
