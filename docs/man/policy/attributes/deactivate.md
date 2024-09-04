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

Deactivation preserves uniqueness of the attribute within policy and all existing relations, essentially reserving them.

However, a deactivation of an attribute means its associated values cannot be entitled in an access decision.

For more general information about attributes, see the `attributes` subcommand.
