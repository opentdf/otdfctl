---
title: Import Key
command:
  name: import
  aliases:
    - i
  flags:
    - name: key-id
      description: A unique, often human-readable, identifier for the key being imported.
      required: true
    - name: algorithm
      shorthand: a
      description: Algorithm for the key being imported (see table below for options).
      required: true
    - name: kas
      description: Specify the Key Access Server (KAS) where the key will be imported. The KAS can be identified by its ID, URI, or Name.
      required: true
    - name: wrapping-key-id
      description: Identifier related to the wrapping key.
      required: true
    - name: wrapping-key
      shorthand: w
      description: The symmetric key material (AES cipher, hex encoded) used to wrap the imported private key.
      required: true
    - name: private-key-pem
      description: The private key PEM to import (encrypted by an AES 32-byte key, then base64 encoded).
      required: true
    - name: public-key-pem
      shorthand: e
      description: The base64 encoded public key PEM.
      required: true
    - name: label
      shorthand: l
      description: Comma-separated key=value pairs for metadata labels to associate with the imported key (e.g., "owner=team-a,env=production").
---

Imports an existing cryptographic key into a specified Key Access Server (KAS).

>[!IMPORTANT]
>Use this command when migrating keys from KAS over to the platform.
>All keys created with import will be of key_mode=**KEY_MODE_CONFIG_ROOT_KEY**

## Examples

### Import a key

```shell
otdfctl policy kas-registry key import --key-id "imported-key" --algorithm "rsa:2048" \
  --kas 891cfe85-b381-4f85-9699-5f7dbfe2a9ab \
  --wrapping-key-id "my-wrapping-key" \
  --wrapping-key "a8c4824daafcfa38ed0d13002e92b08720e6c4fcee67d52e954c1a6e045907d1" \
  --public-key-pem <base64 encoded public key pem> \
  --private-key-pem <base64 encoded private key pem> \
```

1. The `"algorithm"` specifies the key algorithm:

    | Key Algorithm  |
    | -------------- |
    | `rsa:2048`     |
    | `rsa:4096`     |
    | `ec:secp256r1` |
    | `ec:secp384r1` |
    | `ec:secp521r1` |

2. The `"mode"` specifies where the key that is encrypting TDFs is stored:

    | Mode         | Description                                                                                             |
    | ------------ | ------------------------------------------------------------------------------------------------------- |
    | `local`      | Root Key is stored within Virtru's database and the symmetric wrapping key is stored in KAS             |
    | `provider`   | Root Key is stored within Virtru's database and the symmetric wrapping key is stored externally         |
    | `remote`     | Root Key and wrapping key are stored remotely                                                           |
    | `public_key` | Root Key and wrapping key are stored remotely. Use this when importing another org's policy information |
