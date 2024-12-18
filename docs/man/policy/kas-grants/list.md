---
title: List KAS Grants

command:
  name: list
  aliases:
    - l
  description: List the Grants of KASes to Attribute Namespaces, Definitions, and Values
  flags:
    - name: kas
      shorthand: k
      description: The optional ID or URI of a KAS to filter the list
    - name: limit
      shorthand: l
      description: Limit retrieved count
    - name: offset
      shorthand: o
      description: Offset (page) quantity from start of the list
---

List the Grants of Registered Key Access Servers (KASes) to attribute namespaces, definitions,
or values.

Omitting `kas` lists all grants known to platform policy, otherwise results are filtered to
the KAS URI or ID specified by the flag value.

For more information, see `kas-registry` and `kas-grants` manuals.

## Example

```shell
otdfctl policy kas-grants list
```

```shell
  SUCCESS                                                                                                                                                                                                                   
╭─────────────────────────────────────────────────┬─────────────────────────────────────────────────┬─────────────────┬─────────────────────────────────────────────────┬──────────────────────────────────────────────────╮
│KAS ID                                           │KAS URI                                          │Assigned To      │Granted Object ID                                │Granted Object FQN                                │
├─────────────────────────────────────────────────┼─────────────────────────────────────────────────┼─────────────────┼─────────────────────────────────────────────────┼──────────────────────────────────────────────────┤
│62857b55-560c-4b67-96e3-33e4670ecb3b             │http://example.com/kas                           │Definition       │a21eb299-3a7d-4035-8a39-c8662c03cb15             │https://opentdf.io/attr/myattribute               │
│62857b55-560c-4b67-96e3-33e4670ecb3b             │http://example.com/kas                           │Value            │0a40b27c-6cc9-49e8-a6ae-663cac2c324b             │https://opentdf.io/attr/myattribute/value/myvalue2│
│62857b55-560c-4b67-96e3-33e4670ecb3b             │http://example.com/kas                           │Namespace        │3d25d33e-2469-4990-a9ed-fdd13ce74436             │https://opentdf.io                                │
╰─────────────────────────────────────────────────┴─────────────────────────────────────────────────┴─────────────────┴─────────────────────────────────────────────────┴──────────────────────────────────────────────────╯
                          
```
