---
title: Create Key
command:
  name: create
  aliases:
    - c
  flags:
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
    - name: kasId
      description: Key Access Server ID
    - name: kasUri
      description: The FQN of the KAS server
    - name: kasName
      description: Key Access Server name
    - name: wrappingKeyId
      description: KeyId that wraps the asymmetric key.
    - name: wrappingKey
      shorthand: w
      description: The key used to wrap the generated private key. (Must be generated with AES cipher, and base64 encoded)
    - name: privatePem
      description: The private key pem, encrypted by an AES 32 byte key, base64 encoded.
    - name: providerConfigId
      shorthand: p
      description: Configuration ID for the key provider, if applicable. When mode is `"provider"` we will use the provider to wrap the key.
    - name: pubPem
      shorthand: e
      description: The base64 encoded public key pem to be used with `"remote"` and `"public_key"` modes.
    - name: label
      shorthand: l
      description: Metadata labels for the provider config 
---

Creates a new key that for a specified Key Access Server, which will be used
for encrypting and decrypting data keys.

## Examples

```shell
otdfctl key create --keyId "aws-key" --alg "rsa:2048" --mode "local" --kasId 891cfe85-b381-4f85-9699-5f7dbfe2a9ab --wrappingKeyId "virtru-stored-key" --wrappingKey "YWVzIGtleQ=="

otdfctl key create --keyId "aws-key" --alg "rsa:2048" --mode "local" --kasUri "https://test-kas.com" --wrappingKeyId "virtru-stored-key" --wrappingKey "YWVzIGtleQ=="
```

