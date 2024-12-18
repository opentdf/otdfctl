---
title: Inspect a TDF file
command:
  name: inspect [file]
  flags:
---

# Inspect a TDF file

Prints the `manifest.json` of the specified TDF for inspection.

This is useful for development and administration.

## Example

```shell
$ otdfctl inspect example.tdf

{
  "manifest": {
    "algorithm": "HS256",
    "keyAccessType": "split",
    "mimeType": "",
    "policy": "eyJ1dWlkIjoiOTk0OWJkYTItN2E4MC00NTcwLWJjMTYtMjAxYmM4ZDA3YTE2IiwiYm9keSI6eyJkYXRhQXR0cmlidXRlcyI6W10sImRpc3NlbSI6W119fQ==",
    "protocol": "zip",
    "segmentHashAlgorithm": "GMAC",
    "signature": "MmEyZTIwYzgwYzIyMGNmMDMzNmQ0Y2U2MzU3Y2Q4YTRjYjFkYmNiNzQ0YzFhNjhlNjU0Y2MxNGM3MmMwYzNiZQ==",
    "type": "reference",
    "method": {
      "algorithm": "AES-256-GCM",
      "iv": "lUpBdhH8csdbqgAP",
      "isStreamable": true
    },
    "integrityInformation": {
      "rootSignature": {
        "alg": "HS256",
        "sig": "MmEyZTIwYzgwYzIyMGNmMDMzNmQ0Y2U2MzU3Y2Q4YTRjYjFkYmNiNzQ0YzFhNjhlNjU0Y2MxNGM3MmMwYzNiZQ=="
      },
      "segmentHashAlg": "GMAC",
      "segmentSizeDefault": 1048576,
      "encryptedSegmentSizeDefault": 1048604,
      "segments": [
        {
          "hash": "Y2RhNWYwMmFhNWE4M2EyYWY5Zjk2OTQ5NjU1MGQ4ODY=",
          "segmentSize": 1618,
          "encryptedSegmentSize": 1646
        }
      ]
    },
    "encryptionInformation": {
      "type": "split",
      "policy": "eyJ1dWlkIjoiOTk0OWJkYTItN2E4MC00NTcwLWJjMTYtMjAxYmM4ZDA3YTE2IiwiYm9keSI6eyJkYXRhQXR0cmlidXRlcyI6W10sImRpc3NlbSI6W119fQ==",
      "keyAccess": [
        {
          "type": "wrapped",
          "url": "http://localhost:8080/kas",
          "protocol": "kas",
          "wrappedKey": "eEjzpg2XloommzdT6b9EVue6q1Lq/MRoZH9pU7EhcKpmt/+w6VHOUrTfk7rD05orQ2T2s2CjajrT6JNTbwQPXeoGCkKVp2xy2xceuNn8GFRJ5Gfz5rm1yI2vuOcn9xX4xbIHeLHQb7tUHyZnpeDMPc0y222VQfu/3Js1ycOBLE6lmgTgU3fXMYWSwXUIIdvWkrCW43eQxCPwZIO3HCOCo7mpWw/1gnzgJSldH/8vnlqeyeQDOvNq3+TDUwk74BV+0O72SAycaPISe/Vhh4SwSpUnRJdRN5mSngD9iuB/Dd9ChbhmNuwPW9KDzFocyz/SM5GsU3jhmjntMGNCMviR6g==",
          "policyBinding": "ODViMjE5N2NiNWQzOWVmZDk0ZmU0OTMxMTM4MDNjNjNlMmZlNGQxYWE2NzIyYTQ3YmRhMTI1NGRhZTdkMmQ5NQ==",
          "encryptedMetadata": "eyJjaXBoZXJ0ZXh0IjoibFVwQmRoSDhjc2RicWdBUGwxYkxtOW9kSHVReCtQclFxbUx3R3c9PSIsIml2IjoibFVwQmRoSDhjc2RicWdBUCJ9"
        }
      ],
      "method": {
        "algorithm": "AES-256-GCM",
        "iv": "lUpBdhH8csdbqgAP",
        "isStreamable": true
      },
      "integrityInformation": {
        "rootSignature": {
          "alg": "HS256",
          "sig": "MmEyZTIwYzgwYzIyMGNmMDMzNmQ0Y2U2MzU3Y2Q4YTRjYjFkYmNiNzQ0YzFhNjhlNjU0Y2MxNGM3MmMwYzNiZQ=="
        },
        "segmentHashAlg": "GMAC",
        "segmentSizeDefault": 1048576,
        "encryptedSegmentSizeDefault": 1048604,
        "segments": [
          {
            "hash": "Y2RhNWYwMmFhNWE4M2EyYWY5Zjk2OTQ5NjU1MGQ4ODY=",
            "segmentSize": 1618,
            "encryptedSegmentSize": 1646
          }
        ]
      }
    }
  },
  "attributes": []
}
```
