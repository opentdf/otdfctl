---
title: Add a Public Key to a Key Access Server
command:
  name: create
  aliases:
    - add
  flags:
    - name: kas
      shorthand: k
      description: Key Access Server ID, Name or URI.
      required: true
    - name: key
      shorthand: p
      description: Public key to add to the KAS. Must be in PEM format. Can be base64 encoded or plain text.
      required: true
    - name: key-id
      shorthand: i
      description: ID of the public key.
    - name: algorithm
      shorthand: a
      description: Algorithm of the public key. (rsa:2048, rsa:4096, ec:secp256r1, ec:secp384r1, ec:secp521r1)
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''

---

For more information about registration of Key Access Servers, see the manual for `kas-registry`.

## Example

```shell
otdfctl policy kas-registry public-key add --kas-id 1 --key "-----BEGIN CERTIFICATE-----\nMIIB...5Q=\n-----END CERTIFICATE-----\n" --algorithm rsa:2048
```
