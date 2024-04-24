---
title: Encrypt file or stdin as a TDF
command:
  name: encrypt [file]
  flags:
    - name: out
      shorthand: o
      description: A filename and extension that will be TDFd (i.e. '-o password.txt' -> 'password.txt.tdf') and placed in the current working directory.
      default: ''
    - name: attr
      shorthand: a
      description: Attribute value Fully Qualified Names (FQNs, i.e. 'https://example.com/attr/attr1/value/value1') to apply to the encrypted data.
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
