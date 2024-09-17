---
title: Create an attribute namespace
command:
  name: create
  aliases:
    - c
    - add
    - new
  flags:
    - name: name
      shorthand: n
      description: Name of the attribute namespace
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

Creation of a `namespace` is required to add attributes or any other policy objects beneath.

For more information, see the `namespaces` subcommand.
