---
title: Authenticate to the platform with the client-credentials flow

command:
  name: client-credentials
  flags:
    - name: client-id
      description: Client ID
      required: true
    - name: client-secret
      description: Client secret
---

Allows the user to login in via client ID and secret. This will subsequently be stored in the
OS-specific keychain by default.
