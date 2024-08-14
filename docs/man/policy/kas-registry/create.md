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
    - name: public-key-local
      shorthand: p
      description: Public key of the Key Access Server
    - name: public-key-remote
      shorthand: r
      description: URI of the public key of the Key Access Server
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

For more information about registration of Key Access Servers, see the manual for `kas-registry`.

> Warning: storage of the public key as `remote` or `local` may soon be deprecated in
> favor of reaching out to the KAS directly for the public key.

Public keys can be stored as either `remote` or `local` under the following JSON structure.

### Remote

```json
{ "remote": "https://mykas.com/public_key" }
```

The JSON value passed to the `--public-key-remote` flag puts the location where the public key
can be accessed for a the registered KAS under the `remote` key.

### Local

```json
{ "local": "myBase64EncodedCert" }
```

The JSON value passed to the `--public-key-local` flag puts a base64-encoded key value under
the `local` key.

The base64 encoding should contain everything `-----BEGIN CERTIFICATE-----\nMIIB...5Q=\n-----END CERTIFICATE-----\n`.
