---
title: Encrypt file or stdin as a TDF
command:
  name: encrypt [file]
  flags:
    - name: out
      shorthand: o
      description: The output file TDF in the current working directory instead of stdout ('-o file.txt' and '-o file.txt.tdf' both write the TDF as file.txt.tdf).
      default: ''
    - name: attr
      shorthand: a
      description: Attribute value Fully Qualified Names (FQNs, i.e. 'https://example.com/attr/attr1/value/value1') to apply to the encrypted data.
    - name: mime-type
      description: The MIME type of the input data. If not provided, the MIME type is inferred from the input data.
    - name: tdf-type
      shorthand: t
      description: The type of tdf to encrypt as. TDF3 supports structured manifests and larger payloads. Nano has a smaller footprint and more performant, but does not support structured manifests or large payloads.
      enum:
        - tdf3
        - nano
      default: tdf3
    - name: kas-url-path
      description: URL path to the KAS service at the platform endpoint domain. Leading slash is required if needed.
      default: /kas
---

Build a Trusted Data Format (TDF) with encrypted content from a specified file or input from stdin utilizing OpenTDF platform.

## Examples:

```bash
# output to stdout
echo "some text" | otdfctl encrypt
otdfctl encrypt hello.txt
# pipe stdout to a bucket
echo "my secret" | otdfctl encrypt | aws s3 cp - s3://my-bucket/secret.txt.tdf

# output hello.txt.tdf in root directory
echo "hello world" | otdfctl encrypt -o hello.txt
cat hello.txt | otdfctl encrypt -o hello.txt
cat hello.txt | otdfctl encrypt -o hello.txt.tdf #.tdf extension is only added once
```
