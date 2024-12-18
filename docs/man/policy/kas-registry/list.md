---
title: List Key Access Server registrations
command:
  name: list
  aliases:
    - l
  flags:
    - name: limit
      shorthand: l
      description: Limit retrieved count
    - name: offset
      shorthand: o
      description: Offset (page) quantity from start of the list
---

For more information about registration of Key Access Servers, see the manual for `kas-registry`.

## Example

```shell
otdfctl policy kas-registry list
```

```shell
  SUCCESS   Found kas-registry list                                                                                                                                                                                         
                                                                                                                                                                                                                            
╭──────────────────────────────────────────────────────────────────┬─────────────────────────────────────────────────────┬────────────────────────────────────────┬────────────────────────────────────────────────────────╮
│ID                                                                │URI                                                  │Name                                    │PublicKey                                               │
├──────────────────────────────────────────────────────────────────┼─────────────────────────────────────────────────────┼────────────────────────────────────────┼────────────────────────────────────────────────────────┤
│f612b628-5459-4342-b20f-3768b30ad588                              │http://localhost:8080/kas                            │alpha                                   │cached:{keys:{pem:"-----BEGIN PUBLIC KEY-----\\nMIIBIjA…│
│62857b55-560c-4b67-96e3-33e4670ecb3b                              │http://example.com/kas                               │example-kas                             │cached:{keys:{pem:"-----BEGIN CERTIFICATE-----\nMIIC/TC…│
│3c39618a-cd8c-48cf-a60c-e8a2f4be4dd5                              │http://example.com/kas2                              │example-kas2                            │remote:"https://example.com/kas2/public_key"            │
╰──────────────────────────────────────────────────────────────────┴─────────────────────────────────────────────────────┴────────────────────────────────────────┴────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy kas-registry get --id=<id> --json' to see all properties
```
