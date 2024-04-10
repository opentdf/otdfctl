---
title: Generate a set of selector expressions for keys and values of a Subject Context
command:
  name: gen
  flags:
    - name: subject
      shorthand: s
      description: A Subject Context string (JSON or JWT, default JSON)
      default: ''
    - name: type
      shorthand: t
      description: 'The type of the Subject Context: [json, jwt]'
      default: json
---

Take in a representation of some Subject Context, such as that provided by
an Identity Provider (idP), LDAP, or OIDC Access Token JWT, and generate
sample [jq syntax expressions](https://jqlang.github.io/jq/manual/) to employ
within Subject Condition Sets to parse that external Subject Context into mapped Attribute
Values.
