---
title: Unassign an obligation from an attribute value
command:
  name: unassign
  flags:
    - name: id
      shorthand: i
      description: ID of the obligation value
      required: true
    - name: attr-val
      description: ID of the attribute value being removed for derived obligation assignment
      required: true
---

Unassigns an obligation value from an attribute value so that the obligation is no longer considered a derived obligation in an
access decision.

For more information about the significance of obligations and how they are utilized for derived obligations on attribute values,
see the parent command.