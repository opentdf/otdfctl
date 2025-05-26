---
title: Base Key Operations

command:
  name: base
  hidden: false
  aliases:
    - k
  flags:
    - name: json
      description: output single command in JSON (overrides configured output format)
      default: 'false'
---

Set of operations to be used for setting and getting base platform keys.
These base platform keys will be used to encrypt data in the following cases:

- No attributes present when encrypting a file
- No keys associated with an attribute
