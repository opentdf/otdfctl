---
title: Assign a grant

command:
  name: assign
  aliases:
    - u
    - update
    - create
    - add
    - new
    - upsert
  description: Assign a grant of a KAS to an Attribute Definition or Value
  flags:
    - name: namespace-id
      shorthand: n
      description: The ID of the Namespace being assigned a KAS Grant
    - name: attribute-id
      shorthand: a
      description: The ID of the Attribute Definition being assigned a KAS Grant
      required: true
    - name: value-id
      shorthand: v
      description: The ID of the Value being assigned a KAS Grant
      required: true
    - name: kas-id
      shorthand: k
      description: The ID of the Key Access Server being assigned to the grant
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
    - name: force-replace-labels
      description: Destructively replace entire set of existing metadata 'labels' with any provided to this command
      default: false
---

Assign a registered Key Access Server (KAS) to an attribute namespace, definition, or value.

For more information, see `kas-registry` and `kas-grants` manuals.