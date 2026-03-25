---
title: Create a resource mapping
command:
  name: create
  aliases:
    - add
    - new
    - c
  flags:
    - name: attribute-value-id
      description: The ID of the attribute value to map to the resource.
      default: ''
    - name: terms
      description: The synonym terms to match for the resource mapping.
      default: ''
    - name: group-id
      description: The ID of the resource mapping group to assign this mapping to
      default: ''
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

Associate an attribute value with a set of plaintext string terms.

For more information about resource mappings, see the `resource-mappings` subcommand.

## Namespacing subject mappings

The following rules must be applied when attempting to namespace a subject mapping:

- Either all policy constructs (action, subject mappings, subject condition set, attribute value) are within the same
  namespace
- Subject mapping, subject condition set, action are all not within a namespace.

You cannot, for example:

- Create a subject mapping that is not within the same namespace as an action that is passed in

## Examples

```shell
otdfctl policy resource-mappings create --attribute-value-id 891cfe85-b381-4f85-9699-5f7dbfe2a9ab --terms term1,term2 --group-id 3ff446fb-8fb1-4c04-8023-47592c90370c


# Namespaced subject condition set creation
otdfctl policy resource-mappings create --attribute-value-id 891cfe85-b381-4f85-9699-5f7dbfe2a9ab --terms term1,term2 --group-id 3ff446fb-8fb1-4c04-8023-47592c90370c --namespace "https://example.com"
```
