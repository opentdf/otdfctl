---
title: Create an obligation fulfillment
command:
  name: create
  flags:
    - name: conditions-json
      description: TODO - Conditions as defined by the protos as JSON
      required: false
    - name: conditions-json-file
      description: TODO - Conditions as defined by the protos from a JSON file
      required: false
    - name: obligation-scope
      shorthand: s
      description: Scope of the obligation as subject or environment [ SUBJECT, ENV ]
      required: true
    - name: label
      description: "Optional metadata 'labels' in the format: key=value"
      shorthand: l
      default: ''
---

For more information about the significance of obligations and how they are utilized for derived obligations on attribute values
or when added directly to a TDF, see the parent command above.