```shell
otdfctl key create --keyId "aws-key" --alg "rsa:2048" --mode "provider" --kasUri "https://test-kas.com" --pubPem "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tXG5NSUlDL1RDQ0FlV2dBd0lCQWdJVVNIVEoyYnpBaDdkUW1tRjAzcTZJcS9uMGw5MHdEUVlKS29aSWh2Y05BUUVMXG5CUUF3RGpFTU1Bb0dBMVVFQXd3RGEyRnpNQjRYRFRJME1EWXdOakUzTkRZMU5Gb1hEVEkxTURZd05qRTNORFkxXG5ORm93RGpFTU1Bb0dBMVVFQXd3RGEyRnpNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDXG5BUUVBeE4zQVBpaFRpb2pjYUg2b1dqMXRNdFpNYWFaK0lBMXF0cUZtcHk1Rmc4RDViRXNQNzM2R3h6VU1Gc01WXG5zaHJLRVh6OGRZOUtwMjN1SXd5ZUMwUlBXTGU1eElmVGtKVWJ5THBxR2RsRWdxajEwUlE4a1NWcTI3MFhQRVMyXG5HWlVpajJEdUpWZndwVHBMemN0aTJQc2dFT29PS0M2Tm5uQUkwTlMxbWFvLzJEeFF4cy9EOWhBSmpHZHB6eW1iXG54aTJUeEdudllidm9mQ1BkOFJkRlRDUHZnd0tMUzcrTXFCY21pYzlWZFg5MVFOT1BtclAzcklvS3RqamQrNVBZXG5sL3o3M1BBeFIzSzNTSXpJWkx2SXRxMmFob2JPT01pU3h3OHNvT2xPZEhOVUpUcEVDY2R1aFJicXVxbUs2ZlR3XG5WT2ZyY1JRaGhVNFRrRHU5MkxJN1NnbE9XUUlEQVFBQm8xTXdVVEFkQmdOVkhRNEVGZ1FVZGd4eDdVNUFRZ2ZpXG5pUVd1M2toaTl5bmVFVm93SHdZRFZSMGpCQmd3Rm9BVWRneHg3VTVBUWdmaWlRV3Uza2hpOXluZUVWb3dEd1lEXG5WUjBUQVFIL0JBVXdBd0VCL3pBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQVRjTFliSG9tSmdMUS9INmlEdmNBXG5JcElTRi9SY3hnaDdObklxUmtCK1RtNHhObE5ISXhsNFN6K0trRVpFUGgwV0tJdEdWRGozMjkzckFyUk9FT1hJXG50Vm1uMk9CdjlNLzVEUWtIajc2UnU0UFEyVGNMMENBQ2wxSktmcVhMc01jNkhIVHA4WlRQOGxNZHBXNGt6RWMzXG5mVnRndnRwSmM0V0hkVUlFekF0VGx6WVJxSWJ5eUJNV2VUalh3YTU0YU12M1JaUWRKK0MwZWh3V1REUURwaDduXG5LWTMrN0cwZW5ORVZ0eVc0ZHR4dlFRYmlkTWFueTBKRXByNlFwUG14QzhlMFoyM2RNRGRrUjFJb1Q5OVBoZFcvXG5RQzh4TWp1TENpUkVWN2E2ZTJNeENHajNmeHJuTVh3T0lxTzNBek5zd2UyYW1jb3oya3R1b3FnRFRZbG8rRmtLXG41dz09XG4tLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tXG4=" --privatePem "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tXG5NSUlDL1RDQ0FlV2dBd0lCQWdJVVNIVEoyYnpBaDdkUW1tRjAzcTZJcS9uMGw5MHdEUVlKS29aSWh2Y05BUUVMXG5CUUF3RGpFTU1Bb0dBMVVFQXd3RGEyRnpNQjRYRFRJME1EWXdOakUzTkRZMU5Gb1hEVEkxTURZd05qRTNORFkxXG5ORm93RGpFTU1Bb0dBMVVFQXd3RGEyRnpNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDXG5BUUVBeE4zQVBpaFRpb2pjYUg2b1dqMXRNdFpNYWFaK0lBMXF0cUZtcHk1Rmc4RDViRXNQNzM2R3h6VU1Gc01WXG5zaHJLRVh6OGRZOUtwMjN1SXd5ZUMwUlBXTGU1eElmVGtKVWJ5THBxR2RsRWdxajEwUlE4a1NWcTI3MFhQRVMyXG5HWlVpajJEdUpWZndwVHBMemN0aTJQc2dFT29PS0M2Tm5uQUkwTlMxbWFvLzJEeFF4cy9EOWhBSmpHZHB6eW1iXG54aTJUeEdudllidm9mQ1BkOFJkRlRDUHZnd0tMUzcrTXFCY21pYzlWZFg5MVFOT1BtclAzcklvS3RqamQrNVBZXG5sL3o3M1BBeFIzSzNTSXpJWkx2SXRxMmFob2JPT01pU3h3OHNvT2xPZEhOVUpUcEVDY2R1aFJicXVxbUs2ZlR3XG5WT2ZyY1JRaGhVNFRrRHU5MkxJN1NnbE9XUUlEQVFBQm8xTXdVVEFkQmdOVkhRNEVGZ1FVZGd4eDdVNUFRZ2ZpXG5pUVd1M2toaTl5bmVFVm93SHdZRFZSMGpCQmd3Rm9BVWRneHg3VTVBUWdmaWlRV3Uza2hpOXluZUVWb3dEd1lEXG5WUjBUQVFIL0JBVXdBd0VCL3pBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQVRjTFliSG9tSmdMUS9INmlEdmNBXG5JcElTRi9SY3hnaDdObklxUmtCK1RtNHhObE5ISXhsNFN6K0trRVpFUGgwV0tJdEdWRGozMjkzckFyUk9FT1hJXG50Vm1uMk9CdjlNLzVEUWtIajc2UnU0UFEyVGNMMENBQ2wxSktmcVhMc01jNkhIVHA4WlRQOGxNZHBXNGt6RWMzXG5mVnRndnRwSmM0V0hkVUlFekF0VGx6WVJxSWJ5eUJNV2VUalh3YTU0YU12M1JaUWRKK0MwZWh3V1REUURwaDduXG5LWTMrN0cwZW5ORVZ0eVc0ZHR4dlFRYmlkTWFueTBKRXByNlFwUG14QzhlMFoyM2RNRGRrUjFJb1Q5OVBoZFcvXG5RQzh4TWp1TENpUkVWN2E2ZTJNeENHajNmeHJuTVh3T0lxTzNBek5zd2UyYW1jb3oya3R1b3FnRFRZbG8rRmtLXG41dz09XG4tLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tXG4=" --wrappingKeyId "openbao-key" --providerConfigId "f86b166a-98a5-407a-939f-ef84916ce1e5"
```

