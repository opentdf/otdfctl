---
title: Update an attribute namespace
command:
  name: update
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute namespace
      required: true
    - name: name
      shorthand: n
      description: Name of the attribute namespace
      required: true
---

# Unsafe Update Warning

Renaming a Namespace means any Attribute Definitions, Values, and any associated mappings underneath will now be tied to the new name.

Any existing TDFs containing attributes under the old namespace will be rendered inaccessible, and any TDFs tied to the new namespace
and already created may now become accessible.

Make sure you know what you are doing.
