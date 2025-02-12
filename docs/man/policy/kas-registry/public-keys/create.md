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

Add a public key to a Key Access Server. The public key must be in PEM format. It can be base64 encoded or plain text.

If a key exists with the same algorithm already, the new key will be marked as active and the existing key will be marked as inactive. The namespace, attribute and value mappings will be updated to point to the new key.

## Example

```shell
# Add a public key to a Key Access Server By ID
otdfctl policy kas-registry public-key create --kas 62857b55-560c-4b67-96e3-33e4670ecb3b  --key-id key-1 --key "-----BEGIN CERTIFICATE-----\nMIIB...5Q=\n-----END CERTIFICATE-----\n" --algorithm rsa:2048
```

```shell
# Add a public key to a Key Access Server By Name
otdfctl policy kas-registry public-key
create --kas kas-1 --key-id key-1 --key "-----BEGIN CERTIFICATE-----\nMIIB...5Q=\n-----END CERTIFICATE-----\n" --algorithm rsa:2048
```

```shell
# Add a public key to a Key Access Server By URI
otdfctl policy kas-registry public-key
create --kas https://example.com/kas --key-id key-1 --key "-----BEGIN CERTIFICATE-----\nMIIB...5Q=\n-----END CERTIFICATE-----\n" --algorithm rsa:2048
```

