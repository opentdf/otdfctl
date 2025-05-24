---
title: Get Key
command:
  name: get
  aliases:
    - g
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

This command retrieves details of a specific key by sending a `GetKeyRequest` to the platform.
