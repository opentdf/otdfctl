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
    - name: public-key-cached
      shorthand: c
      description: Cached public keys for the KAS
    - name: public-key-remote
      shorthand: r
      description: Remote URI where the public key can be retrieved for the KAS
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

For more information about registration of Key Access Servers, see the manual for `kas-registry`.

Public keys can be stored as either `remote` or `cached` under the following JSON structure.

### Remote

```json
{ "remote": "https://mykas.com/public_key" }
```

The JSON value passed to the `--public-key-remote` flag puts the location where the public key
can be accessed for a the registered KAS under the `remote` key.

### Cached

```json5
{
  "cached": {
    // One or more known public keys for the KAS
    "keys":[
      {
        // x509 ASN.1 content in PEM envelope, usually
        "pem": "base64encodedCert",
        // key identifier 
        "kid": "<your key id>",
        // algorithm (either: 1 for rsa:2048, 2 for ec:secp256r1)
        "alg": 1
      }
    ]
  }
}
```

The JSON value passed to the `--public-key-cached` flag stores the set of public keys for the KAS.

The PEM base64 encoding should contain everything `-----BEGIN CERTIFICATE-----\nMIIB...5Q=\n-----END CERTIFICATE-----\n`.

### Local

Deprecated.
