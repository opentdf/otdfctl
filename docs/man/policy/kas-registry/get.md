---
title: Get a registered Key Access Server
command:
  name: get
  aliases:
    - g
  flags:
    - name: id
      shorthand: i
      description: ID of the Key Access Server registration
      required: true
---

For more information about registration of Key Access Servers, see the manual for `kas-registry`.

## Example

```shell
otdfctl policy kas-registry get --id=62857b55-560c-4b67-96e3-33e4670ecb3b
```

```shell
  SUCCESS   Found kas-registry: 62857b55-560c-4b67-96e3-33e4670ecb3b                                                                                                                                                                                                                                                                                                                                                           
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │62857b55-560c-4b67-96e3-33e4670ecb3b                                                                                                            │
│URI                                                                      │http://example.com/kas                                                                                                                          │
│PublicKey                                                                │cached:{keys:{pem:"-----BEGIN CERTIFICATE-----\nMIIC/TCCAeWgAwIBAgIUSHTJ2bzAh7dQmmF03q6Iq/n0l90wDQYJKoZIhvcNAQEL\nBQAwDjEMMAoGA1UEAwwDa2FzMB4XD…│
│Name                                                                     │example-kas                                                                                                                                     │
│Created At                                                               │Wed Dec 18 04:51:22 UTC 2024                                                                                                                    │
│Updated At                                                               │Wed Dec 18 04:51:22 UTC 2024                                                                                                                    │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy kas-registry get --id=62857b55-560c-4b67-96e3-33e4670ecb3b --json' to see all properties 
```
