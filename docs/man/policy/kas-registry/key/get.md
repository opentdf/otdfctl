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
    
---

This command retrieves details of a specific key by sending a `GetKeyRequest` to the platform.
