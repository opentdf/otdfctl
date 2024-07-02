---
title: Update an attribute definition
command:
  name: update
  flags:
    - name: id
      shorthand: i
      description: ID of the attribute definition
      required: true
    - name: name
      shorthand: n
      description: Name of the attribute definition (new)
      required: false
    - name: rule
      shorthand: r
      description: Rule of the attribute definition (new)
      required: false
    - name: values-order
      shorthand: o
      description: Order of the attribute values (new)
      required: false
---

# Unsafe Update Warning

## Name Update

Renaming an Attribute Definition means any Values and any associated mappings underneath will now be tied to the new name.

Any existing TDFs containing attributes under the old definition name will be rendered inaccessible, and any TDFs tied to the new name
and already created may now become accessible.

## Rule Update

Altering a rule of an Attribute Definition changes the evaluation of entitlement to data. Existing TDFs of the same definition name
and values will now be accessible based on the updated rule. An `anyOf` rule becoming `hierarchy` or vice versa, for example, have
entirely different meanings and access evaluations.

## Values-Order Update

In the case of a `hierarchy` Attribute Definition Rule, the order of Values on the attribute has significant impact on data access.
Changing this order (complete, destructive replacement of the existing order) will impact access to data.

Make sure you know what you are doing.
