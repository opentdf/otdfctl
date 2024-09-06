---
title: Deactivate an attribute
command:
  name: deactivate
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute
      required: true
---

# Deactivate an attribute definition

Deactivation preserves uniqueness of the attribute and values underneath within policy and all existing relations,
essentially reserving them.

However, a deactivation of an attribute means its associated values cannot be entitled in an access decision.

For information about reactivation, see the `unsafe reactivate` subcommand.

For more general information about attributes, see the `attributes` subcommand.