```shell
otdfctl key create --keyId "aws-key" --alg "rsa:2048" --mode "remote" --kasUri "https://test-kas.com" --wrappingKeyId "openbao-key" --providerConfigId "f86b166a-98a5-407a-939f-ef84916ce1e5" --pubPem "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tXG5NSUlDL1RDQ0FlV2dBd0lCQWdJVVNIVEoyYnpBaDdkUW1tRjAzcTZJcS9uMGw5MHdEUVlKS29aSWh2Y05BUUVMXG5CUUF3RGpFTU1Bb0dBMVVFQXd3RGEyRnpNQjRYRFRJME1EWXdOakUzTkRZMU5Gb1hEVEkxTURZd05qRTNORFkxXG5ORm93RGpFTU1Bb0dBMVVFQXd3RGEyRnpNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDXG5BUUVBeE4zQVBpaFRpb2pjYUg2b1dqMXRNdFpNYWFaK0lBMXF0cUZtcHk1Rmc4RDViRXNQNzM2R3h6VU1Gc01WXG5zaHJLRVh6OGRZOUtwMjN1SXd5ZUMwUlBXTGU1eElmVGtKVWJ5THBxR2RsRWdxajEwUlE4a1NWcTI3MFhQRVMyXG5HWlVpajJEdUpWZndwVHBMemN0aTJQc2dFT29PS0M2Tm5uQUkwTlMxbWFvLzJEeFF4cy9EOWhBSmpHZHB6eW1iXG54aTJUeEdudllidm9mQ1BkOFJkRlRDUHZnd0tMUzcrTXFCY21pYzlWZFg5MVFOT1BtclAzcklvS3RqamQrNVBZXG5sL3o3M1BBeFIzSzNTSXpJWkx2SXRxMmFob2JPT01pU3h3OHNvT2xPZEhOVUpUcEVDY2R1aFJicXVxbUs2ZlR3XG5WT2ZyY1JRaGhVNFRrRHU5MkxJN1NnbE9XUUlEQVFBQm8xTXdVVEFkQmdOVkhRNEVGZ1FVZGd4eDdVNUFRZ2ZpXG5pUVd1M2toaTl5bmVFVm93SHdZRFZSMGpCQmd3Rm9BVWRneHg3VTVBUWdmaWlRV3Uza2hpOXluZUVWb3dEd1lEXG5WUjBUQVFIL0JBVXdBd0VCL3pBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQVRjTFliSG9tSmdMUS9INmlEdmNBXG5JcElTRi9SY3hnaDdObklxUmtCK1RtNHhObE5ISXhsNFN6K0trRVpFUGgwV0tJdEdWRGozMjkzckFyUk9FT1hJXG50Vm1uMk9CdjlNLzVEUWtIajc2UnU0UFEyVGNMMENBQ2wxSktmcVhMc01jNkhIVHA4WlRQOGxNZHBXNGt6RWMzXG5mVnRndnRwSmM0V0hkVUlFekF0VGx6WVJxSWJ5eUJNV2VUalh3YTU0YU12M1JaUWRKK0MwZWh3V1REUURwaDduXG5LWTMrN0cwZW5ORVZ0eVc0ZHR4dlFRYmlkTWFueTBKRXByNlFwUG14QzhlMFoyM2RNRGRrUjFJb1Q5OVBoZFcvXG5RQzh4TWp1TENpUkVWN2E2ZTJNeENHajNmeHJuTVh3T0lxTzNBek5zd2UyYW1jb3oya3R1b3FnRFRZbG8rRmtLXG41dz09XG4tLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tXG4="
```

