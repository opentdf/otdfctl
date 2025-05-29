---
title: Set Base Key
command:
  name: set
  aliases:
    - s
  flags:
    - name: id
      shorthand: i
      description: ID of the key to retrieve
    - name: keyId
      shorthand: k
      description: KeyID of the key to retrieve
    - name: kasUri
      shorthand: u
      description: URI of the Key access server that the key is assigned to.
    - name: kasName
      shorthand: n
      description: Name of the Key access server that the key is assigned to.
    - name: kasId
      shorthand: d
      description: Id of the Key access server that the key is assigned to.
    
---

Command for setting a base key to be used for encryption operations on data where no attributes are present or where no keys are present on found attributes.

## Examples

```
otdfctl policy kas-registry key base set --id 8af2059f-5d0b-46c2-84f0-bed8a6101d90

otdfctl policy kas-registry key base set --kasUri
```
