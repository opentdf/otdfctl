---
title: Encrypt file, text, or stdin as a TDF
command:
  name: encrypt
  flags:
    - name: file
      shorthand: f
      description: A file to encrypt at a path relative to $HOME.
      default: ''
    - name: text
      shorthand: t
      description: A string of text to encrypt.
      default: ''
    - name: out
      shorthand: o
      description: A filename and extension that will be TDFd (i.e. '-o password.txt' -> 'password.txt.tdf', default 'sensitive.txt.tdf' or <file>.tdf) and placed in $HOME.
      default: '' # default is set dynamically to allow filename parsing
    - name: attr
      shorthand: a
      description: Attribute value Fully Qualified Names (FQNs, i.e. 'https://example.com/attr/attr1/value/value1') to apply to the encrypted data.
---

Build a Trusted Data Format (TDF) with encrypted content from a file, string of text, or input from stdin utilizing OpenTDF platform.
