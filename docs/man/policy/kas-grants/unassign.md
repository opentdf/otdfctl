---
title: Unassign a grant

command:
  name: unassign
  aliases:
    - delete
    - remove
  description: Remove a grant assignment of a KAS to an Attribute Definition or Value
  flags:
    - name: namespace-id
      shorthand: n
      description: The ID of the Namespace being unassigned a KAS Grant
    - name: attribute-id
      shorthand: a
      description: The ID of the Attribute Definition being unassigned the KAS grant
      required: true
    - name: value-id
      shorthand: v
      description: The ID of the Value being unassigned the KAS Grant
      required: true
    - name: kas-id
      shorthand: k
      description: The Key Access Server (KAS) ID being unassigned a grant
      required: true
---

Unassign a registered Key Access Server (KAS) to an attribute definition or value.

For more information, see `kas-registry` and `kas-grants` manuals.
