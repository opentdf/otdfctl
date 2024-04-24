---
title: Decrypt a TDF file
command:
  name: decrypt [file]
  flags:
    - name: out
      shorthand: o
      description: 'The file destination for decrypted content to be written.'
      default: ''
---

Decrypt a Trusted Data Format (TDF) file and output the contents to stdout or a file in the current working directory.

The first argument is the TDF file with path from the current working directory being decrypted (default 'sensitive.txt.tdf').

Examples:

```bash
# default to sensitive.txt.tdf, then print to stdout
otdfctl decrypt

# specify the TDF to decrypt then output decrypted contents
otdfctl decrypt hello.txt.tdf # print to stdout
otdfctl decrypt hello.txt.tdf > hello.txt # consume stdout to write to hello.txt file
otdfctl decrypt hello.txt.tdf -o hello.txt # write to hello.txt file
```
