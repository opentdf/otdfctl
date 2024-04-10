---
title: Test resolution of a set of selector expressions for keys and values of a Subject Context.
command:
  name: test
  flags:
    - name: subject
      shorthand: s
      description: A Subject Context string (JSON or JWT, default JSON)
      default: ''
    - name: type
      shorthand: t
      description: 'The type of the Subject Context: [json, jwt]'
      default: json
    - name: selector
      shorthand: x
      description: 'Individual selectors to test against the Subject Context (i.e. .key, .example[1].group)'
---

Test a given representation of some Subject Context, such as that provided by
an Identity Provider (idP), LDAP, or OIDC Access Token JWT, against provided [jq syntax
'selector' expressions](https://jqlang.github.io/jq/manual/) to validate their resolution
to field values on the Subject Context.
