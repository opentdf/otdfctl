---
title: Decrypt a TDF file
command:
  name: decrypt [file]
  flags:
    - name: out
      shorthand: o
      description: 'The file destination for decrypted content to be written instead of stdout.'
      default: ''
    - name: tdf-type
      shorthand: t
      description:  The type of tdf to decrypt as
      enum:
        - tdf3
        - nano
      default: tdf3
---

Decrypt a Trusted Data Format (TDF) file and output the contents to stdout or a file in the current working directory.

The first argument is the TDF file with path from the current working directory being decrypted.

## Examples:

```bash
# specify the TDF to decrypt then output decrypted contents
otdfctl decrypt hello.txt.tdf # write to stdout
otdfctl decrypt hello.txt.tdf > hello.txt # consume stdout to write to hello.txt file
otdfctl decrypt hello.txt.tdf -o hello.txt # write to hello.txt file instead of stdout

# pipe the TDF to decrypt
cat hello.txt.tdf | otdfctl decrypt > hello.txt
```
