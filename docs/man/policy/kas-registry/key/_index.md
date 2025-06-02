---
title: Key management for KAS Registry

command:
  name: key
  hidden: true
  aliases:
    - k
  flags:
    - name: json
      description: Output the result of a subcommand in JSON format (overrides configured output format). This is an inherited flag.
      default: 'false'
---

Provides a set of subcommands for managing cryptographic keys within the Key Access Server (KAS) registry.
These keys are essential for encryption and decryption operations within the OpenTDF platform.
Operations include creating, retrieving, listing, updating, and managing the platform's base key.
