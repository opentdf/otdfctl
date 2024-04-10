---
title: Decrypt a TDF
command:
  name: decrypt
  flags:
    - name: tdf
      shorthand: t
      description: The TDF file with path from $HOME being decrypted (default 'sensitive.txt.tdf')
      default: 'sensitive.txt.tdf'
    - name: output
      shorthand: o
      description: "The decrypted output destination (default 'file', options: 'file', 'stdout')"
      default: 'file'
---

Decrypt a TDF and output the contents to stdout or a file without the .tdf extension in the same directory as the TDF.
