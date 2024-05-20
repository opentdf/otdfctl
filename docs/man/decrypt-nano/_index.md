---
title: Decrypt a NanoTDF file
command:
  name: decrypt-nano [file]
  flags:
    - name: out
      shorthand: o
      description: 'The file destination for decrypted content to be written instead of stdout.'
      default: ''
---

Decrypt a Nano Trusted Data Format (NanoTDF) file and output the contents to stdout or a file in the current working directory.

The first argument is the NanoTDF file with path from the current working directory being decrypted.

## Examples:

```bash
# specify the NanoTDF to decrypt then output decrypted contents
otdfctl decrypt-nano hello.txt.tdf # write to stdout
otdfctl decrypt-nano hello.txt.tdf > hello.txt # consume stdout to write to hello.txt file
otdfctl decrypt-nano hello.txt.tdf -o hello.txt # write to hello.txt file instead of stdout

# pipe the NanoTDF to decrypt
cat hello.txt.tdf | otdfctl decrypt-nano > hello.txt
```
