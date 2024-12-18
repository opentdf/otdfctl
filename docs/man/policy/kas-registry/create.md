---
title: Create a Key Access Server registration
command:
  name: create
  aliases:
    - c
    - add
    - new
  flags:
    - name: uri
      shorthand: u
      description: URI of the Key Access Server
      required: true
    - name: public-keys
      shorthand: c
      description: One or more public keys saved for the KAS
    - name: public-key-remote
      shorthand: r
      description: Remote URI where the public key can be retrieved for the KAS
    - name: label
    - name: name
      shorthand: n
      description: Optional name of the registered KAS (must be unique within policy)
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

Public keys can be stored as either `remote` or `cached` under the following JSON structure.

### Remote

The value passed to the `--public-key-remote` flag puts the hosted location where the public key
can be retrieved for the registered KAS under the `remote` key, such as `https://kas.io/public_key`

### Cached

```json
{
  "cached": {
    // One or more known public keys for the KAS
    "keys": [
      {
        // x509 ASN.1 content in PEM envelope, usually
        "pem": "<your PEM certificate>",
        // key identifier
        "kid": "<your key id>",
        // key algorithm (see table below)
        "alg": 1
      }
    ]
  }
}
```

The JSON value passed to the `--public-keys` flag stores the set of public keys for the KAS.

1. The `"pem"` value should contain the entire certificate `-----BEGIN CERTIFICATE-----\nMIIB...5Q=\n-----END CERTIFICATE-----\n`.

2. The `"kid"` value is a named key identifier, which is useful for key rotations.

3. The `"alg"` specifies the key algorithm:

| Key Algorithm  | `alg` Value |
| -------------- | ----------- |
| `rsa:2048`     | 1           |
| `ec:secp256r1` | 5           |

### Local

Deprecated.

For more information about registration of Key Access Servers, see the manual for `kas-registry`.

## Examples

```shell
otdfctl policy kas-registry create --uri http://example.com/kas --name example-kas --public-keys '{
        "cached": {
          "keys": [
                {
                  "pem": "-----BEGIN CERTIFICATE-----\nMIIC/TCCAeWgAwIBAgIUSHTJ2bzAh7dQmmF03q6Iq/n0l90wDQYJKoZIhvcNAQEL\nBQAwDjEMMAoGA1UEAwwDa2FzMB4XDTI0MDYwNjE3NDY1NFoXDTI1MDYwNjE3NDY1\nNFowDjEMMAoGA1UEAwwDa2FzMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKC\nAQEAxN3APihTiojcaH6oWj1tMtZMaaZ+IA1qtqFmpy5Fg8D5bEsP736GxzUMFsMV\nshrKEXz8dY9Kp23uIwyeC0RPWLe5xIfTkJUbyLpqGdlEgqj10RQ8kSVq270XPES2\nGZUij2DuJVfwpTpLzcti2PsgEOoOKC6NnnAI0NS1mao/2DxQxs/D9hAJjGdpzymb\nxi2TxGnvYbvofCPd8RdFTCPvgwKLS7+MqBcmic9VdX91QNOPmrP3rIoKtjjd+5PY\nl/z73PAxR3K3SIzIZLvItq2ahobOOMiSxw8soOlOdHNUJTpECcduhRbquqmK6fTw\nVOfrcRQhhU4TkDu92LI7SglOWQIDAQABo1MwUTAdBgNVHQ4EFgQUdgxx7U5AQgfi\niQWu3khi9yneEVowHwYDVR0jBBgwFoAUdgxx7U5AQgfiiQWu3khi9yneEVowDwYD\nVR0TAQH/BAUwAwEB/zANBgkqhkiG9w0BAQsFAAOCAQEATcLYbHomJgLQ/H6iDvcA\nIpISF/Rcxgh7NnIqRkB+Tm4xNlNHIxl4Sz+KkEZEPh0WKItGVDj3293rArROEOXI\ntVmn2OBv9M/5DQkHj76Ru4PQ2TcL0CACl1JKfqXLsMc6HHTp8ZTP8lMdpW4kzEc3\nfVtgvtpJc4WHdUIEzAtTlzYRqIbyyBMWeTjXwa54aMv3RZQdJ+C0ehwWTDQDph7n\nKY3+7G0enNEVtyW4dtxvQQbidMany0JEpr6QpPmxC8e0Z23dMDdkR1IoT99PhdW/\nQC8xMjuLCiREV7a6e2MxCGj3fxrnMXwOIqO3AzNswe2amcoz2ktuoqgDTYlo+FkK\n5w==\n-----END CERTIFICATE-----\n",
                  "kid": "k1",
                  "alg": 1
                }
          ]
        }
  }'
```

```shell
  SUCCESS   Created kas-registry: 62857b55-560c-4b67-96e3-33e4670ecb3b                                                                                                                                                                                                                                                                                                                                                             
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

With a remote public key:
```shell
otdfctl policy kas-registry create --uri http://example.com/kas2 --name example-kas2 --public-key-remote "https://example.com/kas2/public_key"
```

```shell
  SUCCESS   Created kas-registry: 3c39618a-cd8c-48cf-a60c-e8a2f4be4dd5                                                                                                                                                                                                                                                                                                                                                                     
╭─────────────────────────────────────────────────────────────────────────┬────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╮
│Property                                                                 │Value                                                                                                                                           │
├─────────────────────────────────────────────────────────────────────────┼────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────┤
│Id                                                                       │3c39618a-cd8c-48cf-a60c-e8a2f4be4dd5                                                                                                            │
│URI                                                                      │http://example.com/kas2                                                                                                                         │
│PublicKey                                                                │remote:"https://example.com/kas2/public_key"                                                                                                    │
│Name                                                                     │example-kas2                                                                                                                                    │
│Created At                                                               │Wed Dec 18 04:57:51 UTC 2024                                                                                                                    │
│Updated At                                                               │Wed Dec 18 04:57:51 UTC 2024                                                                                                                    │
╰─────────────────────────────────────────────────────────────────────────┴────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────────╯
  NOTE   Use 'otdfctl policy kas-registry get --id=3c39618a-cd8c-48cf-a60c-e8a2f4be4dd5 --json' to see all properties 
```
