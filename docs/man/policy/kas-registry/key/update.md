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
    - name: status
      shorthand: s
      description: The status of the key
    - name: label
      shorthand: l
      description: Metadata labels for the provider config
---

Updates an existing key KAS key.

1. The `"status"` specifies the key status:

| Key Status     |
| -------------- |
| `active`       |
| `inactive`     |
| `compromised`  |
