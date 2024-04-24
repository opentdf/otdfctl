---
title: Encrypt file or stdin as a TDF
command:
  name: encrypt [file]
  flags:
    - name: out
      shorthand: o
      description: A filename and extension that will be TDFd (i.e. '-o password.txt' -> 'password.txt.tdf', default 'sensitive.txt.tdf' or <file>.tdf) and placed in the current working directory.
      default: '' # default is set dynamically to allow filename parsing
    - name: attr
      shorthand: a
      description: Attribute value Fully Qualified Names (FQNs, i.e. 'https://example.com/attr/attr1/value/value1') to apply to the encrypted data.
---

Build a Trusted Data Format (TDF) with encrypted content from a specified file or input from stdin utilizing OpenTDF platform.

## Examples:

```bash
# default to sensitive.txt.tdf
echo "some text" | otdfctl encrypt

# output hello.txt.tdf in root directory
echo "hello world" | otdfctl encrypt -o hello.txt
cat hello.txt | otdfctl encrypt -o hello.txt
otdfctl encrypt hello.txt
```

The `.tdf` itself is always added to the directory where this tool is executed and provided to stdout for piping.
