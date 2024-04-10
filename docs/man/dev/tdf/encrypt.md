---
title: Encrypt file or text
command:
  name: encrypt
  flags:
    - name: file
      shorthand: f
      description: A file to encrypt that will be encrypted and saved as '<filename>.<extension>.tdf' in the same directory.
      default: ''
    - name: text
      shorthand: t
      description: A string of text to encrypt that will be saved as 'sensitive.txt.tdf' in the $HOME directory.
      default: ''
    - name: attr-value
      shorthand: v
      description: Attribute value Fully Qualified Names (FQNs, i.e. 'https://example.com/attr/attr1/value/value1') to apply to the encrypted data.
---

Encrypt a file or string of text utilizing the TDF and OpenTDF platform.