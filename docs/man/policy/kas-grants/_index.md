---
title: Manage Key Access Server grants

command:
  name: kas-grants
  aliases:
    - kasg
    - kas-grant
---

## Background

Once Key Access Servers (KASs) have been registered within a platform's policy,
they can be assigned grants to various attribute objects (namespaces, definitions, values).

> See `kas-registry` command within `policy` to manage the KASs known to the platform.

Grants are utilized by the platform at two important points when engaging with a TDF.

## Utilization

The steps below are driven by the SDK on encrypt, and they are the same steps followed
on decrypt by a KAS making a decision request on a key release (once the decision
is found to be permissible):

1. look up the attributes on the TDF within the platform
2. find any associated grants for those attributes' values, definitions, namespaces
3. retrieve the public key of each KAS granted to those attribute objects
4. determine based on the specificity matrix below which keys to utilize

## Specificity

When KAS grants are considered, they follow a most-to-least specificity matrix. Grants to
Attribute Values supersede any grants to Definitions which also supersede any grants to a Namespace.

Grants to Attribute Objects:

| Namespace Grant | Attr Definition Grant | Attr Value Grant | Data Encryption Key Utilized |
| --------------- | --------------------- | ---------------- | ---------------------------- |
| yes             | no                    | no               | namespace                    |
| yes             | yes                   | no               | attr definition              |
| no              | yes                   | no               | attr definition              |
| yes             | yes                   | yes              | value                        |
| no              | yes                   | yes              | value                        |
| no              | no                    | yes              | value                        |
| no              | no                    | no               | default KAS/platform key     |

> Note:
> A namespace grant may soon be required with deprecation of a default KAS/platform key.
