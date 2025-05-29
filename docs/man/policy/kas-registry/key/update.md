---
title: Update KAS Key
command:
  name: update
  aliases:
    - u
  flags:
    - name: id
      shorthand: i
      description: ID of the key to update
      required: true
    - name: label
      shorthand: l
      description: Metadata labels for the provider config
---

Updates an existing key KAS key's metadata.
