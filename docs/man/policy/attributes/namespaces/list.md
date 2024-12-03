---
title: List attribute namespaces
command:
  name: list
  aliases:
    - ls
    - l
  flags:
    - name: state
      shorthand: s
      description: Filter by state [active, inactive, any]
    - name: limit
      shorthand: l
      description: Limit retrieved count (default set by platform if not provided)
    - name: offset
      shorthand: o
      description: Offset quantity from start of the list (page)
---

For more general information, see the `namespaces` subcommand.
