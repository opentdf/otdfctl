---
title: Authenticate to the platform with the client-credentials flow

command:
  name: client-credentials
  flags:
    - name: client-id
      description: Client ID
      shorthand: i
      required: true
    - name: client-secret
      description: Client secret
      shorthand: s
    - name: no-cache
      description: Do not cache credentials on the native OS and print token value to stdout instead
---

Allows the user to login in via Client ID and Secret. The client credentials and OIDC Access Token will be stored
in the OS-specific keychain by default, otherwise printed to `stdout` if `--no-cache` is passed.
