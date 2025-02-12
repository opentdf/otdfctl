---
title: Update Public Key Metadata
command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      shorthand: i
      description: ID of the Public Key
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---

Update the metadata of a public key. The public key information itself cannot be updated. To update a public key create a new key with the updated information.

## Example 

```shell
otdfctl policy kas-registry public-key update --id=62857b55-560c-4b67-96e3-33e4670ecb3b --label key=value
```

