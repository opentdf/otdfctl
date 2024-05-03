---
title: Authenticate to the platform with the client-credentials flow

command:
  name: client-credentials
  flags:
    - name: client-id
      description: A clientId for use in client-credentials auth flow
      shorthand: i
      required: true
    - name: client-secret
      description: A clientSecret for use in client-credentials auth flow
      shorthand: s
    - name: no-cache
      description: Do not cache credentials on the native OS and print access token to stdout instead
---

Allows the user to login in via Client ID and Secret. The client credentials and OIDC Access Token will be stored
in the OS-specific keychain by default, otherwise printed to `stdout` if `--no-cache` is passed.
