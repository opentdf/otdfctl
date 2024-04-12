---
title: Decrypt a TDF file
command:
  name: decrypt
  flags:
    - name: file
      shorthand: f
      description: The TDF file with path from the current working directory being decrypted (default 'sensitive.txt.tdf')
      default: 'sensitive.txt.tdf'
    - name: output
      shorthand: o
      description: "The decrypted output destination. Default: 'file'. Options: ['file', 'stdout']"
      default: 'stdout'
---

Decrypt a Trusted Data Format (TDF) file and output the contents to stdout or a file in the current working directory.
