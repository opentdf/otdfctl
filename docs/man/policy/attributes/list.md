---
title: List attributes
command:
  name: list
  aliases:
    - l
  flags:
    - name: state
      shorthand: s
      description: Filter by state
      enum:
        - active
        - inactive
        - any
      default: active
---

# List the known attributes

By default, the list will only provide `active` attributes if unspecified, but the filter can be controlled with the `--state` flag.

For more general information about attributes, see the `attributes` subcommand.
