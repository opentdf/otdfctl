---
title: Deactivate an attribute namespace
command:
  name: deactivate
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute namespace
      required: true
    - name: force
      description: Force deletion without interactive confirmation (dangerous)
---

# Deactivate an attribute namespace

Deactivating an Attribute Namespace will make the namespace name inactive as well as any attribute definitions and values beneath.

Deactivation of a Namespace renders any existing TDFs of those attributes inaccessible.

Deactivation will permanently reserve the Namespace name within a platform. Reactivation and deletion are both considered "unsafe"
behaviors.

For reactivation, see the `unsafe` command.
