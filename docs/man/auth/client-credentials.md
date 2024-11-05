---
title: Authenticate to the platform with the client-credentials flow

command:
  name: client-credentials
  args: 
    - client-id
  arbitrary_args:
    - client-secret
---

> [!NOTE]
> Requires experimental profiles feature.
>
> | OS | Keychain | State |
> | --- | --- | --- |
> | MacOS | Keychain | Stable |
> | Windows | Credential Manager | Alpha |
> | Linux | Secret Service | Not yet supported |

Allows the user to login in via Client Credentials flow. The client credentials will be stored safely
in the OS keyring for future use.

## Examples

Authenticate with client credentials (secret provided interactively)

```shell
opentdf auth client-credentials --client-id <client-id>
```

Authenticate with client credentials (secret provided as argument)

```shell
opentdf auth client-credentials --client-id <client-id> --client-secret <client-secret>
```
