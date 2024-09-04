---
title: Delete a subject mapping by id
command:
  name: delete
  flags:
    - name: id
      description: The ID of the subject mapping to delete
      shorthand: i
      required: true
      default: ''
---

# Delete a subject mapping

Delete a Subject Mapping to remove entitlement of an entity (via Subject Condition Set) to an Attribute Value.

For more information about subject mappings, see the `subject-mappings` subcommand.

For more information about subject condition sets, see the `subject-condition-sets` subcommand.