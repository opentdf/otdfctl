---
title: Reactivate an attribute namespace
command:
  name: reactivate
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute namespace
      required: true
---

# Unsafe Reactivate Warning

Reactivating a Namespace can potentially open up an access path to any existing TDFs referencing attributes under that Namespace.

The Active/Inactive state of any Attribute Definitions or Values under this Namespace will NOT be changed.

Make sure you know what you are doing.
