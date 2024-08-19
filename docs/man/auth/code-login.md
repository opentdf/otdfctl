---
title: Open a browser and login with Auth Code PKCE

command:
  name: code-login
  flags:
    - name: client-id
      description: A clientId for use in auth code flow (default = platform well-known public_client_id)
      shorthand: i
      required: false
    - name: no-cache
      description: Do not cache credentials on the native OS (print access token to stdout)
      default: false
---

Authenticate for use of the OpenTDF Platform through a browser (required).

Provide a specific public 'client-id' known to support the Auth Code PKCE flow and recognized
by the OpenTDF Platform, or use the default `opentdf-public` client if not specified.

The OIDC Access Token will be stored in the OS-specific keychain by default, otherwise printed to `stdout` if `--no-cache` is passed.
