---
title: Update an attribute value
command:
  name: update
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute value
      required: true
    - name: value
      shorthand: v
      description: The new value replacing the current value
---

# Unsafe Update Warning

## Value Update

Changing an Attribute Value means any associated mappings underneath will now be tied to the new value.

Any existing TDFs containing attributes under the old value will be rendered inaccessible, and any TDFs tied to the new value
and already created may now become accessible.

Make sure you know what you are doing.
