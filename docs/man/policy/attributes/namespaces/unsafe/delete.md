---
title: Delete an attribute namespace
command:
  name: delete
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute namespace
      required: true
---

# Unsafe Delete Warning

Deleting a Namespace cascades deletion of any Attribute Definitions, Values, and any associated mappings underneath.

Any existing TDFs containing attributes under this namespace will be rendered inaccessible until it has been recreated.

Make sure you know what you are doing.

For more general information, see the `namespaces` subcommand.