```shell
otdfctl key create --keyId "aws-key" --alg "rsa:2048" --mode "public_key" --kasUri "https://test-kas.com" --pubPem "LS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tXG5NSUlDL1RDQ0FlV2dBd0lCQWdJVVNIVEoyYnpBaDdkUW1tRjAzcTZJcS9uMGw5MHdEUVlKS29aSWh2Y05BUUVMXG5CUUF3RGpFTU1Bb0dBMVVFQXd3RGEyRnpNQjRYRFRJME1EWXdOakUzTkRZMU5Gb1hEVEkxTURZd05qRTNORFkxXG5ORm93RGpFTU1Bb0dBMVVFQXd3RGEyRnpNSUlCSWpBTkJna3Foa2lHOXcwQkFRRUZBQU9DQVE4QU1JSUJDZ0tDXG5BUUVBeE4zQVBpaFRpb2pjYUg2b1dqMXRNdFpNYWFaK0lBMXF0cUZtcHk1Rmc4RDViRXNQNzM2R3h6VU1Gc01WXG5zaHJLRVh6OGRZOUtwMjN1SXd5ZUMwUlBXTGU1eElmVGtKVWJ5THBxR2RsRWdxajEwUlE4a1NWcTI3MFhQRVMyXG5HWlVpajJEdUpWZndwVHBMemN0aTJQc2dFT29PS0M2Tm5uQUkwTlMxbWFvLzJEeFF4cy9EOWhBSmpHZHB6eW1iXG54aTJUeEdudllidm9mQ1BkOFJkRlRDUHZnd0tMUzcrTXFCY21pYzlWZFg5MVFOT1BtclAzcklvS3RqamQrNVBZXG5sL3o3M1BBeFIzSzNTSXpJWkx2SXRxMmFob2JPT01pU3h3OHNvT2xPZEhOVUpUcEVDY2R1aFJicXVxbUs2ZlR3XG5WT2ZyY1JRaGhVNFRrRHU5MkxJN1NnbE9XUUlEQVFBQm8xTXdVVEFkQmdOVkhRNEVGZ1FVZGd4eDdVNUFRZ2ZpXG5pUVd1M2toaTl5bmVFVm93SHdZRFZSMGpCQmd3Rm9BVWRneHg3VTVBUWdmaWlRV3Uza2hpOXluZUVWb3dEd1lEXG5WUjBUQVFIL0JBVXdBd0VCL3pBTkJna3Foa2lHOXcwQkFRc0ZBQU9DQVFFQVRjTFliSG9tSmdMUS9INmlEdmNBXG5JcElTRi9SY3hnaDdObklxUmtCK1RtNHhObE5ISXhsNFN6K0trRVpFUGgwV0tJdEdWRGozMjkzckFyUk9FT1hJXG50Vm1uMk9CdjlNLzVEUWtIajc2UnU0UFEyVGNMMENBQ2wxSktmcVhMc01jNkhIVHA4WlRQOGxNZHBXNGt6RWMzXG5mVnRndnRwSmM0V0hkVUlFekF0VGx6WVJxSWJ5eUJNV2VUalh3YTU0YU12M1JaUWRKK0MwZWh3V1REUURwaDduXG5LWTMrN0cwZW5ORVZ0eVc0ZHR4dlFRYmlkTWFueTBKRXByNlFwUG14QzhlMFoyM2RNRGRrUjFJb1Q5OVBoZFcvXG5RQzh4TWp1TENpUkVWN2E2ZTJNeENHajNmeHJuTVh3T0lxTzNBek5zd2UyYW1jb3oya3R1b3FnRFRZbG8rRmtLXG41dz09XG4tLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tXG4="
```

1. The `"alg"` specifies the key algorithm:

    | Key Algorithm  |
    | -------------- |
    | `rsa:2048`     |
    | `rsa:4096`     |
    | `ec:secp256r1` |
    | `ec:secp384r1` |
    | `ec:secp521r1` |

2. The `"mode"` specifies where the key that is encrypting TDFs is stored. All keys will be encrypted when stored in Virtru's DB, for modes `"local"` and `"provider"`

    | Mode             | Description                                                                                             |
    | ---------------- | ------------------------------------------------------------------------------------------------------- |
    | `local`          | Root Key is stored within Virtru's database and the symmetric wrapping key is stored in KAS             |
    | `provider`       | Root Key is stored within Virtru's database and the symmetric wrapping key is stored externally         |
    | `remote`         | Root Key and wrapping key are stored remotely                                                           |
    | `public_key`     | Root Key and wrapping key are stored remotely. Use this when importing another org's policy information |
