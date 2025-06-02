---
title: Platform Base Key Management

command:
  name: base
  hidden: true
  flags:
    - name: json
      description: Output the result of a subcommand in JSON format (overrides configured output format). This is an inherited flag.
      default: 'false'
---

Provides subcommands for managing the platform's base cryptographic key.
This base key is a fallback used for encryption operations in specific scenarios:

- No attributes present when encrypting a file
- No keys associated with an attribute

Available operations include `get` to retrieve the current base key and `set` to designate a new base key.
