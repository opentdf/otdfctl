---
title: Set up client credentials

command:
  name: client-credentials
  flags:
    - name: clientId
      description: Client ID
      required: true
    - name: clientSecret
      description: Client secret
---

Allows the user to login in via clientId and clientSecret. This will subsequently be stored in the
OS-specific keychain by default.
