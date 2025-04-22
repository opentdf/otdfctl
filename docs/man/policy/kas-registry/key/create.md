---
title: Create Key
command:
  name: create
  aliases:
    - c
  flags:
    - name: kasId
      description: Key Access Server ID
      required: true
    - name: keyId
      description: Name for the Key
      required: true
    - name: alg
      shorthand: a
      description: Algorithm for the key
      required: true
    - name: mode
      shorthand: m
      description: Describes how the private key is managed
      required: true
    - name: publicKeyCtx
      description: Public Key Context in JSON form, needs to contains the public key for remote KEKs. Public key is generated automatically internal keys.
    - name: privateKeyCtx
      description: Private Key Context in JSON form. Private Key is generated automatically for internal keys.
    - name: wrappingKey
      shorthand: w
      description: The key used to wrap the generated private key. (Must be generated with AES cipher, and base64 encoded)
    - name: providerConfigId
      shorthand: p
      description: Configuration ID for the key provider, if applicable
    - name: label
      shorthand: l
      description: Metadata labels for the provider config 
---

Creates a new key that for a specified Key Access Server, which will be used
for encrypting and decrypting data keys.

## Examples

```shell
otdfctl key create --kasId 891cfe85-b381-4f85-9699-5f7dbfe2a9ab --keyId "aws-key" --algorithm 1 --mode "remote" --wrappingKey
```

1. The `"alg"` specifies the key algorithm:

    | Key Algorithm  |
    | -------------- |
    | `rsa:2048`     |
    | `rsa:4096`     |
    | `ec:secp256r1` |
    | `ec:secp384r1` |
    | `ec:secp521r1` |

2. The `"mode"` specifies whether the KEK is stored in Virtru's database or remotely in an external KMS.

    | Mode           |
    | -------------- |
    | `local`        |
    | `remote`       |
