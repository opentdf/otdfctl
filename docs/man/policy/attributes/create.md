---
title: Create an attribute
command:
  name: create
  aliases:
    - new
    - add
    - c
  flags:
    - name: name
      shorthand: n
      description: Name of the attribute
      required: true
    - name: rule
      shorthand: r
      description: Rule of the attribute
      enum:
        - ANY_OF
        - ALL_OF
        - HIERARCHY
      required: true
    - name: value
      shorthand: v
      description: Value of the attribute (i.e. 'value1')
      required: true
    - name: namespace
      shorthand: s
      description: Namespace ID of the attribute
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

# Create an attribute definition

Under a namespace, create an attribute with a rule.

### Rules

#### ANY_OF

If an Attribute is defined with logical rule `ANY_OF`, an Entity who is mapped to `any` of the associated Values of the Attribute
on TDF'd Resource Data will be Entitled.

#### ALL_OF

If an Attribute is defined with logical rule `ALL_OF`, an Entity must be mapped to `all` of the associated Values of the Attribute
on TDF'd Resource Data to be Entitled.

### HIERARCHY

If an Attribute is defined with logical rule `HIERARCHY`, an Entity must be mapped to the same level Value or a level above in hierarchy
compared to a given Value on TDF'd Resource Data. Hierarchical values are considered highest at index 0 and lowest at the last index.

For more general information about attributes, see the `attributes` subcommand.
