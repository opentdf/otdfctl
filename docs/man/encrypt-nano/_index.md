---
title: Encrypt file or stdin as a NanoTDF
command:
  name: encrypt-nano [file]
  flags:
    - name: out
      shorthand: o
      description: The output file NanoTDF in the current working directory instead of stdout ('-o file.txt' and '-o file.txt.tdf' both write the NanoTDF as file.txt.tdf).
      default: ''
    - name: attr
      shorthand: a
      description: Attribute value Fully Qualified Names (FQNs, i.e. 'https://example.com/attr/attr1/value/value1') to apply to the encrypted data.
---

Build a Nano Trusted Data Format (NanoTDF) with encrypted content from a specified file or input from stdin utilizing OpenTDF platform.

## Examples:

```bash
# output to stdout
echo "some text" | otdfctl encrypt-nano
otdfctl encrypt-nano hello.txt
# pipe stdout to a bucket
echo "my secret" | otdfctl encrypt-nano | aws s3 cp - s3://my-bucket/secret.txt.tdf

# output hello.txt.tdf in root directory
echo "hello world" | otdfctl encrypt-nano -o hello.txt
cat hello.txt | otdfctl encrypt-nano -o hello.txt
cat hello.txt | otdfctl encrypt-nano -o hello.txt.tdf #.tdf extension is only added once
```
