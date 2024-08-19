---
title: Print the cached OIDC access token (if found)

command:
  name: print-access-token
  flags:
    - name: json
      description: Print the full token in JSON format
      default: false
---

Retrieves a new OIDC Access Token using the client credentials from the OS-specific keychain and prints to stdout if found.
