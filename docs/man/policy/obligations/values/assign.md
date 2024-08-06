---
title: Assign an obligation to an attribute value
command:
  name: assign
  flags:
    - name: id
      shorthand: i
      description: ID of the obligation value
      required: true
    - name: attr-val
      description: ID of the attribute value(s) being assigned for derived obligations
      required: true
---

Assigns an existing obligation value to one or more attribute values for derived obligations in an access decision.

For more information about the significance of obligations and how they are utilized for derived obligations on attribute values,
see the parent command.
