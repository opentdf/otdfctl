---
title: Deactivate an attribute value
command:
  name: deactivate
  flags:
    - name: id
      shorthand: i
      description: The ID of the attribute value to deactivate
---

Deactivation preserves uniqueness of the attribute value within policy and all existing relations, essentially reserving it.

However, a deactivation of an attribute value means it cannot be entitled in an access decision.

For information about reactivation, see the `unsafe reactivate` subcommand.

For more information on attribute values, see the `values` subcommand.
