---
title: Open a browser and login

command:
  name: login
  flags:
    - name: client-id
      description: A clientId for use in auth code flow (e.g. cli-client)
      shorthand: i
      required: true
---

> [!NOTE]
> Requires experimental profiles feature.
>
> | OS | Keychain | State |
> | --- | --- | --- |
> | MacOS | Keychain | Stable |
> | Windows | Credential Manager | Alpha |
> | Linux | Secret Service | Not yet supported |

Authenticate for use of the OpenTDF Platform through a browser (required).

Provide a specific public 'client-id' known to support the Auth Code PKCE flow and recognized
by the OpenTDF Platform, or use the default public client in the platform well-known configuration if not specified.

The OIDC Access Token will be stored in the OS-specific keychain by default (Linux not yet supported).
