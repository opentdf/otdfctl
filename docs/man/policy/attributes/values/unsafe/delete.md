---
title: Delete an attribute value
command:
  name: delete
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute value
      required: true
---

# Unsafe Delete Warning

Deleting an Attribute Value cascades deletion of any associated mappings underneath.

Any existing TDFs containing the deleted attribute of this value will be rendered inaccessible until it has been recreated.

Make sure you know what you are doing.

For more information on attribute values, see the `values` subcommand.
